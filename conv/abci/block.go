package abci

import (
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"

	"github.com/celestiaorg/optimint/types"
)

// ToABCIHeader converts Optimint header to Header format defined in ABCI.
// Caller should fill all the fields that are not available in Optimint header (like ChainID).
func ToABCIHeader(header *types.Header) (tmproto.Header, error) {
	hash := header.Hash()
	return tmproto.Header{
		Version: tmversion.Consensus{
			Block: header.Version.Block,
			App:   header.Version.App,
		},
		Height: int64(header.Height),
		Time:   time.Unix(int64(header.Time), 0),
		LastBlockId: tmproto.BlockID{
			Hash: hash[:],
			PartSetHeader: tmproto.PartSetHeader{
				Total: 0,
				Hash:  nil,
			},
		},
		LastCommitHash:     header.LastCommitHash[:],
		DataHash:           header.DataHash[:],
		ValidatorsHash:     nil,
		NextValidatorsHash: nil,
		ConsensusHash:      header.ConsensusHash[:],
		AppHash:            header.AppHash[:],
		LastResultsHash:    header.LastResultsHash[:],
		EvidenceHash:       nil,
		ProposerAddress:    header.ProposerAddress,
	}, nil
}

func ToABCIBlock(block *types.Block) (tmproto.Block, error) {
	abciHeader, err := ToABCIHeader(&block.Header)
	if err != nil {
		return tmproto.Block{}, err
	}
	hash := block.Hash()
	abciBlock := tmproto.Block{
		Header: abciHeader,
		Evidence: tmproto.EvidenceList{
			Evidence: nil,
		},
		LastCommit: &tmproto.Commit{
			Height: int64(block.LastCommit.Height),
			Round:  0,
			BlockID: tmproto.BlockID{
				Hash: hash[:],
				PartSetHeader: tmproto.PartSetHeader{
					Total: 0,
					Hash:  nil,
				},
			},
			Signatures: nil,
		},
	}
	abciBlock.Data.Txs = make([][]byte, len(block.Data.Txs))
	for i := range block.Data.Txs {
		abciBlock.Data.Txs[i] = block.Data.Txs[i]
	}

	return abciBlock, nil
}