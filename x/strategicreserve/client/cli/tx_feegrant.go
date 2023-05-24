package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/spf13/cobra"
)

// CmdInvokeFeeGrant CLI registration for invoking a feegrant command
func CmdInvokeFeeGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoke-feegrant [ticket]",
		Short: "invoke a new feegrant from data fee collector module account",
		Long:  "invoke a new feegrant from data fee collector module account with ticket.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgInvokeFeeGrant(
				clientCtx.GetFromAddress().String(),
				args[0],
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

// CmdRevokeFeeGrant CLI registration for revoking a feegrant command
func CmdRevokeFeeGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-feegrant [ticket]",
		Short: "revoke an existing feegrant from data fee collector module account",
		Long:  fmt.Sprintf("revoke an existing feegrant from data fee collector with ticket. this takes %d minutes to be removed from feegrant module of the cosmos.", types.DefaultAllowanceExpiration),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevokeFeeGrant(
				clientCtx.GetFromAddress().String(),
				args[0],
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
