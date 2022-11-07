package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cobra"
)

// CmdPlaceBet implements a command to place and store a single bet
func CmdPlaceBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bet [uid] [amount] [odds_type] [ticket]",
		Short: "Place a new bet",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// Get value arguments
			uid := args[0]
			argAmount := args[1]
			argOddsType := args[2]
			argTicket := args[3]

			oddsType, err := strconv.ParseInt(argOddsType, 10, 32)
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

// CmdSettleBet implements a command to settle a bet
func CmdSettleBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settle-bet [bet-uid]",
		Short: "Settle a bet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBetUID := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSettleBet(
				clientCtx.GetFromAddress().String(),
				argBetUID,
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
