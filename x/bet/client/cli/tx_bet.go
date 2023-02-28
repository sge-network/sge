package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// CmdPlaceBet implements a command to place and store a single bet
func CmdPlaceBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bet [uid] [amount] [odds_type] [ticket]",
		Short: "Place bet",
		Long:  "Place bet uuid, amount, odds type and ticket required.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			uid := args[0]
			argAmount := args[1]
			argOddsType := args[2]
			argTicket := args[3]

			oddsType, err := cast.ToInt32E(argOddsType)
			if err != nil {
				return types.ErrInvalidOddsType
			}

			argAmountCosmosInt, ok := sdk.NewIntFromString(argAmount)
			if !ok {
				return types.ErrInvalidAmount
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBet(
				clientCtx.GetFromAddress().String(),
				types.PlaceBetFields{
					UID:      uid,
					Amount:   argAmountCosmosInt,
					OddsType: types.OddsType(oddsType),
					Ticket:   argTicket,
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
