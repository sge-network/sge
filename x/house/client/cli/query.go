package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/house/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	houseQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the house module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	houseQueryCmd.AddCommand(
		GetCmdQueryDepsoits(),
		GetCmdQueryDepositorDepsoits(),
		GetCmdQueryDepositorWithdrawals(),
	)

	return houseQueryCmd
}

// GetCmdQueryDepsoits implements the query all deposits command.
func GetCmdQueryDepsoits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Short: "Query for all deposits",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all deposits on a network.

Example:
$ %s query house deposits
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			result, err := queryClient.Deposits(cmd.Context(), &types.QueryDepositsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "deposits")

	return cmd
}

// GetCmdQueryDepositorDepsoits implements the command to query all the deposits made from one depositor.
func GetCmdQueryDepositorDepsoits() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "depositor-deposits [depositor-addr]",
		Short: "Query all deposits made by one depositor",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query deposits for an individual depositor on all houses.

Example:
$ %s query house depositor-deposits %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
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

			params := &types.QueryDepositorDepositsRequest{
				DepositorAddress: depAddr.String(),
				Pagination:       pageReq,
			}

			res, err := queryClient.DepositorDeposits(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "deposits")

	return cmd
}

// GetCmdQueryDepositorWithdrawals implements the command to query all the withdrawals made from one depositor.
func GetCmdQueryDepositorWithdrawals() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "depositor-withdrawals [depositor-addr]",
		Short: "Query all withdrawals made by one depositor",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query withdrawals for an individual depositor on all houses.

Example:
$ %s query house depositor-withdrawals %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
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

			params := &types.QueryDepositorWithdrawalsRequest{
				DepositorAddress: depAddr.String(),
				Pagination:       pageReq,
			}

			res, err := queryClient.DepositorWithdrawals(cmd.Context(), params)
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
