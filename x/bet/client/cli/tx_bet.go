package cli

import (
	"encoding/json"
	"strings"

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
		Use:   "place-bet [uid] [amount] [ticket]",
		Short: "Place a new bet",
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

			msg := types.NewMsgPlaceBet(
				clientCtx.GetFromAddress().String(),
				types.BetPlaceFields{
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

// CmdPlaceBetSlip implements a command to place and store multiple bets
func CmdPlaceBetSlip() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bet-slip [bets]",
		Short: "Place multiple bets (for each bet: uid, amount, ticket)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBets := args[0]
			bets := []*types.BetPlaceFields{}
			err = json.Unmarshal([]byte(argBets), &bets)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBetSlip(
				clientCtx.GetFromAddress().String(),
				bets,
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

// CmdSettleBetBulk implements a command to settle multiple bets
func CmdSettleBetBulk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settle-bet-bulk [bet-uids]",
		Short: "Settle multiple bets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBetUIDs := args[0]
			argBetUIDs = strings.ReplaceAll(argBetUIDs, " ", "")
			betUIDs := []string{}
			if argBetUIDs != "" {
				betUIDs = strings.Split(argBetUIDs, ",")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSettleBetBulk(
				clientCtx.GetFromAddress().String(),
				betUIDs,
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
