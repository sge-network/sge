package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// GetCmdQueryOrderBookParticipationExposures implements the command to query all the participation exposures to a specific orderbook.
func GetCmdQueryOrderBookParticipationExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orderbook-participation-exposures [order-book-id]",
		Short: "Query all participation exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participation exposures on an individual order book.

Example:
$ %s query orderbook participation-exposures %s
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

			params := &types.QueryOrderBookParticipationExposuresRequest{
				OrderBookUid: orderBookUID,
				Pagination:   pageReq,
			}

			res, err := queryClient.OrderBookParticipationExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participation exposures")

	return cmd
}

// GetCmdQueryParticipationExposures implements the participation-exposure query command.
func GetCmdQueryParticipationExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participation-exposures [order-book-id] [participation-index]",
		Short: "Query all participation exposures for a specific order book and participation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a participation exposure.

Example:
$ %s query orderbook participation-exposure %s %d
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
			participationIndex, err := cast.ToUint64E(args[1])
			if err != nil || participationIndex < 1 {
				return fmt.Errorf(
					"participation index argument provided must be a non-negative-integer: %v",
					err,
				)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryParticipationExposuresRequest{
				OrderBookUid:       orderBookUID,
				ParticipationIndex: participationIndex,
				Pagination:         pageReq,
			}
			res, err := queryClient.ParticipationExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participation exposures")

	return cmd
}

// GetCmdQueryHistoricalParticipationExposures implements the command to query all the historical participation exposures to a specific orderbook.
func GetCmdQueryHistoricalParticipationExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-participation-exposures [order-book-id]",
		Short: "Query all historical participation exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query historical participation exposures on an individual order book.

Example:
$ %s query orderbook historical-participation-exposures %s
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

			params := &types.QueryHistoricalParticipationExposuresRequest{
				OrderBookUid: orderBookUID,
				Pagination:   pageReq,
			}

			res, err := queryClient.HistoricalParticipationExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "historical participation exposures")

	return cmd
}
