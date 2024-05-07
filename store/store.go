package store

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"

	abci "github.com/cometbft/cometbft/abci/types"
	ds "github.com/ipfs/go-datastore"

	"github.com/celestiaorg/go-header"

	"github.com/rollkit/rollkit/types"
	pb "github.com/rollkit/rollkit/types/pb/rollkit"
)

var (
	blockPrefix          = "b"
	indexPrefix          = "i"
	commitPrefix         = "c"
	extendedCommitPrefix = "ec"
	statePrefix          = "s"
	responsesPrefix      = "r"
	metaPrefix           = "m"
)

// DefaultStore is a default store implmementation.
type DefaultStore struct {
	db     ds.TxnDatastore
	height atomic.Uint64
}

var _ Store = &DefaultStore{}

// New returns new, default store.
func New(ds ds.TxnDatastore) Store {
	return &DefaultStore{
		db: ds,
	}
}

// Close safely closes underlying data storage, to ensure that data is actually saved.
func (s *DefaultStore) Close() error {
	return s.db.Close()
}

// SetHeight sets the height saved in the Store if it is higher than the existing height
func (s *DefaultStore) SetHeight(ctx context.Context, height uint64) {
	for {
		storeHeight := s.height.Load()
		if height <= storeHeight {
			break
		}
		if s.height.CompareAndSwap(storeHeight, height) {
			break
		}
	}
}

// Height returns height of the highest block saved in the Store.
func (s *DefaultStore) Height() uint64 {
	return s.height.Load()
}

// SaveBlock adds block to the store along with corresponding commit.
// Stored height is updated if block height is greater than stored value.
func (s *DefaultStore) SaveBlock(ctx context.Context, block *types.Block, commit *types.Commit) error {
	hash := block.Hash()
	blockBlob, err := block.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal Block to binary: %w", err)
	}

	commitBlob, err := commit.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal Commit to binary: %w", err)
	}

	bb, err := s.db.NewTransaction(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to create a new batch for transaction: %w", err)
	}
	defer bb.Discard(ctx)

	err = bb.Put(ctx, ds.NewKey(getBlockKey(hash)), blockBlob)
	if err != nil {
		return fmt.Errorf("failed to create a new key for Block Blob: %w", err)
	}
	err = bb.Put(ctx, ds.NewKey(getCommitKey(hash)), commitBlob)
	if err != nil {
		return fmt.Errorf("failed to create a new key for Commit Blob: %w", err)
	}
	err = bb.Put(ctx, ds.NewKey(getIndexKey(block.Height())), hash[:])
	if err != nil {
		return fmt.Errorf("failed to create a new key using height of the block: %w", err)
	}

	if err = bb.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetBlock returns block at given height, or error if it's not found in Store.
// TODO(tzdybal): what is more common access pattern? by height or by hash?
// currently, we're indexing height->hash, and store blocks by hash, but we might as well store by height
// and index hash->height
func (s *DefaultStore) GetBlock(ctx context.Context, height uint64) (*types.Block, error) {
	h, err := s.loadHashFromIndex(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("failed to load hash from index: %w", err)
	}
	return s.GetBlockByHash(ctx, h)
}

// GetBlockByHash returns block with given block header hash, or error if it's not found in Store.
func (s *DefaultStore) GetBlockByHash(ctx context.Context, hash types.Hash) (*types.Block, error) {
	blockData, err := s.db.Get(ctx, ds.NewKey(getBlockKey(hash)))
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
func (s *DefaultStore) SaveBlockResponses(ctx context.Context, height uint64, responses *abci.ResponseFinalizeBlock) error {
	data, err := responses.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	return s.db.Put(ctx, ds.NewKey(getResponsesKey(height)), data)
}

// GetBlockResponses returns block results at given height, or error if it's not found in Store.
func (s *DefaultStore) GetBlockResponses(ctx context.Context, height uint64) (*abci.ResponseFinalizeBlock, error) {
	data, err := s.db.Get(ctx, ds.NewKey(getResponsesKey(height)))
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
func (s *DefaultStore) GetCommit(ctx context.Context, height uint64) (*types.Commit, error) {
	hash, err := s.loadHashFromIndex(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("failed to load hash from index: %w", err)
	}
	return s.GetCommitByHash(ctx, hash)
}

// GetCommitByHash returns commit for a block with given block header hash, or error if it's not found in Store.
func (s *DefaultStore) GetCommitByHash(ctx context.Context, hash types.Hash) (*types.Commit, error) {
	commitData, err := s.db.Get(ctx, ds.NewKey(getCommitKey(hash)))
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

// SaveExtendedCommit saves extended commit information in Store.
func (s *DefaultStore) SaveExtendedCommit(ctx context.Context, height uint64, commit *abci.ExtendedCommitInfo) error {
	bytes, err := commit.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal Extended Commit: %w", err)
	}
	return s.db.Put(ctx, ds.NewKey(getExtendedCommitKey(height)), bytes)
}

// GetExtendedCommit returns extended commit (commit with vote extensions) for a block at given height.
func (s *DefaultStore) GetExtendedCommit(ctx context.Context, height uint64) (*abci.ExtendedCommitInfo, error) {
	bytes, err := s.db.Get(ctx, ds.NewKey(getExtendedCommitKey(height)))
	if err != nil {
		return nil, fmt.Errorf("failed to load extended commit data: %w", err)
	}
	extendedCommit := new(abci.ExtendedCommitInfo)
	err = extendedCommit.Unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal extended commit: %w", err)
	}
	return extendedCommit, nil
}

// UpdateState updates state saved in Store. Only one State is stored.
// If there is no State in Store, state will be saved.
func (s *DefaultStore) UpdateState(ctx context.Context, state types.State) error {
	pbState, err := state.ToProto()
	if err != nil {
		return fmt.Errorf("failed to marshal state to JSON: %w", err)
	}
	data, err := pbState.Marshal()
	if err != nil {
		return err
	}
	return s.db.Put(ctx, ds.NewKey(getStateKey()), data)
}

// GetState returns last state saved with UpdateState.
func (s *DefaultStore) GetState(ctx context.Context) (types.State, error) {
	blob, err := s.db.Get(ctx, ds.NewKey(getStateKey()))
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
	return state, err
}

// SetMetadata saves arbitrary value in the store.
//
// Metadata is separated from other data by using prefix in KV.
func (s *DefaultStore) SetMetadata(ctx context.Context, key string, value []byte) error {
	err := s.db.Put(ctx, ds.NewKey(getMetaKey(key)), value)
	if err != nil {
		return fmt.Errorf("failed to set metadata for key '%s': %w", key, err)
	}
	return nil
}

// GetMetadata returns values stored for given key with SetMetadata.
func (s *DefaultStore) GetMetadata(ctx context.Context, key string) ([]byte, error) {
	data, err := s.db.Get(ctx, ds.NewKey(getMetaKey(key)))
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for key '%s': %w", key, err)
	}
	return data, nil
}

// loadHashFromIndex returns the hash of a block given its height
func (s *DefaultStore) loadHashFromIndex(ctx context.Context, height uint64) (header.Hash, error) {
	blob, err := s.db.Get(ctx, ds.NewKey(getIndexKey(height)))

	if err != nil {
		return nil, fmt.Errorf("failed to load block hash for height %v: %w", height, err)
	}
	if len(blob) != 32 {
		return nil, errors.New("invalid hash length")
	}
	return blob, nil
}

func getBlockKey(hash types.Hash) string {
	return GenerateKey([]string{blockPrefix, hex.EncodeToString(hash[:])})
}

func getCommitKey(hash types.Hash) string {
	return GenerateKey([]string{commitPrefix, hex.EncodeToString(hash[:])})
}

func getExtendedCommitKey(height uint64) string {
	return GenerateKey([]string{extendedCommitPrefix, strconv.FormatUint(height, 10)})
}

func getIndexKey(height uint64) string {
	return GenerateKey([]string{indexPrefix, strconv.FormatUint(height, 10)})
}

func getStateKey() string {
	return statePrefix
}

func getResponsesKey(height uint64) string {
	return GenerateKey([]string{responsesPrefix, strconv.FormatUint(height, 10)})
}

func getMetaKey(key string) string {
	return GenerateKey([]string{metaPrefix, key})
}
