package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/payload"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
	"github.com/waku-org/go-waku/waku/v2/utils"
	"google.golang.org/protobuf/proto"
)

// App struct
type App struct {
	ctx   context.Context
	node  *node.WakuNode
	topic protocol.ContentTopic
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	contentTopic, err := protocol.NewContentTopic("wachat", "0.1", "1to1", "proto")
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
}

// Greet returns a greeting for the given name
func (a *App) Send(message string) (string, error) {
	const version = 0
	wakuPayload := new(payload.Payload)
	wakuPayload.Data = []byte(message)
	wakuPayload.Key = &payload.KeyInfo{Kind: payload.None}

	payloadBytes, err := wakuPayload.Encode(version)
	if err != nil {
		log.Error("Error encoding the payload", err)
		return "", err
	}

	msg := &pb.WakuMessage{
		Payload:      payloadBytes,
		Version:      proto.Uint32(version),
		ContentTopic: a.topic.String(),
		Timestamp:    utils.GetUnixEpoch(a.node.Timesource()),
	}

	msgHash, err := a.node.Relay().Publish(a.ctx, msg, relay.WithDefaultPubsubTopic())
	if err != nil {
		log.Error("Error push a message", err)
		return "", err
	}

	return hexutil.Encode(msgHash), nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
