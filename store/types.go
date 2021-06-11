package store

import "github.com/lazyledger/optimint/types"

// Store is minimal interface for storing and retrieving blocks and commits.
type Store interface {
	// Height returns height of the highest block in store.
	Height() uint64

	// SaveBlock saves block along with it's commit.
	SaveBlock(block *types.Block, commit *types.Commit) error

	// LoadBlock returns block at given height, or error if it's not found in Store.
	LoadBlock(height uint64) (*types.Block, error)
	// LoadBlockByHash returns block with given block header hash, or error if it's not found in Store.
	LoadBlockByHash(hash [32]byte) (*types.Block, error)

	// LoadCommit returns commit for a block at given height, or error if it's not found in Store.
	LoadCommit(height uint64) (*types.Commit, error)
	// LoadCommit returns commit for a block with given block header hash, or error if it's not found in Store.
	LoadCommitByHash(hash [32]byte) (*types.Commit, error)
}
