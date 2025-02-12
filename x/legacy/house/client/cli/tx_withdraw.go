package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/legacy/house/types"
)

func CmdWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [market_uid] [participation_index] [ticket] [mode] [amount]",
		Args:  cobra.RangeArgs(4, 5),
		Short: "Withdraw tokens from a deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw coins of unused amount corresponding to a deposit.

				Example:
				$ %s tx house withdraw bc79a72c-ad7e-4cf5-91a2-98af2751e812 1 {ticket string} 1 1000 --from mykey
				`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argMarketUID := args[0]

			particiapntIndex, err := cast.ToUint64E(args[1])
			if err != nil || particiapntIndex < 1 {
				return fmt.Errorf("participant number should be a positive number")
			}

			argTicket := args[2]

			mode, err := cast.ToInt32E(args[3])
			if err != nil {
				return fmt.Errorf("mode provided must be a non-negative-integer: %v", mode)
			}

			var argAmountCosmosInt sdkmath.Int
			if mode == int32(types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL) {
				if len(args) != 5 {
					return fmt.Errorf("amount is mandatory for partial mode")
				}

				var ok bool
				argAmountCosmosInt, ok = sdkmath.NewIntFromString(args[4])
				if !ok {
					return types.ErrInvalidAmount
				}
			}

			depAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgWithdraw(depAddr.String(), argMarketUID, argAmountCosmosInt,
				particiapntIndex, types.WithdrawalMode(mode), argTicket)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
