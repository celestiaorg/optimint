module github.com/celestiaorg/optimint

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.40.0-rc5
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/ipfs/go-log v1.0.4
	github.com/lazyledger/lazyledger-app v0.0.0-20210615142030-0e60861ff7f5
	github.com/libp2p/go-libp2p v0.13.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-discovery v0.5.0
	github.com/libp2p/go-libp2p-kad-dht v0.11.1
	github.com/libp2p/go-libp2p-pubsub v0.4.1
	github.com/minio/sha256-simd v0.1.1
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/pelletier/go-toml v1.9.3
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.13
	go.uber.org/multierr v1.6.0
	google.golang.org/grpc v1.39.0
)

replace (
	github.com/cosmos/cosmos-sdk v0.40.0-rc5 => github.com/celestiaorg/cosmos-sdk v0.40.0-rc5.0.20210121152417-3addd7f65d1c
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
