package store

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/rollkit/rollkit/types"
)

// Store is minimal interface for storing and retrieving blocks, commits and state.
type Store interface {
	// Height returns height of the highest block in store.
	Height() uint64

	// SetHeight sets the height saved in the Store if it is higher than the existing height.
	SetHeight(ctx context.Context, height uint64)

	// SaveBlock saves block along with its seen commit (which will be included in the next block).
	SaveBlock(ctx context.Context, block *types.Block, commit *types.Commit) error

	// GetBlock returns block at given height, or error if it's not found in Store.
	GetBlock(ctx context.Context, height uint64) (*types.Block, error)
	// GetBlockByHash returns block with given block header hash, or error if it's not found in Store.
	GetBlockByHash(ctx context.Context, hash types.Hash) (*types.Block, error)

	// SaveBlockResponses saves block responses (events, tx responses, validator set updates, etc) in Store.
	SaveBlockResponses(ctx context.Context, height uint64, responses *abci.ResponseFinalizeBlock) error

	// GetBlockResponses returns block results at given height, or error if it's not found in Store.
	GetBlockResponses(ctx context.Context, height uint64) (*abci.ResponseFinalizeBlock, error)

	// GetCommit returns commit for a block at given height, or error if it's not found in Store.
	GetCommit(ctx context.Context, height uint64) (*types.Commit, error)
	// GetCommitByHash returns commit for a block with given block header hash, or error if it's not found in Store.
	GetCommitByHash(ctx context.Context, hash types.Hash) (*types.Commit, error)

	// SaveExtendedCommit saves extended commit information in Store.
	SaveExtendedCommit(ctx context.Context, height uint64, commit *abci.ExtendedCommitInfo) error

	// GetExtendedCommit returns extended commit (commit with vote extensions) for a block at given height.
	GetExtendedCommit(ctx context.Context, height uint64) (*abci.ExtendedCommitInfo, error)

	// UpdateState updates state saved in Store. Only one State is stored.
	// If there is no State in Store, state will be saved.
	UpdateState(ctx context.Context, state types.State) error
	// GetState returns last state saved with UpdateState.
	GetState(ctx context.Context) (types.State, error)

	// SetMetadata saves arbitrary value in the store.
	//
	// This method enables rollkit to safely persist any information.
	SetMetadata(ctx context.Context, key string, value []byte) error

	// GetMetadata returns values stored for given key with SetMetadata.
	GetMetadata(ctx context.Context, key string) ([]byte, error)

	// Close safely closes underlying data storage, to ensure that data is actually saved.
	Close() error
}
