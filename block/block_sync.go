package block

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/celestiaorg/go-header"
	goheaderp2p "github.com/celestiaorg/go-header/p2p"
	goheaderstore "github.com/celestiaorg/go-header/store"
	goheadersync "github.com/celestiaorg/go-header/sync"
	"github.com/cometbft/cometbft/libs/log"
	cmtypes "github.com/cometbft/cometbft/types"
	ds "github.com/ipfs/go-datastore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/conngater"
	"go.uber.org/multierr"

	"github.com/rollkit/rollkit/config"
	"github.com/rollkit/rollkit/p2p"
	"github.com/rollkit/rollkit/types"
)

// P2P Sync Service for block that implements the go-header interface.
// Contains a block store where synced blocks are stored.
// Uses the go-header library for handling all P2P logic.
type BlockSyncService struct {
	conf       config.NodeConfig
	genesis    *cmtypes.GenesisDoc
	p2p        *p2p.Client
	ex         *goheaderp2p.Exchange[*types.Block]
	sub        *goheaderp2p.Subscriber[*types.Block]
	p2pServer  *goheaderp2p.ExchangeServer[*types.Block]
	blockStore *goheaderstore.Store[*types.Block]

	syncer       *goheadersync.Syncer[*types.Block]
	syncerStatus *SyncerStatus

	logger log.Logger
	ctx    context.Context
}

func NewBlockSyncService(ctx context.Context, store ds.TxnDatastore, conf config.NodeConfig, genesis *cmtypes.GenesisDoc, p2p *p2p.Client, logger log.Logger) (*BlockSyncService, error) {
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
	ss, err := goheaderstore.NewStore[*types.Block](storeBatch, goheaderstore.WithStorePrefix("blockSync"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the block store: %w", err)
	}

	return &BlockSyncService{
		conf:         conf,
		genesis:      genesis,
		p2p:          p2p,
		ctx:          ctx,
		blockStore:   ss,
		logger:       logger,
		syncerStatus: new(SyncerStatus),
	}, nil
}

// BlockStore returns the blockstore of the BlockSyncService
func (bSyncService *BlockSyncService) BlockStore() *goheaderstore.Store[*types.Block] {
	return bSyncService.blockStore
}

func (bSyncService *BlockSyncService) initBlockStoreAndStartSyncer(ctx context.Context, initial *types.Block) error {
	if initial == nil {
		return fmt.Errorf("failed to initialize the blockstore and start syncer")
	}
	if err := bSyncService.blockStore.Init(ctx, initial); err != nil {
		return err
	}
	if err := bSyncService.StartSyncer(); err != nil {
		return err
	}
	return nil
}

// Initialize block store if needed and broadcasts provided block.
// Note: Only returns an error in case block store can't be initialized. Logs error if there's one while broadcasting.
func (bSyncService *BlockSyncService) WriteToBlockStoreAndBroadcast(ctx context.Context, block *types.Block) error {
	// For genesis block initialize the store and start the syncer
	if int64(block.Height()) == bSyncService.genesis.InitialHeight {
		if err := bSyncService.blockStore.Init(ctx, block); err != nil {
			return fmt.Errorf("failed to initialize block store")
		}

		if err := bSyncService.StartSyncer(); err != nil {
			return fmt.Errorf("failed to start syncer after initializing block store")
		}
	}

	// Broadcast for subscribers
	if err := bSyncService.sub.Broadcast(ctx, block); err != nil {
		bSyncService.logger.Error("failed to broadcast block", "error", err)
	}
	return nil
}

func (bSyncService *BlockSyncService) isInitialized() bool {
	return bSyncService.blockStore.Height() > 0
}

// OnStart is a part of Service interface.
func (bSyncService *BlockSyncService) Start() error {
	// have to do the initializations here to utilize the p2p node which is created on start
	ps := bSyncService.p2p.PubSub()
	chainIDBlock := bSyncService.genesis.ChainID + "-block"
	bSyncService.sub = goheaderp2p.NewSubscriber[*types.Block](ps, pubsub.DefaultMsgIdFn, chainIDBlock)
	if err := bSyncService.sub.Start(bSyncService.ctx); err != nil {
		return fmt.Errorf("error while starting subscriber: %w", err)
	}
	if _, err := bSyncService.sub.Subscribe(); err != nil {
		return fmt.Errorf("error while subscribing: %w", err)
	}

	if err := bSyncService.blockStore.Start(bSyncService.ctx); err != nil {
		return fmt.Errorf("error while starting block store: %w", err)
	}

	var err error
	_, _, network, err := bSyncService.p2p.Info()
	if err != nil {
		return fmt.Errorf("error while fetching the network: %w", err)
	}
	networkIDBlock := network + "-block"

	if bSyncService.p2pServer, err = newBlockP2PServer(bSyncService.p2p.Host(), bSyncService.blockStore, networkIDBlock); err != nil {
		return fmt.Errorf("error while creating p2p server: %w", err)
	}
	if err := bSyncService.p2pServer.Start(bSyncService.ctx); err != nil {
		return fmt.Errorf("error while starting p2p server: %w", err)
	}

	peerIDs := bSyncService.p2p.PeerIDs()
	if bSyncService.ex, err = newBlockP2PExchange(bSyncService.p2p.Host(), peerIDs, networkIDBlock, chainIDBlock, bSyncService.p2p.ConnectionGater()); err != nil {
		return fmt.Errorf("error while creating exchange: %w", err)
	}
	if err := bSyncService.ex.Start(bSyncService.ctx); err != nil {
		return fmt.Errorf("error while starting exchange: %w", err)
	}

	if bSyncService.syncer, err = newBlockSyncer(bSyncService.ex, bSyncService.blockStore, bSyncService.sub, goheadersync.WithBlockTime(bSyncService.conf.BlockTime)); err != nil {
		return fmt.Errorf("error while creating syncer: %w", err)
	}

	if bSyncService.isInitialized() {
		if err := bSyncService.StartSyncer(); err != nil {
			return fmt.Errorf("error while starting the syncer: %w", err)
		}
		return nil
	}

	// Look to see if trusted hash is passed, if not get the genesis block
	var trustedBlock *types.Block
	// Try fetching the trusted block from peers if exists
	if len(peerIDs) > 0 {
		if bSyncService.conf.TrustedHash != "" {
			trustedHashBytes, err := hex.DecodeString(bSyncService.conf.TrustedHash)
			if err != nil {
				return fmt.Errorf("failed to parse the trusted hash for initializing the blockstore: %w", err)
			}

			if trustedBlock, err = bSyncService.ex.Get(bSyncService.ctx, header.Hash(trustedHashBytes)); err != nil {
				return fmt.Errorf("failed to fetch the trusted block for initializing the blockStore: %w", err)
			}
		} else {
			// Try fetching the genesis block if available, otherwise fallback to blocks
			if trustedBlock, err = bSyncService.ex.GetByHeight(bSyncService.ctx, uint64(bSyncService.genesis.InitialHeight)); err != nil {
				// Full/light nodes have to wait for aggregator to publish the genesis block
				// proposing aggregator can init the store and start the syncer when the first block is published
				return fmt.Errorf("failed to fetch the genesis block: %w", err)
			}
		}
		return bSyncService.initBlockStoreAndStartSyncer(bSyncService.ctx, trustedBlock)
	}
	return nil
}

// OnStop is a part of Service interface.
func (bSyncService *BlockSyncService) Stop() error {
	err := bSyncService.blockStore.Stop(bSyncService.ctx)
	err = multierr.Append(err, bSyncService.p2pServer.Stop(bSyncService.ctx))
	err = multierr.Append(err, bSyncService.ex.Stop(bSyncService.ctx))
	err = multierr.Append(err, bSyncService.sub.Stop(bSyncService.ctx))
	if bSyncService.syncerStatus.isStarted() {
		err = multierr.Append(err, bSyncService.syncer.Stop(bSyncService.ctx))
	}
	return err
}

// newBlockP2PServer constructs a new ExchangeServer using the given Network as a protocolID suffix.
func newBlockP2PServer(
	host host.Host,
	store *goheaderstore.Store[*types.Block],
	network string,
	opts ...goheaderp2p.Option[goheaderp2p.ServerParameters],
) (*goheaderp2p.ExchangeServer[*types.Block], error) {
	opts = append(opts,
		goheaderp2p.WithNetworkID[goheaderp2p.ServerParameters](network),
	)
	return goheaderp2p.NewExchangeServer[*types.Block](host, store, opts...)
}

func newBlockP2PExchange(
	host host.Host,
	peers []peer.ID,
	network, chainID string,
	conngater *conngater.BasicConnectionGater,
	opts ...goheaderp2p.Option[goheaderp2p.ClientParameters],
) (*goheaderp2p.Exchange[*types.Block], error) {
	opts = append(opts,
		goheaderp2p.WithNetworkID[goheaderp2p.ClientParameters](network),
		goheaderp2p.WithChainID(chainID),
	)
	return goheaderp2p.NewExchange[*types.Block](host, peers, conngater, opts...)
}

// newBlockSyncer constructs new Syncer for blocks.
func newBlockSyncer(
	ex header.Exchange[*types.Block],
	store header.Store[*types.Block],
	sub header.Subscriber[*types.Block],
	opt goheadersync.Options,
) (*goheadersync.Syncer[*types.Block], error) {
	return goheadersync.NewSyncer[*types.Block](ex, store, sub, opt)
}

func (bSyncService *BlockSyncService) StartSyncer() error {
	bSyncService.syncerStatus.m.Lock()
	defer bSyncService.syncerStatus.m.Unlock()
	if bSyncService.syncerStatus.started {
		return nil
	}
	err := bSyncService.syncer.Start(bSyncService.ctx)
	if err != nil {
		return err
	}
	bSyncService.syncerStatus.started = true
	return nil
}
