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
	"wachat/pb"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/jpillora/overseer"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/waku-org/go-waku/waku/v2/dnsdisc"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/payload"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/lightpush"
	wpb "github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
	"github.com/waku-org/go-waku/waku/v2/utils"
	"google.golang.org/protobuf/proto"
)

// App struct
type App struct {
	ctx      context.Context
	node     *node.WakuNode
	topic    protocol.ContentTopic
	username string
	isOnline bool
}

type Message struct {
	Hash      string `json:"hash"`
	Content   string `json:"content"`
	Name      string `json:"name"`
	Timestamp uint64 `json:"timestamp"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	contentTopic, err := protocol.NewContentTopic("toy-chat", "2", "huilong", "proto")
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

	a.ctx = ctx
	a.node = wakuNode
	a.topic = contentTopic
	a.isOnline = a.isNetworkOnline()

	go a.readMessages()
	go a.discoverNodes()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			currentStatus := a.isNetworkOnline()
			if a.isOnline != currentStatus {
				fmt.Println("Network is online", a.isOnline, currentStatus)
				a.isOnline = currentStatus
				runtime.EventsEmit(a.ctx, "isOnline", currentStatus)
				if currentStatus {
					fmt.Println("Network is online")
					overseer.Restart()
				}
			}
		}
	}()
}

func (a *App) shutdown(ctx context.Context) {
	a.node.Stop()
}

// Greet returns a greeting for the given name
func (a *App) Send(message string) (string, error) {
	const version = 0

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

	payloadBytes, err := wakuPayload.Encode(version)
	if err != nil {
		log.Error("Error encoding the payload", err)
		return "", err
	}

	msg := &wpb.WakuMessage{
		Payload:      payloadBytes,
		Version:      proto.Uint32(version),
		ContentTopic: a.topic.String(),
		Timestamp:    utils.GetUnixEpoch(a.node.Timesource()),
	}

	msgHash, err := a.node.Lightpush().Publish(a.ctx, msg, lightpush.WithDefaultPubsubTopic())
	if err != nil {
		log.Error("Error push a message", err)
		return "", err
	}

	return hexutil.Encode(msgHash), nil
}

const testServer = "http://www.google.com"

func (a *App) isNetworkOnline() bool {
	// Set a timeout for the HTTP request
	online := true
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Attempt to perform a GET request to a known server
	_, err := client.Get(testServer)
	if err != nil {
		fmt.Println("Error:", err)
		online = false
	}

	return online
}

func (a *App) CreateUser(name string) error {
	a.username = name
	return nil
}

func (a *App) GetMessages() []Message {
	return []Message{}
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

		msg := Message{
			Hash:      hexutil.Encode(envelope.Hash()),
			Content:   string(msgDecoded.Payload),
			Name:      msgDecoded.Nick,
			Timestamp: msgDecoded.Timestamp,
		}

		fmt.Println("emit msg", msg)
		runtime.EventsEmit(a.ctx, "newMessage", msg)
	}
}

func (a *App) discoverNodes() {
	dnsDiscoveryUrl := "enrtree://ANEDLO25QVUGJOUTQFRYKWX6P4Z4GKVESBMHML7DZ6YK4LGS5FC5O@prod.wakuv2.nodes.status.im"
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

			ctx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
			defer cancel()

			err = a.node.DialPeerWithInfo(ctx, info)
			if err != nil {
				log.Error("Error dialing peer", err)
				return
			}
		}(a.ctx, node.PeerInfo)
	}
	wg.Wait()
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
