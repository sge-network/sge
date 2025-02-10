package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/legacy/house/types"
)

func CmdDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [market_uid] [amount] [ticket]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens in a market order book to be the house",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit tokens in a market order book to be the house.

				Example:
				$ %s tx house deposit bc79a72c-ad7e-4cf5-91a2-98af2751e812 1000usge {ticket string} --from mykey
				`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argMarketUID := args[0]

			argAmountCosmosInt, ok := sdkmath.NewIntFromString(args[1])
			if !ok {
				return types.ErrInvalidAmount
			}

			argTicket := args[2]

			depAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgDeposit(depAddr.String(), argMarketUID, argAmountCosmosInt, argTicket)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
