package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
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

			params := &types.QueryOrderBookRequest{BookUid: orderBookUID}
			res, err := queryClient.OrderBook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Orderbook)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryBookParticipations implements the command to query all the participations to a specific orderbook.
func GetCmdQueryBookParticipations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-participations [order-book-id]",
		Short: "Query all book participations for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participations on an individual order book.

Example:
$ %s query orderbook book-participations %s
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

			params := &types.QueryBookParticipationsRequest{
				BookUid:    orderBookUID,
				Pagination: pageReq,
			}

			res, err := queryClient.BookParticipations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "book participations")

	return cmd
}

// GetCmdQueryBookParticipation implements the bookparticipation query command.
func GetCmdQueryBookParticipation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-participation [order-book-id] [participation-index]",
		Short: "Query a book participation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a book participation.

Example:
$ %s query orderbook book-participation %s %d
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000", 1,
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

			particiapationIndex, err := cast.ToUint64E(args[1])
			if err != nil || particiapationIndex < 1 {
				return fmt.Errorf("particiapnt index argument provided must be a non-negative-integer: %v", err)
			}

			params := &types.QueryBookParticipationRequest{BookUid: orderBookUID, ParticipationIndex: particiapationIndex}
			res, err := queryClient.BookParticipation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.BookParticipation)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryBookExposures implements the command to query all the exposures to a specific orderbook.
func GetCmdQueryBookExposures() *cobra.Command {
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

			params := &types.QueryBookExposuresRequest{
				BookUid:    orderBookUID,
				Pagination: pageReq,
			}

			res, err := queryClient.BookExposures(cmd.Context(), params)
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

// GetCmdQueryBookExposure implements the bookexposure query command.
func GetCmdQueryBookExposure() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-exposure [order-book-id] [odd-uid]",
		Short: "Query a book exposure",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a book exposure.

Example:
$ %s query orderbook book-exposure %s %s
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000", "9991c60f-2025-48ce-ae79-1dc110f16990",
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

			params := &types.QueryBookExposureRequest{BookUid: orderBookUID, OddsUid: oddsUID}
			res, err := queryClient.BookExposure(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.BookExposure)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
