package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/house/types"
)

// GetCmdQueryWithdrawal implements the command to query the withdrawal of a
// certain market and participation with id.
func GetCmdQueryWithdrawal() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "withdrawal [depositor_address] [market_uid] [participation_index] [id]",
		Short: "Query withdrawal made by one account on certain participation and market with id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query withdrawal for an individual account on a certain house and participation by id.

Example:
$ %s query house withdrawal %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p 7d81a666-101c-11ee-be56-0242ac120002 1 1
`,
				version.AppName, bech32PrefixAccAddr,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			argDepositorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			argMarketUID := args[1]

			argParticipationIndex, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			argID, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			params := &types.QueryWithdrawalRequest{
				DepositorAddress:   argDepositorAddr.String(),
				MarketUid:          argMarketUID,
				ParticipationIndex: argParticipationIndex,
				Id:                 argID,
			}

			res, err := queryClient.Withdrawal(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "withdrawals")

	return cmd
}

// GetCmdQueryWithdrawalsByAccount implements the command to query all the withdrawals made from one account.
func GetCmdQueryWithdrawalsByAccount() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "withdrawals-by-account [account]",
		Short: "Query all withdrawals made by one account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query withdrawals for an individual accounts on all houses.

Example:
$ %s query house withdrawals-by-account %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.AppName, bech32PrefixAccAddr,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			depAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryWithdrawalsByAccountRequest{
				Address:    depAddr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.WithdrawalsByAccount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "withdrawals")

	return cmd
}
