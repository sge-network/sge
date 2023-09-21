package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetCmdQueryOrderBookParticipations implements the command to query all the participations to a specific orderbook.
func GetCmdQueryOrderBookParticipations() *cobra.Command {
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

			params := &types.QueryOrderBookParticipationsRequest{
				OrderBookUid: orderBookUID,
				Pagination:   pageReq,
			}

			res, err := queryClient.OrderBookParticipations(cmd.Context(), params)
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

// GetCmdQueryOrderBookParticipation implements the book-participation query command.
func GetCmdQueryOrderBookParticipation() *cobra.Command {
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

			participationIndex, err := cast.ToUint64E(args[1])
			if err != nil || participationIndex < 1 {
				return fmt.Errorf(
					"participant index argument provided must be a non-negative-integer: %v",
					err,
				)
			}

			params := &types.QueryOrderBookParticipationRequest{
				OrderBookUid:       orderBookUID,
				ParticipationIndex: participationIndex,
			}
			res, err := queryClient.OrderBookParticipation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.OrderBookParticipation)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
