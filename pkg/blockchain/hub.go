package blockchain

import (
	"context"
	"crypto/sha256"
	"errors"
	"strings"
	"sync"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type MessageHub struct {
	sync.Mutex

	blockchain, public *room
	ps                 *pubsub.PubSub
	otpKey             string
	maxsize            int
	keyLength          int
	interval           int
	joinPublic         bool

	ctxCancel                context.CancelFunc
	Messages, PublicMessages chan *Message
}

// roomBufSize is the number of incoming messages to buffer for each topic.
const roomBufSize = 128

func NewHub(otp string, maxsize, keyLength, interval int, joinPublic bool) *MessageHub {
	return &MessageHub{otpKey: otp, maxsize: maxsize, keyLength: keyLength, interval: interval,
		Messages: make(chan *Message, roomBufSize), PublicMessages: make(chan *Message, roomBufSize), joinPublic: joinPublic}
}

func (m *MessageHub) topicKey(salts ...string) string {
	totp := TOTP(sha256.New, m.keyLength, m.interval, m.otpKey)
	if len(salts) > 0 {
		return MD5(totp + strings.Join(salts, ":"))
	}
	return MD5(totp)
}

func (m *MessageHub) joinRoom(host host.Host) error {
	m.Lock()
	defer m.Unlock()

	if m.ctxCancel != nil {
		m.ctxCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.ctxCancel = cancel

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, host, pubsub.WithMaxMessageSize(m.maxsize))
	if err != nil {
		return err
	}

	// join the "chat" room
	cr, err := connect(ctx, ps, host.ID(), m.topicKey(), m.Messages)
	if err != nil {
		return err
	}

	m.blockchain = cr

	if m.joinPublic {
		cr2, err := connect(ctx, ps, host.ID(), m.topicKey("public"), m.PublicMessages)
		if err != nil {
			return err
		}
		m.public = cr2
	}

	m.ps = ps

	return nil
}

func (m *MessageHub) Start(ctx context.Context, host host.Host) error {
	c := make(chan interface{})
	go func(c context.Context, cc chan interface{}) {
		k := ""
		for {
			select {
			default:
				currentKey := m.topicKey()
				if currentKey != k {
					k = currentKey
					cc <- nil
				}
				time.Sleep(1 * time.Second)
			case <-ctx.Done():
				close(cc)
				return
			}
		}
	}(ctx, c)

	for range c {
		m.joinRoom(host)
	}

	// Close eventual open contexts
	if m.ctxCancel != nil {
		m.ctxCancel()
	}
	return nil
}

func (m *MessageHub) PublishMessage(mess *Message) error {
	m.Lock()
	defer m.Unlock()
	if m.blockchain != nil {
		return m.blockchain.publishMessage(mess)
	}
	return errors.New("no message room available")
}

func (m *MessageHub) PublishPublicMessage(mess *Message) error {
	m.Lock()
	defer m.Unlock()
	if m.public != nil {
		return m.public.publishMessage(mess)
	}
	return errors.New("no message room available")
}

func (m *MessageHub) ListPeers() ([]peer.ID, error) {
	m.Lock()
	defer m.Unlock()
	if m.blockchain != nil {
		return m.blockchain.Topic.ListPeers(), nil
	}
	return nil, errors.New("no message room available")
}
