package node

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/celestiaorg/go-header"
	goheaderp2p "github.com/celestiaorg/go-header/p2p"
	goheaderstore "github.com/celestiaorg/go-header/store"
	goheadersync "github.com/celestiaorg/go-header/sync"
	ds "github.com/ipfs/go-datastore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/conngater"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"go.uber.org/multierr"

	"github.com/rollkit/rollkit/config"
	"github.com/rollkit/rollkit/p2p"
	"github.com/rollkit/rollkit/types"
)

type BlockExchangeService struct {
	conf         config.NodeConfig
	genesis      *tmtypes.GenesisDoc
	p2p          *p2p.Client
	ex           *goheaderp2p.Exchange[*types.Block]
	syncer       *goheadersync.Syncer[*types.Block]
	sub          *goheaderp2p.Subscriber[*types.Block]
	p2pServer    *goheaderp2p.ExchangeServer[*types.Block]
	blockStore   *goheaderstore.Store[*types.Block]
	syncerStatus *SyncerStatus

	logger log.Logger
	ctx    context.Context
}

func NewBlockExchangeService(ctx context.Context, store ds.TxnDatastore, conf config.NodeConfig, genesis *tmtypes.GenesisDoc, p2p *p2p.Client, logger log.Logger) (*BlockExchangeService, error) {
	// store is TxnDatastore, but we require Batching, hence the type assertion
	// note, the badger datastore impl that is used in the background implements both
	storeBatch, ok := store.(ds.Batching)
	if !ok {
		return nil, errors.New("failed to access the datastore")
	}
	ss, err := goheaderstore.NewStore[*types.Block](storeBatch)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the block store: %w", err)
	}

	return &BlockExchangeService{
		conf:         conf,
		genesis:      genesis,
		p2p:          p2p,
		ctx:          ctx,
		blockStore:   ss,
		logger:       logger,
		syncerStatus: new(SyncerStatus),
	}, nil
}

func (bExService *BlockExchangeService) initBlockStoreAndStartSyncer(ctx context.Context, initial *types.Block) error {
	if err := bExService.blockStore.Init(ctx, initial); err != nil {
		return err
	}
	if err := bExService.syncer.Start(bExService.ctx); err != nil {
		return err
	}
	bExService.syncerStatus.m.Lock()
	defer bExService.syncerStatus.m.Unlock()
	bExService.syncerStatus.started = true
	return nil
}

func (bExService *BlockExchangeService) tryInitBlockStoreAndStartSyncer(ctx context.Context, trustedBlock *types.Block) {
	if trustedBlock != nil {
		if err := bExService.initBlockStoreAndStartSyncer(ctx, trustedBlock); err != nil {
			bExService.logger.Error("failed to initialize the headerstore and start syncer", "error", err)
		}
	}
}

func (bExService *BlockExchangeService) writeToBlockStoreAndBroadcast(ctx context.Context, block *types.Block) {
	// For genesis block initialize the store and start the syncer
	if block.Height() == bExService.genesis.InitialHeight {
		if err := bExService.blockStore.Init(ctx, block); err != nil {
			bExService.logger.Error("failed to initialize block store", "error", err)
		}

		if err := bExService.syncer.Start(bExService.ctx); err != nil {
			bExService.logger.Error("failed to start syncer after initializing block store", "error", err)
		}
	}

	// Broadcast for subscribers
	if err := bExService.sub.Broadcast(ctx, block); err != nil {
		bExService.logger.Error("failed to broadcast block", "error", err)
	}
}

// OnStart is a part of Service interface.
func (bExService *BlockExchangeService) Start() error {
	var err error
	// have to do the initializations here to utilize the p2p node which is created on start
	ps := bExService.p2p.PubSub()
	bExService.sub = goheaderp2p.NewSubscriber[*types.Block](ps, pubsub.DefaultMsgIdFn, bExService.genesis.ChainID)
	if err = bExService.sub.Start(bExService.ctx); err != nil {
		return fmt.Errorf("error while starting subscriber: %w", err)
	}
	if _, err := bExService.sub.Subscribe(); err != nil {
		return fmt.Errorf("error while subscribing: %w", err)
	}

	if err = bExService.blockStore.Start(bExService.ctx); err != nil {
		return fmt.Errorf("error while starting block store: %w", err)
	}

	_, _, network := bExService.p2p.Info()
	if bExService.p2pServer, err = newBlockP2PServer(bExService.p2p.Host(), bExService.blockStore, network); err != nil {
		return err
	}
	if err = bExService.p2pServer.Start(bExService.ctx); err != nil {
		return fmt.Errorf("error while starting p2p server: %w", err)
	}

	peerIDs := bExService.p2p.PeerIDs()
	if bExService.ex, err = newBlockP2PExchange(bExService.p2p.Host(), peerIDs, network, bExService.genesis.ChainID, bExService.p2p.ConnectionGater()); err != nil {
		return err
	}
	if err = bExService.ex.Start(bExService.ctx); err != nil {
		return fmt.Errorf("error while starting exchange: %w", err)
	}

	if bExService.syncer, err = newBlockSyncer(bExService.ex, bExService.blockStore, bExService.sub, goheadersync.WithBlockTime(bExService.conf.BlockTime)); err != nil {
		return err
	}

	// Check if the blockStore is not initialized and try initializing
	if bExService.blockStore.Height() > 0 {
		if err := bExService.syncer.Start(bExService.ctx); err != nil {
			return fmt.Errorf("error while starting the syncer: %w", err)
		}
		bExService.syncerStatus.started = true
		return nil
	}

	// Look to see if trusted hash is passed, if not get the genesis block
	var trustedBlock *types.Block
	// Try fetching the trusted block from peers if exists
	if len(peerIDs) > 0 {
		if bExService.conf.TrustedHash != "" {
			trustedHashBytes, err := hex.DecodeString(bExService.conf.TrustedHash)
			if err != nil {
				return fmt.Errorf("failed to parse the trusted hash for initializing the blockstore: %w", err)
			}

			if trustedBlock, err = bExService.ex.Get(bExService.ctx, header.Hash(trustedHashBytes)); err != nil {
				return fmt.Errorf("failed to fetch the trusted block for initializing the blockStore: %w", err)
			}
		} else {
			// Try fetching the genesis block if available, otherwise fallback to blocks
			if trustedBlock, err = bExService.ex.GetByHeight(bExService.ctx, uint64(bExService.genesis.InitialHeight)); err != nil {
				// Full/light nodes have to wait for aggregator to publish the genesis block
				// proposing aggregator can init the store and start the syncer when the first block is published
				bExService.logger.Info("failed to fetch the genesis block", "error", err)
			}
		}
	}
	go bExService.tryInitBlockStoreAndStartSyncer(bExService.ctx, trustedBlock)

	return nil
}

// OnStop is a part of Service interface.
func (bExService *BlockExchangeService) Stop() error {
	err := bExService.blockStore.Stop(bExService.ctx)
	err = multierr.Append(err, bExService.p2pServer.Stop(bExService.ctx))
	err = multierr.Append(err, bExService.ex.Stop(bExService.ctx))
	err = multierr.Append(err, bExService.sub.Stop(bExService.ctx))
	bExService.syncerStatus.m.Lock()
	defer bExService.syncerStatus.m.Unlock()
	if bExService.syncerStatus.started {
		err = multierr.Append(err, bExService.syncer.Stop(bExService.ctx))
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
