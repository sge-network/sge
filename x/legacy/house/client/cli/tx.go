package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/sge-network/sge/x/legacy/house/types"
)

// NewTxCmd returns a root CLI command handler for all x/house transaction commands.
func NewTxCmd() *cobra.Command {
	houseTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "House transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	houseTxCmd.AddCommand(
		CmdDeposit(),
		CmdWithdraw(),
	)

	return houseTxCmd
}
