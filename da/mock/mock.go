package mock

import (
	"github.com/celestiaorg/optimint/da"
	"github.com/celestiaorg/optimint/log"
	"github.com/celestiaorg/optimint/store"
	"github.com/celestiaorg/optimint/types"
)

// MockDataAvailabilityLayerClient is intended only for usage in tests.
// It does actually ensures DA - it stores data in-memory.
type MockDataAvailabilityLayerClient struct {
	logger log.Logger

	Blocks     map[[32]byte]*types.Block
	BlockIndex map[uint64][32]byte

	dalcKV *store.PrefixKV
}

var _ da.DataAvailabilityLayerClient = &MockDataAvailabilityLayerClient{}
var _ da.BlockRetriever = &MockDataAvailabilityLayerClient{}

// Init is called once to allow DA client to read configuration and initialize resources.
func (m *MockDataAvailabilityLayerClient) Init(config []byte, dalcKV *store.PrefixKV, logger log.Logger) error {
	m.logger = logger
	m.Blocks = make(map[[32]byte]*types.Block)
	m.BlockIndex = make(map[uint64][32]byte)
	m.dalcKV = dalcKV
	return nil
}

// Start implements DataAvailabilityLayerClient interface.
func (m *MockDataAvailabilityLayerClient) Start() error {
	m.logger.Debug("Mock Data Availability Layer Client starting")
	return nil
}

// Stop implements DataAvailabilityLayerClient interface.
func (m *MockDataAvailabilityLayerClient) Stop() error {
	m.logger.Debug("Mock Data Availability Layer Client stopped")
	return nil
}

// TODO(jbowen93): This is where storage to the kvStore needs to take place

// SubmitBlock submits the passed in block to the DA layer.
// This should create a transaction which (potentially)
// triggers a state transition in the DA layer.
func (m *MockDataAvailabilityLayerClient) SubmitBlock(block *types.Block) da.ResultSubmitBlock {
	m.logger.Debug("Submitting block to DA layer!", "height", block.Header.Height)

	hash := block.Header.Hash()
	m.Blocks[hash] = block
	m.BlockIndex[block.Header.Height] = hash

	return da.ResultSubmitBlock{
		DAResult: da.DAResult{
			Code:    da.StatusSuccess,
			Message: "OK",
		},
	}
}

// CheckBlockAvailability queries DA layer to check data availability of block corresponding to given header.
func (m *MockDataAvailabilityLayerClient) CheckBlockAvailability(header *types.Header) da.ResultCheckBlock {
	_, ok := m.Blocks[header.Hash()]
	return da.ResultCheckBlock{DAResult: da.DAResult{Code: da.StatusSuccess}, DataAvailable: ok}
}

// TODO(jbowen93): This is where retrieval from the kvStore needs to take place

// RetrieveBlock returns block at given height from data availability layer.
func (m *MockDataAvailabilityLayerClient) RetrieveBlock(height uint64) da.ResultRetrieveBlock {
	hash := m.BlockIndex[height]
	return da.ResultRetrieveBlock{DAResult: da.DAResult{Code: da.StatusSuccess}, Block: m.Blocks[hash]}
}
