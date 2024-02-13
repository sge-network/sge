package cli

import (
	"github.com/spf13/cobra"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/sge-network/sge/x/bet/types"
)

// CmdPriceLockTopUp implements a command to top up the price lock module account.
func CmdPriceLockTopUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "TopUpPriceLockPool [funder_address] [amount]",
		Short: "Top up price lock pool",
		Long:  "Top up price lock pool module account. the funder and amount required.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			argFunder := args[0]
			argAmount := args[1]

			argAmountCosmosInt, ok := sdkmath.NewIntFromString(argAmount)
			if !ok {
				return types.ErrInvalidAmount
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPriceLockPoolTopUp(
				clientCtx.GetFromAddress().String(),
				argFunder,
				argAmountCosmosInt,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
