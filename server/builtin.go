package server

import (
	"github.com/Wikiwoop/wikiwoop-chain/consensus"
	consensusDev "github.com/Wikiwoop/wikiwoop-chain/consensus/dev"
	consensusDummy "github.com/Wikiwoop/wikiwoop-chain/consensus/dummy"
	consensusIBFT "github.com/Wikiwoop/wikiwoop-chain/consensus/ibft"
	"github.com/Wikiwoop/wikiwoop-chain/secrets"
	"github.com/Wikiwoop/wikiwoop-chain/secrets/awsssm"
	"github.com/Wikiwoop/wikiwoop-chain/secrets/gcpssm"
	"github.com/Wikiwoop/wikiwoop-chain/secrets/hashicorpvault"
	"github.com/Wikiwoop/wikiwoop-chain/secrets/local"
)

type ConsensusType string

const (
	DevConsensus   ConsensusType = "dev"
	IBFTConsensus  ConsensusType = "ibft"
	DummyConsensus ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:   consensusDev.Factory,
	IBFTConsensus:  consensusIBFT.Factory,
	DummyConsensus: consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
