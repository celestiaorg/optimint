package block

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/cometbft/cometbft/libs/log"
	cmtypes "github.com/cometbft/cometbft/types"
	ds "github.com/ipfs/go-datastore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/conngater"
	"github.com/multiformats/go-multiaddr"

	"github.com/celestiaorg/go-header"
	goheaderp2p "github.com/celestiaorg/go-header/p2p"
	goheaderstore "github.com/celestiaorg/go-header/store"
	goheadersync "github.com/celestiaorg/go-header/sync"

	"github.com/rollkit/rollkit/config"
	"github.com/rollkit/rollkit/p2p"
	"github.com/rollkit/rollkit/types"
)

type syncType string

const (
	headerSync syncType = "headerSync"
	blockSync  syncType = "blockSync"
)

// SyncService is the P2P Sync Service for blocks and headers.
//
// Uses the go-header library for handling all P2P logic.
type SyncService[H header.Header[H]] struct {
	conf      config.NodeConfig
	genesis   *cmtypes.GenesisDoc
	p2p       *p2p.Client
	ex        *goheaderp2p.Exchange[H]
	sub       *goheaderp2p.Subscriber[H]
	p2pServer *goheaderp2p.ExchangeServer[H]
	store     *goheaderstore.Store[H]
	syncType  syncType

	syncer       *goheadersync.Syncer[H]
	syncerStatus *SyncerStatus

	logger log.Logger
	ctx    context.Context
}

// BlockSyncService is the P2P Sync Service for blocks.
type BlockSyncService = SyncService[*types.Block]

// HeaderSyncService is the P2P Sync Service for headers.
type HeaderSyncService = SyncService[*types.SignedHeader]

// NewBlockSyncService returns a new BlockSyncService.
func NewBlockSyncService(ctx context.Context, store ds.TxnDatastore, conf config.NodeConfig, genesis *cmtypes.GenesisDoc, p2p *p2p.Client, logger log.Logger) (*BlockSyncService, error) {
	return newSyncService[*types.Block](ctx, store, blockSync, conf, genesis, p2p, logger)
}

// NewHeaderSyncService returns a new HeaderSyncService.
func NewHeaderSyncService(ctx context.Context, store ds.TxnDatastore, conf config.NodeConfig, genesis *cmtypes.GenesisDoc, p2p *p2p.Client, logger log.Logger) (*HeaderSyncService, error) {
	return newSyncService[*types.SignedHeader](ctx, store, headerSync, conf, genesis, p2p, logger)
}

func newSyncService[H header.Header[H]](ctx context.Context, store ds.TxnDatastore, syncType syncType, conf config.NodeConfig, genesis *cmtypes.GenesisDoc, p2p *p2p.Client, logger log.Logger) (*SyncService[H], error) {
	if genesis == nil {
		return nil, errors.New("genesis doc cannot be nil")
	}
	if p2p == nil {
		return nil, errors.New("p2p client cannot be nil")
	}
	// store is TxnDatastore, but we require Batching, hence the type assertion
	// note, the badger datastore impl that is used in the background implements both
	storeBatch, ok := store.(ds.Batching)
	if !ok {
		return nil, errors.New("failed to access the datastore")
	}
	ss, err := goheaderstore.NewStore[H](
		storeBatch,
		goheaderstore.WithStorePrefix(string(syncType)),
		goheaderstore.WithMetrics(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the %s store: %w", syncType, err)
	}

	return &SyncService[H]{
		conf:         conf,
		genesis:      genesis,
		p2p:          p2p,
		ctx:          ctx,
		store:        ss,
		syncType:     syncType,
		logger:       logger,
		syncerStatus: new(SyncerStatus),
	}, nil
}

// Store returns the store of the SyncService
func (syncService *SyncService[H]) Store() *goheaderstore.Store[H] {
	return syncService.store
}

func (syncService *SyncService[H]) initStoreAndStartSyncer(ctx context.Context, initial H) error {
	if initial.IsZero() {
		return fmt.Errorf("failed to initialize the store and start syncer")
	}
	if err := syncService.store.Init(ctx, initial); err != nil {
		return err
	}
	if err := syncService.StartSyncer(); err != nil {
		return err
	}
	return nil
}

// WriteToStoreAndBroadcast initializes store if needed and broadcasts  provided header or block.
// Note: Only returns an error in case store can't be initialized. Logs error if there's one while broadcasting.
func (syncService *SyncService[H]) WriteToStoreAndBroadcast(ctx context.Context, headerOrBlock H) error {
	if syncService.genesis.InitialHeight < 0 {
		return fmt.Errorf("invalid initial height; cannot be negative")
	}
	isGenesis := headerOrBlock.Height() == uint64(syncService.genesis.InitialHeight)
	// For genesis header/block initialize the store and start the syncer
	if isGenesis {
		if err := syncService.store.Init(ctx, headerOrBlock); err != nil {
			return fmt.Errorf("failed to initialize the store")
		}

		if err := syncService.StartSyncer(); err != nil {
			return fmt.Errorf("failed to start syncer after initializing the store")
		}
	}

	// Broadcast for subscribers
	if err := syncService.sub.Broadcast(ctx, headerOrBlock); err != nil {
		// for the genesis header, broadcast error is expected as we have already initialized the store
		// for starting the syncer. Hence, we ignore the error.
		// exact reason: validation failed, err header verification failed: known header: '1' <= current '1'
		if isGenesis && errors.Is(err, pubsub.ValidationError{Reason: pubsub.RejectValidationFailed}) {
			return nil
		}
		return fmt.Errorf("failed to broadcast: %w", err)
	}
	return nil
}

func (syncService *SyncService[H]) isInitialized() bool {
	return syncService.store.Height() > 0
}

// Start is a part of Service interface.
func (syncService *SyncService[H]) Start() error {
	// have to do the initializations here to utilize the p2p node which is created on start
	ps := syncService.p2p.PubSub()
	chainID := syncService.genesis.ChainID + "-" + string(syncService.syncType)

	var err error
	syncService.sub, err = goheaderp2p.NewSubscriber[H](
		ps,
		pubsub.DefaultMsgIdFn,
		goheaderp2p.WithSubscriberNetworkID(chainID),
		goheaderp2p.WithSubscriberMetrics(),
	)
	if err != nil {
		return err
	}

	if err := syncService.sub.Start(syncService.ctx); err != nil {
		return fmt.Errorf("error while starting subscriber: %w", err)
	}
	if _, err := syncService.sub.Subscribe(); err != nil {
		return fmt.Errorf("error while subscribing: %w", err)
	}

	if err := syncService.store.Start(syncService.ctx); err != nil {
		return fmt.Errorf("error while starting store: %w", err)
	}

	_, _, network, err := syncService.p2p.Info()
	if err != nil {
		return fmt.Errorf("error while fetching the network: %w", err)
	}
	networkID := network + "-" + string(syncService.syncType)

	if syncService.p2pServer, err = newP2PServer(syncService.p2p.Host(), syncService.store, networkID); err != nil {
		return fmt.Errorf("error while creating p2p server: %w", err)
	}
	if err := syncService.p2pServer.Start(syncService.ctx); err != nil {
		return fmt.Errorf("error while starting p2p server: %w", err)
	}

	peerIDs := syncService.p2p.PeerIDs()
	if !syncService.conf.Aggregator {
		peerIDs = append(peerIDs, getSeedNodes(syncService.conf.P2P.Seeds, syncService.logger)...)
	}
	if syncService.ex, err = newP2PExchange[H](syncService.p2p.Host(), peerIDs, networkID, syncService.genesis.ChainID, syncService.p2p.ConnectionGater()); err != nil {
		return fmt.Errorf("error while creating exchange: %w", err)
	}
	if err := syncService.ex.Start(syncService.ctx); err != nil {
		return fmt.Errorf("error while starting exchange: %w", err)
	}

	if syncService.syncer, err = newSyncer[H](
		syncService.ex,
		syncService.store,
		syncService.sub,
		[]goheadersync.Option{goheadersync.WithBlockTime(syncService.conf.BlockTime)},
	); err != nil {
		return fmt.Errorf("error while creating syncer: %w", err)
	}

	if syncService.isInitialized() {
		if err := syncService.StartSyncer(); err != nil {
			return fmt.Errorf("error while starting the syncer: %w", err)
		}
		return nil
	}

	// Look to see if trusted hash is passed, if not get the genesis header/block
	var trusted H
	// Try fetching the trusted header/block from peers if exists
	if len(peerIDs) > 0 {
		if syncService.conf.TrustedHash != "" {
			trustedHashBytes, err := hex.DecodeString(syncService.conf.TrustedHash)
			if err != nil {
				return fmt.Errorf("failed to parse the trusted hash for initializing the store: %w", err)
			}

			if trusted, err = syncService.ex.Get(syncService.ctx, header.Hash(trustedHashBytes)); err != nil {
				return fmt.Errorf("failed to fetch the trusted header/block for initializing the store: %w", err)
			}
		} else {
			// Try fetching the genesis header/block if available, otherwise fallback to block
			if trusted, err = syncService.ex.GetByHeight(syncService.ctx, uint64(syncService.genesis.InitialHeight)); err != nil {
				// Full/light nodes have to wait for aggregator to publish the genesis block
				// proposing aggregator can init the store and start the syncer when the first block is published
				return fmt.Errorf("failed to fetch the genesis: %w", err)
			}
		}
		return syncService.initStoreAndStartSyncer(syncService.ctx, trusted)
	}
	return nil
}

// Stop is a part of Service interface.
//
// `store` is closed last because it's used by other services.
func (syncService *SyncService[H]) Stop() error {
	err := errors.Join(
		syncService.p2pServer.Stop(syncService.ctx),
		syncService.ex.Stop(syncService.ctx),
		syncService.sub.Stop(syncService.ctx),
	)
	if syncService.syncerStatus.isStarted() {
		err = errors.Join(err, syncService.syncer.Stop(syncService.ctx))
	}
	err = errors.Join(err, syncService.store.Stop(syncService.ctx))
	return err
}

// newP2PServer constructs a new ExchangeServer using the given Network as a protocolID suffix.
func newP2PServer[H header.Header[H]](
	host host.Host,
	store *goheaderstore.Store[H],
	network string,
	opts ...goheaderp2p.Option[goheaderp2p.ServerParameters],
) (*goheaderp2p.ExchangeServer[H], error) {
	opts = append(opts,
		goheaderp2p.WithNetworkID[goheaderp2p.ServerParameters](network),
		goheaderp2p.WithMetrics[goheaderp2p.ServerParameters](),
	)
	return goheaderp2p.NewExchangeServer[H](host, store, opts...)
}

func newP2PExchange[H header.Header[H]](
	host host.Host,
	peers []peer.ID,
	network, chainID string,
	conngater *conngater.BasicConnectionGater,
	opts ...goheaderp2p.Option[goheaderp2p.ClientParameters],
) (*goheaderp2p.Exchange[H], error) {
	opts = append(opts,
		goheaderp2p.WithNetworkID[goheaderp2p.ClientParameters](network),
		goheaderp2p.WithChainID(chainID),
		goheaderp2p.WithMetrics[goheaderp2p.ClientParameters](),
	)
	return goheaderp2p.NewExchange[H](host, peers, conngater, opts...)
}

// newSyncer constructs new Syncer for headers/blocks.
func newSyncer[H header.Header[H]](
	ex header.Exchange[H],
	store header.Store[H],
	sub header.Subscriber[H],
	opts []goheadersync.Option,
) (*goheadersync.Syncer[H], error) {
	opts = append(opts,
		goheadersync.WithMetrics(),
	)
	return goheadersync.NewSyncer[H](ex, store, sub, opts...)
}

// StartSyncer starts the SyncService's syncer
func (syncService *SyncService[H]) StartSyncer() error {
	if syncService.syncerStatus.isStarted() {
		return nil
	}
	err := syncService.syncer.Start(syncService.ctx)
	if err != nil {
		return err
	}
	syncService.syncerStatus.started.Store(true)
	return nil
}

func getSeedNodes(seeds string, logger log.Logger) []peer.ID {
	var peerIDs []peer.ID
	for _, seed := range strings.Split(seeds, ",") {
		maddr, err := multiaddr.NewMultiaddr(seed)
		if err != nil {
			logger.Error("failed to parse peer", "address", seed, "error", err)
			continue
		}
		addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			logger.Error("failed to create addr info for peer", "address", maddr, "error", err)
			continue
		}
		peerIDs = append(peerIDs, addrInfo.ID)
	}
	return peerIDs
}
