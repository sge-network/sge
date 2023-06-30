package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group mint queries under a subcommand
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
		GetCmdQueryParams(),
		GetCmdQueryInflation(),
		GetCmdQueryPhaseStep(),
		GetCmdQueryPhaseProvisions(),
		GetCmdQueryEndPhaseStatus(),
	)

	return cmd
}
