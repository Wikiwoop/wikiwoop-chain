package peers

import (
	"github.com/Wikiwoop/wikiwoop-chain/command/helper"
	"github.com/Wikiwoop/wikiwoop-chain/command/peers/add"
	"github.com/Wikiwoop/wikiwoop-chain/command/peers/list"
	"github.com/Wikiwoop/wikiwoop-chain/command/peers/status"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "Top level command for interacting with the network peers. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(peersCmd)

	registerSubcommands(peersCmd)

	return peersCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// peers status
		status.GetCommand(),
		// peers list
		list.GetCommand(),
		// peers add
		add.GetCommand(),
	)
}
