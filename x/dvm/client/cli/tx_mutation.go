package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cobra"
)

// CmdMutation is the command object for mutation of public keys
func CmdMutation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mutation [txs]",
		Short: "updates list of public keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTxs := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMutation(
				clientCtx.GetFromAddress().String(),
				argTxs,
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
