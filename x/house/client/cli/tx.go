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
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/house/types"
)

// NewTxCmd returns a root CLI command handler for all x/house transaction commands.
func NewTxCmd() *cobra.Command {
	houseTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "House transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	houseTxCmd.AddCommand(
		NewDepositCmd(),
		NewWithdrawalCmd(),
	)

	return houseTxCmd
}

func NewDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [sport_event_uid] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens to be part of a house",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit coins to be part of a house corresponding to a sport event.

				Example:
				$ %s tx house deposit bc79a72c-ad7e-4cf5-91a2-98af2751e812 1000usge --from mykey
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

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			depAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgDeposit(depAddr, sportEventUID, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewWithdrawalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [sport_event_uid] [participant_number] [mode] [amount]",
		Args:  cobra.RangeArgs(3, 4),
		Short: "Withdraw tokens from a deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw coins of unused amount corresponding to a deposit.

				Example:
				$ %s tx house withdraw bc79a72c-ad7e-4cf5-91a2-98af2751e812 1 1 1000usge --from mykey
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

			var amount sdk.Coin
			if mode == int64(types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL) {
				if len(args) != 4 {
					return fmt.Errorf("amount is mandatory for partial mode")
				}
				amount, err = sdk.ParseCoinNormalized(args[3])
				if err != nil {
					return err
				}
			}

			depAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgWithdraw(depAddr, sportEventUID, amount, particiapntNumber, types.WithdrawalMode(mode))

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
