package root

import (
	"fmt"
	"os"
	"github.com/Wikiwoop/wikiwoop-chain/command/backup"
	"github.com/Wikiwoop/wikiwoop-chain/command/genesis"
	"github.com/Wikiwoop/wikiwoop-chain/command/helper"
	"github.com/Wikiwoop/wikiwoop-chain/command/ibft"
	"github.com/Wikiwoop/wikiwoop-chain/command/license"
	"github.com/Wikiwoop/wikiwoop-chain/command/loadbot"
	"github.com/Wikiwoop/wikiwoop-chain/command/monitor"
	"github.com/Wikiwoop/wikiwoop-chain/command/peers"
	"github.com/Wikiwoop/wikiwoop-chain/command/secrets"
	"github.com/Wikiwoop/wikiwoop-chain/command/server"
	"github.com/Wikiwoop/wikiwoop-chain/command/status"
	"github.com/Wikiwoop/wikiwoop-chain/command/txpool"
	"github.com/Wikiwoop/wikiwoop-chain/command/version"
	"github.com/Wikiwoop/wikiwoop-chain/command/whitelist"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Wikiwoop Chain is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		monitor.GetCommand(),
		loadbot.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		whitelist.GetCommand(),
		license.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
