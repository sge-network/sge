package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/sge-network/sge/x/house/types"
	"github.com/spf13/cobra"
)

func CmdWithdrawal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [sport_event_uid] [participant_number] [mode] [amount]",
		Args:  cobra.RangeArgs(3, 4),
		Short: "Withdraw tokens from a deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw coins of unused amount corresponding to a deposit.

				Example:
				$ %s tx house withdraw bc79a72c-ad7e-4cf5-91a2-98af2751e812 1 1 1000 --from mykey
				`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sportEventUID := args[0]

			particiapntNumber, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || particiapntNumber < 1 {
				return fmt.Errorf("particiapnt number argument provided must be a non-negative-integer: %v", err)
			}

			mode, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return fmt.Errorf("mode argument provided must be a non-negative-integer: %v", mode)
			}

			var argAmountCosmosInt sdk.Int
			if mode == int64(types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL) {
				if len(args) != 4 {
					return fmt.Errorf("amount is mandatory for partial mode")
				}

				var ok bool
				argAmountCosmosInt, ok = sdk.NewIntFromString(args[1])
				if !ok {
					return types.ErrInvalidAmount
				}
			}

			depAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgWithdraw(depAddr.String(), sportEventUID, argAmountCosmosInt, particiapntNumber, types.WithdrawalMode(mode))

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}