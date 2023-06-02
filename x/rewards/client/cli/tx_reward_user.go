package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sge-network/sge/x/rewards/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRewardUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-user",
		Short: "accounts amounts type meta incentiveId ticket",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg, err := types.NewMsgRewardUser(
				clientCtx.GetFromAddress().String(), args[0], args[1], args[2], args[3], args[4], args[5],
			)
			if err != nil {
				return err
			}
			fmt.Println(args[0], args[1], args[2], args[3])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			creator, _ := sdk.AccAddressFromBech32(msg.Creator)
			fmt.Println(creator)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
