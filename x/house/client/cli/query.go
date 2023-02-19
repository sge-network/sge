package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/house/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	houseQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the house module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	houseQueryCmd.AddCommand(
		GetCmdQueryDepsoits(),
		GetCmdQueryDepositorDepsoits(),
		GetCmdQueryDepositorWithdrawals(),
	)

	return houseQueryCmd
}
