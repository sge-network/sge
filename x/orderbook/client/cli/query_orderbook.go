package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetCmdQueryOrderBooks implements the query all order books command.
func GetCmdQueryOrderBooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orderbooks",
		Short: "Query for all order books",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all order books on a network.

Example:
$ %s query orderbook orderbooks
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

			result, err := queryClient.OrderBooks(cmd.Context(), &types.QueryOrderBooksRequest{
				// Leaving status empty on purpose to query all orderbooks.
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "orderbooks")

	return cmd
}

// GetCmdQueryOrderBook implements the orderbook query command.
func GetCmdQueryOrderBook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orderbook [order-book-id]",
		Short: "Query a orderbook",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a order book.

Example:
$ %s query orderbook orderbook %s
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000",
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookUID := args[0]

			params := &types.QueryOrderBookRequest{OrderBookUid: orderBookUID}
			res, err := queryClient.OrderBook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.OrderBook)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
