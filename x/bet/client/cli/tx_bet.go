package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/bet/types"
)

// CmdWager implements a command to place and store a single bet
func CmdWager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wager [uid] [amount] [ticket]",
		Short: "Wager on an odds",
		Long:  "Wager on an odds. the uuid, amount and ticket required.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			uid := args[0]
			argAmount := args[1]
			argTicket := args[2]

			argAmountCosmosInt, ok := sdkmath.NewIntFromString(argAmount)
			if !ok {
				return types.ErrInvalidAmount
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWager(
				clientCtx.GetFromAddress().String(),
				types.WagerProps{
					UID:    uid,
					Amount: argAmountCosmosInt,
					Ticket: argTicket,
				},
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
