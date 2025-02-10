package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	housetypes "github.com/sge-network/sge/x/legacy/house/types"
	"github.com/sge-network/sge/x/legacy/subaccount/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		TxCreate(),
		TxTopup(),
		TxWager(),
		TxHouseDeposit(),
		TxHouseWithdraw(),
		TxWithdraw(),
	)

	return cmd
}

func TxCreate() *cobra.Command {
	const (
		flagFunds        = "funds"
		flagLockDuration = "lock-duration"
	)

	cmd := &cobra.Command{
		Use:     "create [subaccount-owner]",
		Short:   "Create a new subaccount",
		Long:    `Create a new subaccount.`,
		Example: fmt.Sprintf(`$ %s tx subaccount create sge123456 --funds 1000000000 --lock-duration 8760h --from subaccount-funder-key`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			subaccountOwner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			fundsStr, err := cmd.Flags().GetString(flagFunds)
			if err != nil {
				return err
			}
			funds, ok := sdkmath.NewIntFromString(fundsStr)
			if !ok {
				return fmt.Errorf("invalid funds amount: %s", fundsStr)
			}
			unlocksAfter, err := cmd.Flags().GetDuration(flagLockDuration)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &types.MsgCreate{
				Creator: clientCtx.GetFromAddress().String(),
				Owner:   subaccountOwner.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: cast.ToUint64(time.Now().Add(unlocksAfter).Unix()),
						Amount:   funds,
					},
				},
			})
		},
	}
	cmd.Flags().String(flagFunds, "", "Funds to lock in the subaccount")
	cmd.Flags().Duration(flagLockDuration, 12*30*24*time.Hour, "duration for which the funds will be locked")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func TxTopup() *cobra.Command {
	const (
		flagFunds        = "funds"
		flagLockDuration = "lock-duration"
	)
	cmd := &cobra.Command{
		Use:     "topup [subaccount-owner]",
		Short:   "Topup a subaccount",
		Long:    `Topup a subaccount.`,
		Example: fmt.Sprintf(`$ %s tx subaccount topup sge123456 --funds 1000000000 --lock-duration 8760h --from funder-address-key`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			subaccountAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			fundsStr, err := cmd.Flags().GetString(flagFunds)
			if err != nil {
				return err
			}
			funds, ok := sdkmath.NewIntFromString(fundsStr)
			if !ok {
				return fmt.Errorf("invalid funds amount: %s", fundsStr)
			}
			unlocksAfter, err := cmd.Flags().GetDuration(flagLockDuration)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &types.MsgTopUp{
				Creator: clientCtx.GetFromAddress().String(),
				Address: subaccountAddress.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: cast.ToUint64(time.Now().Add(unlocksAfter).Unix()),
						Amount:   funds,
					},
				},
			})
		},
	}

	cmd.Flags().String(flagFunds, "", "Funds to lock in the subaccount")
	cmd.Flags().Duration(flagLockDuration, 12*30*24*time.Hour, "duration for which the funds will be locked")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func TxWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw --from subaccount-owner-key",
		Short:   "Withdraw unlocked funds from a subaccount",
		Long:    `Withdraw unlocked funds from a subaccount.`,
		Example: fmt.Sprintf(`$ %s tx subaccount withdraw --from subaccount-owner-key`, version.AppName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &types.MsgWithdrawUnlockedBalances{
				Creator: clientCtx.GetFromAddress().String(),
			})
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// TxWager implements a command to place and store a single bet
func TxWager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wager [ticket] --from subaccount-owner-key",
		Short: "Wager on an odds",
		Long:  "Wager on an odds. the uuid, amount and ticket required.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			argTicket := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWager(
				clientCtx.GetFromAddress().String(),
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

func TxHouseDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "house-deposit [market_uid] [amount] [ticket] --from subaccount-owner-key",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens in a market order book to be the house",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit tokens in a market order book to be the house.

				Example:
				$ %s tx subaccount deposit bc79a72c-ad7e-4cf5-91a2-98af2751e812 1000usge {ticket string} --from mykey
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

			argAmountCosmosInt, ok := sdkmath.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount: %s", args[1])
			}

			argTicket := args[2]

			depAddr := clientCtx.GetFromAddress()

			msg := &types.MsgHouseDeposit{Msg: housetypes.NewMsgDeposit(depAddr.String(), argMarketUID, argAmountCosmosInt, argTicket)}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxHouseWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "house-withdraw [market_uid] [participation_index] [ticket] [mode] [amount]",
		Args:  cobra.RangeArgs(4, 5),
		Short: "Withdraw tokens from a deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw coins of unused amount corresponding to a deposit.

				Example:
				$ %s tx subaccount withdraw bc79a72c-ad7e-4cf5-91a2-98af2751e812 1 {ticket string} 1 1000 --from mykey
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
			if mode == int32(housetypes.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL) {
				if len(args) != 5 {
					return fmt.Errorf("amount is mandatory for partial mode")
				}

				var ok bool
				argAmountCosmosInt, ok = sdkmath.NewIntFromString(args[4])
				if !ok {
					return fmt.Errorf("invalid amount: %s", args[4])
				}
			}

			depAddr := clientCtx.GetFromAddress()

			msg := housetypes.NewMsgWithdraw(depAddr.String(), argMarketUID, argAmountCosmosInt,
				particiapntIndex, housetypes.WithdrawalMode(mode), argTicket)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &types.MsgHouseWithdraw{Msg: msg})
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
