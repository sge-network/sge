package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sge-network/sge/x/market/types"
)

// CmdListMarkets implements a command to return all markets
func CmdListMarkets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markets",
		Short: "list markets",
		Long:  "Get list of markets in paginated response.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryMarketsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Markets(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdGetMarket implements a command to return a specific market based on its UID
func CmdGetMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market [uid]",
		Short: "get market",
		Long:  "Get market meta by uid.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argUID := args[0]

			params := &types.QueryMarketRequest{
				Uid: argUID,
			}

			res, err := queryClient.Market(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdListMarketByUIDs returns command object for querying markets by uid list
func CmdListMarketByUIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markets-by-uids [uid]",
		Short: "Query markets list by UIDs",
		Long:  "Get list of markets by list of uids.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqUID := strings.Split(args[0], listSeparator)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryMarketsByUIDsRequest{
				Uids: reqUID,
			}

			res, err := queryClient.MarketsByUIDs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
