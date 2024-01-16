package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
	"wachat/database"
	"wachat/params"
	"wachat/pb"
	"wachat/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/waku-org/go-waku/waku/v2/dnsdisc"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/payload"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	wpb "github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	wakuutils "github.com/waku-org/go-waku/waku/v2/utils"
	"google.golang.org/protobuf/proto"
)

// App struct
type App struct {
	ctx      context.Context
	node     *node.WakuNode
	topic    protocol.ContentTopic
	username string
	isOnline bool
	sqlite   *sqlx.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	DbMigrate()

	sqlDbPath, err := utils.SQLiteDatabasePath()
	if err != nil {
		// utils.Sugar.Fatal(err)
	}
	sqlDb, err := sqlx.Connect("sqlite", sqlDbPath)
	if err != nil {
		// utils.Sugar.Fatal(err)
	}
	a.sqlite = sqlDb

	user, err := database.GetSelectedUser(a.sqlite)
	if err == nil {
		a.username = user.Name
	}

	a.isOnline = a.isNetworkOnline()
	if a.isOnline {
		a.startWaku(ctx)
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			currentStatus := a.isNetworkOnline()
			if a.isOnline != currentStatus {
				fmt.Println("Network status previous and current", a.isOnline, currentStatus)
				a.isOnline = currentStatus
				runtime.EventsEmit(a.ctx, "isOnline", currentStatus)
				if currentStatus {
					fmt.Println("Network is online")
					a.node.Stop()
					a.startWaku(ctx)
				}
			}
		}
	}()
}

func (a *App) startWaku(ctx context.Context) {
	contentTopic, err := protocol.NewContentTopic("toy-chat", "3", "mingde", "proto")
	if err != nil {
		fmt.Println("Invalid Content Topic")
		panic(err)
	}
	hostAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	key, err := randomHex(32)
	if err != nil {
		fmt.Println("Could not generate random key")
		panic(err)
	}
	prvKey, err := crypto.HexToECDSA(key)
	if err != nil {
		fmt.Println("Could not generate private key")
		panic(err)
	}
	wakuNode, err := node.New(
		node.WithPrivateKey(prvKey),
		node.WithHostAddress(hostAddr),
		node.WithNTP(),
		node.WithWakuRelay(),
		node.WithWakuRelayAndMinPeers(1),
	)
	if err != nil {
		fmt.Println("Could not create waku node")
		panic(err)
	}

	if err := wakuNode.Start(ctx); err != nil {
		fmt.Println("Could not start waku node")
		panic(err)
	}

	a.node = wakuNode
	a.topic = contentTopic

	go a.readMessages()
	go a.discoverNodes()
	go a.checkMessages()
}

func (a *App) shutdown(ctx context.Context) {
	a.node.Stop()
}

// Greet returns a greeting for the given name
func (a *App) Send(message string) (string, error) {
	wakuPayload := new(payload.Payload)
	pbMessage := &pb.Chat2Message{
		Timestamp: uint64(a.node.Timesource().Now().Unix()),
		Nick:      a.username,
		Payload:   []byte(message),
	}
	pbMsgBytes, err := proto.Marshal(pbMessage)
	if err != nil {
		return "", err
	}
	wakuPayload.Data = pbMsgBytes
	wakuPayload.Key = &payload.KeyInfo{Kind: payload.None}

	payloadBytes, err := wakuPayload.Encode(utils.MessageVersion)
	if err != nil {
		log.Error("Error encoding the payload", err)
		return "", err
	}

	msg := &wpb.WakuMessage{
		Payload:      payloadBytes,
		Version:      proto.Uint32(utils.MessageVersion),
		ContentTopic: a.topic.String(),
		Timestamp:    wakuutils.GetUnixEpoch(a.node.Timesource()),
	}

	msgHash, err := a.node.Relay().Publish(a.ctx, msg, relay.WithDefaultPubsubTopic())
	if err != nil {
		log.Error("Error push a message", err)
		return "", err
	}

	return hexutil.Encode(msgHash), nil
}

func (a *App) isNetworkOnline() bool {
	// Set a timeout for the HTTP request
	online := true
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Attempt to perform a GET request to a known server
	_, err := client.Get(utils.ConnectivityCheckServer)
	if err != nil {
		fmt.Println("Error:", err)
		online = false
	}

	return online
}

func (a *App) CreateUser(name string) (string, error) {
	user := params.User{
		Name:     name,
		Selected: true,
	}
	err := database.SaveUser(a.sqlite, user)
	if err != nil {
		return "", err
	}
	a.username = name
	return name, nil
}

func (a *App) GetUser() string {
	return a.username
}

func (a *App) GetMessages() []params.Message {
	return []params.Message{}
}

func (a *App) readMessages() {
	sub, err := a.node.Relay().Subscribe(a.ctx, protocol.NewContentFilter(relay.DefaultWakuTopic))
	if err != nil {
		log.Error("Error subscribing to the default waku topic", err)
		return
	}

	for envelope := range sub[0].Ch {
		if envelope.Message().ContentTopic != a.topic.String() {
			continue
		}

		msgPayload, err := payload.DecodePayload(envelope.Message(), &payload.KeyInfo{Kind: payload.None})
		if err != nil {
			log.Error("Error decoding the payload", err)
			continue
		}

		fmt.Println("Received message: ", string(msgPayload.Data))

		msgDecoded := &pb.Chat2Message{}
		if err := proto.Unmarshal(msgPayload.Data, msgDecoded); err != nil {
			log.Error("Error decoding the payload", err)
			continue
		}

		msg := params.Message{
			Hash:          hexutil.Encode(envelope.Hash()),
			Content:       string(msgDecoded.Payload),
			Name:          msgDecoded.Nick,
			Timestamp:     msgDecoded.Timestamp,
			WakuTimestamp: uint64(envelope.Message().GetTimestamp()),
		}
		err = database.SaveMessage(a.sqlite, msg)
		if err != nil {
			fmt.Println("Error saving message", err)
			continue
		}

		fmt.Println("emit msg", msg)
		runtime.EventsEmit(a.ctx, "newMessage", msg)
	}
}

func (a *App) discoverNodes() {
	dnsDiscoveryUrl := "enrtree://AL65EKLJAUXKKPG43HVTML5EFFWEZ7L4LOKTLZCLJASG4DSESQZEC@prod.status.nodes.status.im"
	nodes, err := dnsdisc.RetrieveNodes(a.ctx, dnsDiscoveryUrl)
	if err != nil {
		log.Error("Error retrieving nodes", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(nodes))
	for _, node := range nodes {
		go func(ctx context.Context, info peer.AddrInfo) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, time.Duration(30)*time.Second)
			defer cancel()

			err = a.node.DialPeerWithInfo(ctx, info)
			if err != nil {
				fmt.Println("Error dialing peer", err)
				return
			}
		}(a.ctx, node.PeerInfo)
	}
	wg.Wait()
}

func (a *App) checkMessages() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("Checking messages")
			a.ensureMessageStored()
		}
	}
}

func (a *App) ensureMessageStored() {
	messages, err := database.GetUnstoredMessages(a.sqlite)
	if err != nil {
		fmt.Println("Error getting unstored messages", err)
		return
	}

	for _, msg := range messages {
		fmt.Println("Retrieve unstored message", msg.Hash)
		now := wakuutils.GetUnixEpoch(a.node.Timesource())
		query := store.Query{
			StartTime:     proto.Int64(int64(msg.Timestamp-86400) * int64(time.Second)),
			EndTime:       proto.Int64(int64(*now) * int64(time.Second)),
			ContentTopics: []string{a.topic.String()},
			PubsubTopic:   relay.DefaultWakuTopic,
		}
		fn := func(wakuMsg *wpb.WakuMessage) (bool, error) {
			wakuMsgHash := wakuMsg.Hash(relay.DefaultWakuTopic)

			return msg.Hash == hexutil.Encode(wakuMsgHash), nil
		}
		found, err := a.node.Store().Find(a.ctx, query, fn)
		if err != nil {
			fmt.Println("Error searching for message", err)
			continue
		}

		if found == nil {
			fmt.Println("Message not found, republishing", msg.Hash)
			wakuPayload := new(payload.Payload)
			pbMessage := &pb.Chat2Message{
				Timestamp: msg.Timestamp,
				Nick:      msg.Name,
				Payload:   []byte(msg.Content),
			}
			pbMsgBytes, err := proto.Marshal(pbMessage)
			if err != nil {
				continue
			}
			wakuPayload.Data = pbMsgBytes
			wakuPayload.Key = &payload.KeyInfo{Kind: payload.None}

			payloadBytes, err := wakuPayload.Encode(utils.MessageVersion)
			if err != nil {
				log.Error("Error encoding the payload", err)
				continue
			}
			msg := &wpb.WakuMessage{
				Payload:      payloadBytes,
				Version:      proto.Uint32(utils.MessageVersion),
				ContentTopic: a.topic.String(),
				Timestamp:    proto.Int64(int64(msg.WakuTimestamp)),
			}

			_, err = a.node.Relay().Publish(a.ctx, msg, relay.WithDefaultPubsubTopic())
			if err != nil {
				log.Error("Error push a message", err)
				continue
			}

			continue
		}

		fmt.Println("Message found, updating", msg.Hash)
		database.UpdateStoredMessage(a.sqlite, msg.Hash)
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
