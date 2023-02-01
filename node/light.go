package node

import (
	"context"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/libp2p/go-libp2p/core/crypto"
	abciclient "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/celestiaorg/rollmint/config"
	"github.com/celestiaorg/rollmint/p2p"
)

var _ Node = &LightNode{}

type LightNode struct {
	service.BaseService

	P2P *p2p.Client

	app abciclient.Client

	ctx context.Context
}

func (ln *LightNode) GetClient() rpcclient.Client {
	return NewLightClient(ln)
}

func newLightNode(
	ctx context.Context,
	conf config.NodeConfig,
	p2pKey crypto.PrivKey,
	appClient abciclient.Client,
	genesis *tmtypes.GenesisDoc,
	logger log.Logger,
) (*LightNode, error) {
	client, err := p2p.NewClient(conf.P2P, p2pKey, genesis.ChainID, logger.With("module", "p2p"))
	if err != nil {
		return nil, err
	}

	node := &LightNode{
		P2P: client,
		app: appClient,
		ctx: ctx,
	}

	node.P2P.SetTxValidator(node.falseValidator())
	node.P2P.SetHeaderValidator(node.falseValidator())
	node.P2P.SetCommitValidator(node.falseValidator())
	node.P2P.SetFraudProofValidator(node.newFraudProofValidator())

	node.BaseService = *service.NewBaseService(logger, "LightNode", node)

	return node, nil
}

func (ln *LightNode) OnStart() error {
	if err := ln.P2P.Start(ln.ctx); err != nil {
		return err
	}

	return nil
}

func (ln *LightNode) falseValidator() p2p.GossipValidator {
	return func(*p2p.GossipMessage) bool {
		return false
	}
}

func (ln *LightNode) newFraudProofValidator() p2p.GossipValidator {
	return func(fraudProofMsg *p2p.GossipMessage) bool {
		ln.Logger.Info("fraud proof received", "from", fraudProofMsg.From, "bytes", len(fraudProofMsg.Data))
		var fraudProof abci.FraudProof
		err := fraudProof.Unmarshal(fraudProofMsg.Data)
		if err != nil {
			ln.Logger.Error("failed to deserialize fraud proof", "error", err)
			return false
		}

		resp, err := ln.app.VerifyFraudProofSync(abci.RequestVerifyFraudProof{
			FraudProof:           &fraudProof,
			ExpectedValidAppHash: fraudProof.ExpectedValidAppHash,
		})
		if err != nil {
			return false
		}

		if resp.Success {
			panic("received valid fraud proof! halting light client")
		}

		return false
	}
}

func newLightNode() (Node, error) {
	return &LightNode{}, nil
}
