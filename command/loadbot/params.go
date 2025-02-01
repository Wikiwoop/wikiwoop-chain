package loadbot

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"github.com/Wikiwoop/wikiwoop-chain/command/loadbot/generator"
	"github.com/Wikiwoop/wikiwoop-chain/types"
	"github.com/umbracle/ethgo/abi"
)

var (
	params = &loadbotParams{}
)

var (
	errInvalidMode   = errors.New("invalid loadbot mode")
	errInvalidValues = errors.New("invalid values")
	errContractPath  = errors.New("contract path not specified")
)

const (
	// arbitrary value for total token supply
	// token has 5 decimals
	// transfers are done with 0.001 amount
	wrc20TokenSupply = "100000000000"
	wrc20TokenName   = "WooCoin"
	wrc20TokenSymbol = "WOC"

	wrc721TokenName   = "WocNFT"
	wrc721TokenSymbol = "WOOC"
)

const (
	tpsFlag      = "tps"
	modeFlag     = "mode"
	detailedFlag = "detailed"
	chainIDFlag  = "chain-id"
	senderFlag   = "sender"
	receiverFlag = "receiver"
	valueFlag    = "value"
	countFlag    = "count"
	maxConnsFlag = "max-conns"
	gasPriceFlag = "gas-price"
	gasLimitFlag = "gas-limit"
	contractFlag = "contract"
	maxWaitFlag  = "max-wait"
)

type loadbotParams struct {
	tps      uint64
	chainID  uint64
	count    uint64
	maxConns uint64
	maxWait  uint64

	contractPath string

	detailed bool

	modeRaw     string
	senderRaw   string
	receiverRaw string
	valueRaw    string
	gasPriceRaw string
	gasLimitRaw string

	mode             Mode
	sender           types.Address
	receiver         types.Address
	value            *big.Int
	gasPrice         *big.Int
	gasLimit         *big.Int
	contractArtifact *generator.ContractArtifact
	constructorArgs  []byte
}

func (p *loadbotParams) validateFlags() error {
	// check if valid mode is selected
	if err := p.isValidMode(); err != nil {
		return err
	}

	// validate the correct mode params
	if err := p.hasValidDeployParams(); err != nil {
		return err
	}

	return nil
}

func (p *loadbotParams) initRawParams() error {
	if err := p.initGasValues(); err != nil {
		return err
	}

	if err := p.initAddressValues(); err != nil {
		return err
	}

	if err := p.initTxnValue(); err != nil {
		return err
	}

	if err := p.initContract(); err != nil {
		return err
	}

	if err := p.initContractArtifactAndArgs(); err != nil {
		return err
	}

	return nil
}

func (p *loadbotParams) initGasValues() error {
	var parseErr error

	// Parse the gas price
	if p.gasPriceRaw != "" {
		if p.gasPrice, parseErr = types.ParseUint256orHex(&p.gasPriceRaw); parseErr != nil {
			return fmt.Errorf("failed to decode gas price to value: %w", parseErr)
		}
	}

	// Parse the gas limit
	if p.gasLimitRaw != "" {
		if p.gasLimit, parseErr = types.ParseUint256orHex(&p.gasLimitRaw); parseErr != nil {
			return fmt.Errorf("failed to decode gas limit to value: %w", parseErr)
		}
	}

	return nil
}

func (p *loadbotParams) initAddressValues() error {
	if err := p.sender.UnmarshalText([]byte(p.senderRaw)); err != nil {
		return fmt.Errorf("failed to decode sender address: %w", err)
	}

	if err := p.initReceiverAddress(); err != nil {
		return fmt.Errorf("failed to decode receiver address: %w", err)
	}

	return nil
}

func (p *loadbotParams) initReceiverAddress() error {
	if p.receiverRaw == "" {
		// No receiving address specified,
		// use the sender address as the receiving address
		p.receiver = p.sender

		return nil
	}

	return p.receiver.UnmarshalText([]byte(p.receiverRaw))
}

func (p *loadbotParams) initTxnValue() error {
	value, err := types.ParseUint256orHex(&p.valueRaw)
	if err != nil {
		return fmt.Errorf("failed to decode to value: %w", err)
	}

	p.value = value

	return nil
}

func (p *loadbotParams) initContract() error {
	var readErr error

	p.contractArtifact = &generator.ContractArtifact{
		Bytecode: generator.DefaultContractBytecode,
	}

	if p.contractPath != "" {
		if p.contractArtifact, readErr = generator.ReadContractArtifact(
			p.contractPath,
		); readErr != nil {
			return fmt.Errorf(
				"failed to read contract bytecode: %w",
				readErr,
			)
		}
	}

	return nil
}

func (p *loadbotParams) getRequiredFlags() []string {
	return []string{
		senderFlag,
	}
}

func (p *loadbotParams) generateConfig(
	jsonRPCAddress string,
	grpcAddress string,
) *Configuration {
	return &Configuration{
		TPS:              p.tps,
		Sender:           p.sender,
		Receiver:         p.receiver,
		Count:            p.count,
		Value:            p.value,
		JSONRPC:          jsonRPCAddress,
		GRPC:             grpcAddress,
		MaxConns:         int(p.maxConns),
		GeneratorMode:    p.mode,
		ChainID:          p.chainID,
		GasPrice:         p.gasPrice,
		GasLimit:         p.gasLimit,
		ContractArtifact: p.contractArtifact,
		ConstructorArgs:  p.constructorArgs,
		MaxWait:          p.maxWait,
	}
}

func (p *loadbotParams) isValidMode() error {
	// Set and validate the correct mode type
	p.mode = Mode(strings.ToLower(p.modeRaw))

	switch p.mode {
	case transfer, deploy, erc20, erc721:
		return nil

	default:
		return errInvalidMode
	}
}

func (p *loadbotParams) hasValidDeployParams() error {
	// fail if mode is deploy but we have no contract
	if p.mode == deploy && p.contractPath == "" {
		return errContractPath
	}

	return nil
}

func (p *loadbotParams) initContractArtifactAndArgs() error {
	var (
		ctrArtifact *generator.ContractArtifact
		ctrArgs     []byte
		err         error
	)

	switch p.mode {
	case wrc20:
		ctrArtifact = &generator.ContractArtifact{
			Bytecode: WRC20BIN,
			ABI:      abi.MustNewABI(ERC20ABI),
		}

		if ctrArgs, err = abi.Encode(
			[]string{wrc20TokenSupply, wrc20TokenName, wrc20TokenSymbol}, ctrArtifact.ABI.Constructor.Inputs); err != nil {
			return fmt.Errorf("failed to encode wrc20 constructor parameters: %w", err)
		}

	case wrc721:
		ctrArtifact = &generator.ContractArtifact{
			Bytecode: WRC721BIN,
			ABI:      abi.MustNewABI(WRC721ABI),
		}

		if ctrArgs, err = abi.Encode(
			[]string{Wrc721TokenName, wrc721TokenSymbol},
			ctrArtifact.ABI.Constructor.Inputs); err != nil {
			return fmt.Errorf("failed to encode wrc721 constructor parameters: %w", err)
		}

	default:
		ctrArtifact = &generator.ContractArtifact{
			Bytecode: generator.DefaultContractBytecode,
		}
		ctrArgs = nil
	}

	p.contractArtifact = ctrArtifact
	p.constructorArgs = ctrArgs

	return nil
}
