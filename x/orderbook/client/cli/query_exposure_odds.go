package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cobra"
)

// GetCmdQueryOrderBookExposure implements the book-exposure query command.
func GetCmdQueryOrderBookExposure() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-exposure [order-book-id] [odds-uid]",
		Short: "Query a book exposure",
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query details about a book exposure.

Example:
$ %s query orderbook book-exposure %s %s
`,
				version.AppName,
				"5531c60f-2025-48ce-ae79-1dc110f16000",
				"9991c60f-2025-48ce-ae79-1dc110f16990",
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookUID := args[0]
			oddsUID := args[1]

			params := &types.QueryOrderBookExposureRequest{OrderBookUid: orderBookUID, OddsUid: oddsUID}
			res, err := queryClient.OrderBookExposure(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.OrderBookExposure)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryOrderBookExposures implements the command to query all the exposures to a specific orderbook.
func GetCmdQueryOrderBookExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-exposures [order-book-id]",
		Short: "Query all book exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query exposures on an individual order book.

Example:
$ %s query orderbook book-exposures %s
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryOrderBookExposuresRequest{
				OrderBookUid: orderBookUID,
				Pagination:   pageReq,
			}

			res, err := queryClient.OrderBookExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "book exposures")

	return cmd
}
