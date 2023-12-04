package cli

import (
	"github.com/spf13/cobra"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/sge-network/sge/x/reward/types"
)

func CmdCreateCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-campaign [uid] [totalfunds] [ticket]",
		Short: "Create a new campaign",
		Long:  "Creating a new campaign with certain amount of funds and the ticket",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			argUID := args[0]

			argTotalFundsCosmosInt, ok := sdkmath.NewIntFromString(args[1])
			if !ok {
				return types.ErrInvalidFunds
			}

			// Get value arguments
			argTicket := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCampaign(
				clientCtx.GetFromAddress().String(),
				argUID,
				argTotalFundsCosmosInt,
				argTicket,
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

func CmdUpdateCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-campaign [uid] [topupfunds] [ticket]",
		Short: "Update a campaign",
		Long:  "Updating a new campaign with certain amount of funds and the ticket",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			argUID := args[0]

			argTopupFundsCosmosInt, ok := sdkmath.NewIntFromString(args[1])
			if !ok {
				return types.ErrInvalidFunds
			}

			// Get value arguments
			argTicket := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCampaign(
				clientCtx.GetFromAddress().String(),
				argUID,
				argTopupFundsCosmosInt,
				argTicket,
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

func CmdWithdrawFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-funds [uid] [ticket]",
		Short: "Withdraw funds from a campaign",
		Long:  "Withdrawal of the funds from a certain campaign",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			argUID := args[0]

			// Get value arguments
			argTicket := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFunds(
				clientCtx.GetFromAddress().String(),
				argUID,
				argTicket,
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
