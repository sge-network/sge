package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/house/types"
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
		NewDepositCmd(),
	)

	return houseTxCmd
}

func NewDepositCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "deposit [sport_event_uid] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens to be part of a house",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit coins to be part of a house corresponding to a sport event.

				Example:
				$ %s tx house deposit bc79a72c-ad7e-4cf5-91a2-98af2751e812 1000usge --from mykey
				`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sportEventUid := args[0]

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			depAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgDeposit(depAddr, sportEventUid, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
