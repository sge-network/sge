package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	orderBookQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the orderbook module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	orderBookQueryCmd.AddCommand(
		GetCmdQueryOrderBooks(),
		GetCmdQueryOrderBook(),
		GetCmdQueryBookParticipants(),
		GetCmdQueryBookParticipant(),
		GetCmdQueryBookExposures(),
		GetCmdQueryBookExposure(),
		GetCmdQueryParticipantExposures(),
		GetCmdQueryParticipantExposure(),
		GetCmdQueryHistoricalParticipantExposures(),
		GetCmdQueryParticipantBets(),
	)

	return orderBookQueryCmd
}

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

			orderBookID := args[0]

			params := &types.QueryOrderBookRequest{BookId: orderBookID}
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

// GetCmdQueryBookParticipants implements the command to query all the participants to a specific orderbook.
func GetCmdQueryBookParticipants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-participants [order-book-id]",
		Short: "Query all book participants for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participants on an individual order book.

Example:
$ %s query orderbook book-participants %s
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

			orderBookID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryBookParticipantsRequest{
				BookId:     orderBookID,
				Pagination: pageReq,
			}

			res, err := queryClient.BookParticipants(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "book participants")

	return cmd
}

// GetCmdQueryBookParticipant implements the bookparticipant query command.
func GetCmdQueryBookParticipant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-participant [order-book-id] [participant-number]",
		Short: "Query a book participant",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a book participant.

Example:
$ %s query orderbook book-participant %s %d
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

			orderBookID := args[0]

			particiapntNumber, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || particiapntNumber < 1 {
				return fmt.Errorf("particiapnt number argument provided must be a non-negative-integer: %v", err)
			}

			params := &types.QueryBookParticipantRequest{BookId: orderBookID, ParticipantNumber: particiapntNumber}
			res, err := queryClient.BookParticipant(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.BookParticipant)
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

			orderBookID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryBookExposuresRequest{
				BookId:     orderBookID,
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
		Use:   "book-exposure [order-book-id] [odd-id]",
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

			orderBookID := args[0]
			oddsID := args[1]

			params := &types.QueryBookExposureRequest{BookId: orderBookID, OddsId: oddsID}
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

// GetCmdQueryParticipantExposures implements the command to query all the participant exposures to a specific orderbook.
func GetCmdQueryParticipantExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participant-exposures [order-book-id]",
		Short: "Query all participant exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participant exposures on an individual order book.

Example:
$ %s query orderbook participant-exposures %s
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

			orderBookID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryParticipantExposuresRequest{
				BookId:     orderBookID,
				Pagination: pageReq,
			}

			res, err := queryClient.ParticipantExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participant exposures")

	return cmd
}

// GetCmdQueryParticipantExposure implements the participantexposure query command.
func GetCmdQueryParticipantExposure() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participant-exposure [order-book-id] [participant-number]",
		Short: "Query a participant exposure",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a participant exposure.

Example:
$ %s query orderbook participant-exposure %s %d
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

			orderBookID := args[0]
			particiapntNumber, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || particiapntNumber < 1 {
				return fmt.Errorf("participant number argument provided must be a non-negative-integer: %v", err)
			}

			params := &types.QueryParticipantExposureRequest{BookId: orderBookID, ParticipantNumber: particiapntNumber}
			res, err := queryClient.ParticipantExposure(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryHistoricalParticipantExposures implements the command to query all the historical participant exposures to a specific orderbook.
func GetCmdQueryHistoricalParticipantExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-participant-exposures [order-book-id]",
		Short: "Query all historical participant exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query historical participant exposures on an individual order book.

Example:
$ %s query orderbook historical-participant-exposures %s
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

			orderBookID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryHistoricalParticipantExposuresRequest{
				BookId:     orderBookID,
				Pagination: pageReq,
			}

			res, err := queryClient.HistoricalParticipantExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "historical participant exposures")

	return cmd
}

// GetCmdQueryParticipantBets implements the command to query all the participant fulfilled bets to a specific orderbook.
func GetCmdQueryParticipantBets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participant-bets [order-book-id] [participant-number]",
		Short: "Query all participant fulfilled bets for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participant fulfilled bets on an individual order book.

Example:
$ %s query orderbook participant-bets %s %d
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000", 2,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookID := args[0]
			particiapntNumber, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || particiapntNumber < 1 {
				return fmt.Errorf("particiapnt number argument provided must be a non-negative-integer: %v", err)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryParticipantFullfilledBetsRequest{
				BookId:            orderBookID,
				ParticipantNumber: particiapntNumber,
				Pagination:        pageReq,
			}

			res, err := queryClient.ParticipantFullfilledBets(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participant fulfilled bets")

	return cmd
}
