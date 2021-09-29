package p2p

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	cdiscovery "github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/multierr"

	"github.com/celestiaorg/optimint/config"
	"github.com/celestiaorg/optimint/log"
)

// TODO(tzdybal): refactor to configuration parameters
const (
	// reAdvertisePeriod defines a period after which P2P client re-attempt advertising namespace in DHT.
	reAdvertisePeriod = 1 * time.Hour

	// peerLimit defines limit of number of peers returned during active peer discovery.
	peerLimit = 60

	// txTopicSuffix is added after namespace to create pubsub topic for TX gossiping.
	txTopicSuffix = "-tx"

	// headerTopicSuffix is added after namespace to create pubsub topic for block header gossiping.
	headerTopicSuffix = "-header"
)

// Client is a P2P client, implemented with libp2p.
//
// Initially, client connects to predefined seed nodes (aka bootnodes, bootstrap nodes).
// Those seed nodes serve Kademlia DHT protocol, and are agnostic to ORU chain. Using DHT
// peer routing and discovery clients find other peers within ORU network.
type Client struct {
	conf    config.P2PConfig
	chainID string
	privKey crypto.PrivKey

	host host.Host
	dht  *dht.IpfsDHT
	disc *discovery.RoutingDiscovery

	txGossip    *Gossip
	txValidator pubsub.Validator

	headerGossip *Gossip

	// cancel is used to cancel context passed to libp2p functions
	// it's required because of discovery.Advertise call
	cancel context.CancelFunc

	logger log.Logger
}

// NewClient creates new Client object.
//
// Basic checks on parameters are done, and default parameters are provided for unset-configuration
// TODO(tzdybal): consider passing entire config, not just P2P config, to reduce number of arguments
func NewClient(conf config.P2PConfig, privKey crypto.PrivKey, chainID string, logger log.Logger) (*Client, error) {
	if privKey == nil {
		return nil, errNoPrivKey
	}
	if conf.ListenAddress == "" {
		conf.ListenAddress = config.DefaultListenAddress
	}
	return &Client{
		conf:    conf,
		privKey: privKey,
		chainID: chainID,
		logger:  logger,
	}, nil
}

// Start establish Client's P2P connectivity.
//
// Following steps are taken:
// 1. Setup libp2p host, start listening for incoming connections.
// 2. Setup gossibsub.
// 3. Setup DHT, establish connection to seed nodes and initialize peer discovery.
// 4. Use active peer discovery to look for peers from same ORU network.
func (c *Client) Start(ctx context.Context) error {
	// create new, cancelable context
	ctx, c.cancel = context.WithCancel(ctx)
	c.logger.Debug("starting P2P client")
	host, err := c.listen(ctx)
	if err != nil {
		return err
	}
	return c.startWithHost(ctx, host)
}

func (c *Client) startWithHost(ctx context.Context, h host.Host) error {
	c.host = h
	for _, a := range c.host.Addrs() {
		c.logger.Info("listening on", "address", fmt.Sprintf("%s/p2p/%s", a, c.host.ID()))
	}

	c.logger.Debug("setting up gossiping")
	err := c.setupGossiping(ctx)
	if err != nil {
		return err
	}

	c.logger.Debug("setting up DHT")
	err = c.setupDHT(ctx)
	if err != nil {
		return err
	}

	c.logger.Debug("setting up active peer discovery")
	err = c.peerDiscovery(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Close gently stops Client.
func (c *Client) Close() error {
	c.cancel()

	return multierr.Combine(
		c.txGossip.Close(),
		c.dht.Close(),
		c.host.Close(),
	)
}

// GossipTx sends the transaction to the P2P network.
func (c *Client) GossipTx(ctx context.Context, tx []byte) error {
	c.logger.Debug("Gossiping TX", "len", len(tx))
	return c.txGossip.Publish(ctx, tx)
}

// SetTxHandler sets the callback function, that will be invoked after transaction is received from P2P network.
func (c *Client) SetTxHandler(handler GossipHandler) {
	c.txGossip.handler = handler
}

func (c *Client) SetTxValidator(val pubsub.Validator) {
	c.txValidator = val
}

// GossipHeader sends the block header to the P2P network.
func (c *Client) GossipHeader(ctx context.Context, headerBytes []byte) error {
	c.logger.Debug("Gossiping block header", "len", len(headerBytes))
	return c.headerGossip.Publish(ctx, headerBytes)
}

// SetHeaderHandler sets the callback function, that will be invoked after block header is received from P2P network.
func (c *Client) SetHeaderHandler(handler GossipHandler) {
	c.headerGossip.handler = handler
}

func (c *Client) listen(ctx context.Context) (host.Host, error) {
	var err error
	maddr, err := multiaddr.NewMultiaddr(c.conf.ListenAddress)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(ctx, libp2p.ListenAddrs(maddr), libp2p.Identity(c.privKey))
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (c *Client) setupDHT(ctx context.Context) error {
	seedNodes := c.getSeedAddrInfo(c.conf.Seeds)
	if len(seedNodes) == 0 {
		c.logger.Info("no seed nodes - only listening for connections")
	}

	for _, sa := range seedNodes {
		c.logger.Debug("seed node", "addr", sa)
	}

	var err error
	c.dht, err = dht.New(ctx, c.host, dht.Mode(dht.ModeServer), dht.BootstrapPeers(seedNodes...))
	if err != nil {
		return fmt.Errorf("failed to create DHT: %w", err)
	}

	err = c.dht.Bootstrap(ctx)
	if err != nil {
		return fmt.Errorf("failed to bootstrap DHT: %w", err)
	}

	c.host = routedhost.Wrap(c.host, c.dht)

	return nil
}

func (c *Client) peerDiscovery(ctx context.Context) error {
	err := c.setupPeerDiscovery(ctx)
	if err != nil {
		return err
	}

	err = c.advertise(ctx)
	if err != nil {
		return err
	}

	err = c.findPeers(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) setupPeerDiscovery(ctx context.Context) error {
	// wait for DHT
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.dht.RefreshRoutingTable():
	}
	c.disc = discovery.NewRoutingDiscovery(c.dht)
	return nil
}

func (c *Client) advertise(ctx context.Context) error {
	discovery.Advertise(ctx, c.disc, c.getNamespace(), cdiscovery.TTL(reAdvertisePeriod))
	return nil
}

func (c *Client) findPeers(ctx context.Context) error {
	peerCh, err := c.disc.FindPeers(ctx, c.getNamespace(), cdiscovery.Limit(peerLimit))
	if err != nil {
		return err
	}

	for peer := range peerCh {
		go c.tryConnect(ctx, peer)
	}

	return nil
}

// tryConnect attempts to connect to a peer and logs error if necessary
func (c *Client) tryConnect(ctx context.Context, peer peer.AddrInfo) {
	err := c.host.Connect(ctx, peer)
	if err != nil {
		c.logger.Error("failed to connect to peer", "peer", peer, "error", err)
	}
}

func (c *Client) setupGossiping(ctx context.Context) error {
	ps, err := pubsub.NewGossipSub(ctx, c.host)
	if err != nil {
		return err
	}

	if c.txValidator != nil {
		err = ps.RegisterTopicValidator(c.getTxTopic(), c.txValidator)
		if err != nil {
			return err
		}
	}

	c.txGossip, err = NewGossip(c.host, ps, c.getTxTopic(), c.logger)
	if err != nil {
		return err
	}
	go c.txGossip.ProcessMessages(ctx)

	c.headerGossip, err = NewGossip(c.host, ps, c.getHeaderTopic(), c.logger)
	if err != nil {
		return err
	}
	go c.headerGossip.ProcessMessages(ctx)

	return nil
}

func (c *Client) getSeedAddrInfo(seedStr string) []peer.AddrInfo {
	if len(seedStr) == 0 {
		return []peer.AddrInfo{}
	}
	seeds := strings.Split(seedStr, ",")
	addrs := make([]peer.AddrInfo, 0, len(seeds))
	for _, s := range seeds {
		maddr, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			c.logger.Error("failed to parse seed node", "address", s, "error", err)
			continue
		}
		addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			c.logger.Error("failed to create addr info for seed", "address", maddr, "error", err)
			continue
		}
		addrs = append(addrs, *addrInfo)
	}
	return addrs
}

// getNamespace returns unique string identifying ORU network.
//
// It is used to advertise/find peers in libp2p DHT.
// For now, chainID is used.
func (c *Client) getNamespace() string {
	return c.chainID
}

func (c *Client) getTxTopic() string {
	return c.getNamespace() + txTopicSuffix
}

func (c *Client) getHeaderTopic() string {
	return c.getNamespace() + headerTopicSuffix
}
