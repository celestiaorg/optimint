package types

import (
	abci "github.com/lazyledger/lazyledger-core/abci/types"
	pb "github.com/lazyledger/optimint/types/pb/optimint"
)

type Header struct {
	// Block and App version
	Version Version
	// NamespaceID identifies this chain e.g. when connected to other rollups via IBC.
	// TODO(ismail): figure out if we want to use namespace.ID here instead (downside is that it isn't fixed size)
	// at least extract the used constants (32, 8) as package variables though.
	NamespaceID [8]byte

	Height uint64
	Time   uint64 // time in tai64 format

	// prev block info
	LastHeaderHash [32]byte

	// hashes of block data
	LastCommitHash [32]byte // commit from aggregator(s) from the last block
	DataHash       [32]byte // Block.Data root aka Transactions
	ConsensusHash  [32]byte // consensus params for current block
	AppHash        [32]byte // state after applying txs from the current block

	// Root hash of all results from the txs from the previous block.
	// This is ABCI specific but smart-contract chains require some way of committing
	// to transaction receipts/results.
	LastResultsHash [32]byte

	// Note that the address can be derived from the pubkey which can be derived
	// from the signature when using secp256k.
	// We keep this in case users choose another signature format where the
	// pubkey can't be recovered by the signature (e.g. ed25519).
	ProposerAddress []byte // original proposer of the block
}

// Version captures the consensus rules for processing a block in the blockchain,
// including all blockchain data structures and the rules of the application's
// state transition machine.
// This is equivalent to the tmversion.Consensus type in Tendermint.
type Version struct {
	Block uint32
	App   uint32
}

type Block struct {
	Header     Header
	Data       Data
	LastCommit *Commit
}

type Data struct {
	Txs                    Txs
	IntermediateStateRoots IntermediateStateRoots
	Evidence               EvidenceData
}

type EvidenceData struct {
	Evidence []Evidence
}

type Commit struct {
	Height     uint64
	HeaderHash [32]byte
	Signatures []Signature // most of the time this is a single signature
}

type Signature []byte

type IntermediateStateRoots struct {
	RawRootsList [][]byte
}

func (h *Header) ToProto() *pb.Header {
	return &pb.Header{
		Version: &pb.Version{
			Block: h.Version.Block,
			App:   h.Version.App,
		},
		NamespaceId:     h.NamespaceID[:],
		Height:          h.Height,
		Time:            h.Time,
		LastHeaderHash:  h.LastHeaderHash[:],
		LastCommitHash:  h.LastCommitHash[:],
		DataHash:        h.DataHash[:],
		ConsensusHash:   h.ConsensusHash[:],
		AppHash:         h.AppHash[:],
		LastResultsHash: h.LastResultsHash[:],
		ProposerAddress: h.ProposerAddress[:],
	}
}

func (b *Block) ToProto() pb.Block {
	return pb.Block{
		Header: b.Header.ToProto(),
		Data: &pb.Data{
			Txs:                    txsToByteSlices(b.Data.Txs),
			IntermediateStateRoots: b.Data.IntermediateStateRoots.RawRootsList,
			Evidence:               evidenceToProto(b.Data.Evidence),
		},
		LastCommit: &pb.Commit{
			Height:     b.LastCommit.Height,
			HaderHash:  b.LastCommit.HeaderHash[:],
			Signatures: signaturesToByteSlices(b.LastCommit.Signatures),
		},
	}
}

func txsToByteSlices(txs Txs) [][]byte {
	bytes := make([][]byte, len(txs))
	for i := range txs {
		bytes[i] = txs[i]
	}
	return bytes
}

func evidenceToProto(evidence EvidenceData) []*abci.Evidence {
	var ret []*abci.Evidence
	for _, e := range evidence.Evidence {
		for _, ae := range e.ABCI() {
			ret = append(ret, &ae)
		}
	}
	return ret
}

func signaturesToByteSlices(sigs []Signature) [][]byte {
	bytes := make([][]byte, len(sigs))
	for i := range sigs {
		bytes[i] = sigs[i]
	}
	return bytes
}
