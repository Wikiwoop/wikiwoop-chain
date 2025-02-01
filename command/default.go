package command

import "github.com/Wikiwoop/wikiwoop-chain/server"

const (
	DefaultGenesisFileName = "genesis.json"
	DefaultChainName       = "wikiwoop-chain"
	DefaultChainID         = 514
	DefaultPremineBalance  = "0xD3C21BCECCEDA100000000000" 
	DefaultConsensus       = server.IBFTConsensus
	DefaultGenesisGasUsed  = 458752  // 0x70000
	DefaultGenesisGasLimit = 5242880 // 0x500000
)

const (
	JSONOutputFlag  = "json"
	GRPCAddressFlag = "grpc-address"
	JSONRPCFlag     = "jsonrpc"
)

// GRPCAddressFlagLEGACY Legacy flag that needs to be present to preserve backwards
// compatibility with running clients
const (
	GRPCAddressFlagLEGACY = "grpc"
)
