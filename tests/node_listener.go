package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/sec"
	"github.com/libp2p/go-libp2p/core/sec/insecure"
	"github.com/libp2p/go-libp2p/core/transport"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/net/upgrader"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

type NodeListener struct {
	PrivKey  crypto.PrivKey
	Address  multiaddr.Multiaddr
	ServerId peer.ID
	Upgrader transport.Upgrader
	Listener transport.Listener
}

func NewNodeListener(connString string) (*NodeListener, error) {
	multiAddr := multiaddr.StringCast(connString)
	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, err
	}
	id, u, err := newUpgrader(privKey)
	if err != nil {
		return nil, err
	}
	return &NodeListener{
		PrivKey:  privKey,
		Address:  multiAddr,
		ServerId: id,
		Upgrader: u,
	}, nil
}

func (ml *NodeListener) Start() error {
	logrus.Infof("[+] NodeListener --> Start()")
	var opts []websocket.Option
	tlsConf, err := generateTLSConfig()
	if err != nil {
		return err
	}
	opts = append(opts, websocket.WithTLSConfig(tlsConf))
	tpt, err := websocket.New(ml.Upgrader, &network.NullResourceManager{}, opts...)
	if err != nil {
		return err
	}

	l, err := tpt.Listen(ml.Address)
	if err != nil {
		return err
	}
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	rm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return err
	}

	host, err := libp2p.New(
		libp2p.Transport(websocket.New),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/3001/ws"),
		libp2p.ResourceManager(rm),
		libp2p.Identity(ml.PrivKey),
		libp2p.Ping(false), // disable built-in ping
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
	)
	if err != nil {
		return err
	}

	host.SetStreamHandler("/websocket/1.0.0", handleStream)
	ml.Listener = l
	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	multiaddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return err
	}
	addr1 := ml.Address.String()
	addr2 := multiaddrs[0].String()
	logrus.Infof("[+] libp2p host address: %s", addr1)
	logrus.Infof("[+] libp2p host address: %s", addr2)
	return nil
}

func handleStream(stream network.Stream) {
	defer stream.Close()

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		logrus.Errorf("[-] Error reading from stream: %s", err)
		return
	}

	logrus.Infof("[+] Received message: %s", string(buf[:n]))
}

func generateTLSConfig() (*tls.Config, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, priv.Public(), priv)
	return &tls.Config{
		Certificates: []tls.Certificate{{
			PrivateKey:  priv,
			Certificate: [][]byte{certDER},
		}},
	}, nil
}

func newUpgrader(privKey crypto.PrivKey) (peer.ID, transport.Upgrader, error) {
	id, err := peer.IDFromPrivateKey(privKey)
	security := []sec.SecureTransport{insecure.NewWithIdentity(insecure.ID, id, privKey)}
	if err != nil {
		return "", nil, err
	}
	upgrader, err := upgrader.New(security, []upgrader.StreamMuxer{{ID: "/yamux", Muxer: yamux.DefaultTransport}}, nil, nil, nil)
	if err != nil {
		return "", nil, err
	}
	return id, upgrader, nil
}
