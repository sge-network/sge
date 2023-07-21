package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/house/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName,
		Short: fmt.Sprintf(
			"Querying commands for the %s module",
			types.ModuleName,
		),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		GetCmdQueryDeposits(),
		GetCmdQueryDepositsByAccount(),
		GetCmdQueryWithdrawal(),
		GetCmdQueryWithdrawalsByAccount(),
	)

	return cmd
}
