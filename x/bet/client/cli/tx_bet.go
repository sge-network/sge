package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cobra"
)

// CmdPlace implements a command to place and store a single bet
func CmdPlace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place [uid] [amount] [ticket]",
		Short: "Place bet",
		Long:  "Place bet uuid, amount and ticket required.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			uid := args[0]
			argAmount := args[1]
			argTicket := args[2]

			argAmountCosmosInt, ok := sdk.NewIntFromString(argAmount)
			if !ok {
				return types.ErrInvalidAmount
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlace(
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
