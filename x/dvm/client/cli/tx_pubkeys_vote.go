package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// CmdChangePubkeysListVote is the command object for voting on change of public keys
func CmdChangePubkeysListVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-pubkeys-proposal [ticket] [pubkey]",
		Short: "creates a vote on the proposal to update list of public keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argTicket := args[0]
			pubkeyTicket := args[1]

			msg := types.NewMsgVotePubkeysChangeRequest(
				clientCtx.GetFromAddress().String(),
				argTicket,
				pubkeyTicket,
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
