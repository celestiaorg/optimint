package store

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"

	abci "github.com/cometbft/cometbft/abci/types"
	ds "github.com/ipfs/go-datastore"
	"go.uber.org/multierr"

	"github.com/celestiaorg/go-header"

	"github.com/rollkit/rollkit/types"
	pb "github.com/rollkit/rollkit/types/pb/rollkit"
)

var (
	blockPrefix     = "b"
	indexPrefix     = "i"
	commitPrefix    = "c"
	statePrefix     = "s"
	responsesPrefix = "r"
)

// DefaultStore is a default store implmementation.
type DefaultStore struct {
	db ds.TxnDatastore

	height uint64
	ctx    context.Context
}

var _ Store = &DefaultStore{}

// New returns new, default store.
func New(ctx context.Context, ds ds.TxnDatastore) Store {
	return &DefaultStore{
		db:  ds,
		ctx: ctx,
	}
}

// SetHeight sets the height saved in the Store if it is higher than the existing height
func (s *DefaultStore) SetHeight(height uint64) {
	storeHeight := atomic.LoadUint64(&s.height)
	if height > storeHeight {
		_ = atomic.CompareAndSwapUint64(&s.height, storeHeight, height)
	}
}

// Height returns height of the highest block saved in the Store.
func (s *DefaultStore) Height() uint64 {
	return atomic.LoadUint64(&s.height)
}

// SaveBlock adds block to the store along with corresponding commit.
// Stored height is updated if block height is greater than stored value.
func (s *DefaultStore) SaveBlock(block *types.Block, commit *types.Commit) error {
	hash := block.Hash()
	blockBlob, err := block.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal Block to binary: %w", err)
	}

	commitBlob, err := commit.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal Commit to binary: %w", err)
	}

	bb, err := s.db.NewTransaction(s.ctx, false)
	if err != nil {
		return fmt.Errorf("failed to create a new batch for transaction: %w", err)
	}

	err = multierr.Append(err, bb.Put(s.ctx, ds.NewKey(getBlockKey(hash)), blockBlob))
	err = multierr.Append(err, bb.Put(s.ctx, ds.NewKey(getCommitKey(hash)), commitBlob))
	err = multierr.Append(err, bb.Put(s.ctx, ds.NewKey(getIndexKey(uint64(block.Height()))), hash[:]))

	if err != nil {
		bb.Discard(s.ctx)
		return err
	}

	if err = bb.Commit(s.ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetBlock returns block at given height, or error if it's not found in Store.
// TODO(tzdybal): what is more common access pattern? by height or by hash?
// currently, we're indexing height->hash, and store blocks by hash, but we might as well store by height
// and index hash->height
func (s *DefaultStore) GetBlock(height uint64) (*types.Block, error) {
	h, err := s.loadHashFromIndex(height)
	if err != nil {
		return nil, fmt.Errorf("failed to load hash from index: %w", err)
	}
	return s.GetBlockByHash(h)
}

// GetBlockByHash returns block with given block header hash, or error if it's not found in Store.
func (s *DefaultStore) GetBlockByHash(hash types.Hash) (*types.Block, error) {
	blockData, err := s.db.Get(s.ctx, ds.NewKey(getBlockKey(hash)))
	if err != nil {
		return nil, fmt.Errorf("failed to load block data: %w", err)
	}
	block := new(types.Block)
	err = block.UnmarshalBinary(blockData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal block data: %w", err)
	}

	return block, nil
}

// SaveBlockResponses saves block responses (events, tx responses, validator set updates, etc) in Store.
func (s *DefaultStore) SaveBlockResponses(height uint64, responses *abci.ResponseFinalizeBlock) error {
	data, err := responses.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	return s.db.Put(s.ctx, ds.NewKey(getResponsesKey(height)), data)
}

// GetBlockResponses returns block results at given height, or error if it's not found in Store.
func (s *DefaultStore) GetBlockResponses(height uint64) (*abci.ResponseFinalizeBlock, error) {
	data, err := s.db.Get(s.ctx, ds.NewKey(getResponsesKey(height)))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block results from height %v: %w", height, err)
	}
	var responses abci.ResponseFinalizeBlock
	err = responses.Unmarshal(data)
	if err != nil {
		return &responses, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &responses, nil
}

// GetCommit returns commit for a block at given height, or error if it's not found in Store.
func (s *DefaultStore) GetCommit(height uint64) (*types.Commit, error) {
	hash, err := s.loadHashFromIndex(height)
	if err != nil {
		return nil, fmt.Errorf("failed to load hash from index: %w", err)
	}
	return s.GetCommitByHash(hash)
}

// GetCommitByHash returns commit for a block with given block header hash, or error if it's not found in Store.
func (s *DefaultStore) GetCommitByHash(hash types.Hash) (*types.Commit, error) {
	commitData, err := s.db.Get(s.ctx, ds.NewKey(getCommitKey(hash)))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve commit from hash %v: %w", hash, err)
	}
	commit := new(types.Commit)
	err = commit.UnmarshalBinary(commitData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Commit into object: %w", err)
	}
	return commit, nil
}

// UpdateState updates state saved in Store. Only one State is stored.
// If there is no State in Store, state will be saved.
func (s *DefaultStore) UpdateState(state types.State) error {
	pbState, err := state.ToProto()
	if err != nil {
		return fmt.Errorf("failed to marshal state to JSON: %w", err)
	}
	data, err := pbState.Marshal()
	if err != nil {
		return err
	}
	return s.db.Put(s.ctx, ds.NewKey(getStateKey()), data)
}

// GetState returns last state saved with UpdateState.
func (s *DefaultStore) GetState() (types.State, error) {
	blob, err := s.db.Get(s.ctx, ds.NewKey(getStateKey()))
	if err != nil {
		return types.State{}, fmt.Errorf("failed to retrieve state: %w", err)
	}
	var pbState pb.State
	err = pbState.Unmarshal(blob)
	if err != nil {
		return types.State{}, fmt.Errorf("failed to unmarshal state from JSON: %w", err)
	}

	var state types.State
	err = state.FromProto(&pbState)
	atomic.StoreUint64(&s.height, uint64(state.LastBlockHeight))
	return state, err
}

// loadHashFromIndex returns the hash of a block given its height
func (s *DefaultStore) loadHashFromIndex(height uint64) (header.Hash, error) {
	blob, err := s.db.Get(s.ctx, ds.NewKey(getIndexKey(height)))

	if err != nil {
		return nil, fmt.Errorf("failed to load block hash for height %v: %w", height, err)
	}
	if len(blob) != 32 {
		return nil, errors.New("invalid hash length")
	}
	return blob, nil
}

func getBlockKey(hash types.Hash) string {
	return GenerateKey([]interface{}{blockPrefix, hex.EncodeToString(hash[:])})
}

func getCommitKey(hash types.Hash) string {
	return GenerateKey([]interface{}{commitPrefix, hex.EncodeToString(hash[:])})
}

func getIndexKey(height uint64) string {
	return GenerateKey([]interface{}{indexPrefix, height})
}

func getStateKey() string {
	return statePrefix
}

func getResponsesKey(height uint64) string {
	return GenerateKey([]interface{}{responsesPrefix, height})
}
