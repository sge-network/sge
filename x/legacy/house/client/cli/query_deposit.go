package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/legacy/house/types"
)

// GetCmdQueryDeposits implements the query deposits command.
func GetCmdQueryDeposits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Short: "Query for all deposits in the house",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all deposits on a network.

Example:
$ %s query house deposits
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

// GetCmdQueryDepositsByAccount implements the command to query all the deposits made from one address.
func GetCmdQueryDepositsByAccount() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "deposits-by-account [account]",
		Short: "Query all deposits made by one account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query deposits for an individual account on all houses.

Example:
$ %s query house deposits-by-account %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
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

			params := &types.QueryDepositsByAccountRequest{
				Address:    depAddr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.DepositsByAccount(cmd.Context(), params)
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
