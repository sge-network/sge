package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/sge-network/sge/x/legacy/reward/types"
)

func CmdGrantReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [uid] [campaign uid] [ticket]",
		Short: "Grant a new reward for the campaign",
		Long:  "Grant a new reward for the campaign with the provided uid",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			argUID := args[0]
			argCampaignUID := args[1]

			// Get value arguments
			argTicket := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgGrantReward(
				clientCtx.GetFromAddress().String(),
				argUID,
				argCampaignUID,
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
