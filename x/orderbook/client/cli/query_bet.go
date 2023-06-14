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

// GetCmdQueryParticipationBets implements the command to query all the participation fulfilled bets to a specific orderbook.
func GetCmdQueryParticipationBets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participation-bets [order-book-id] [participation-index]",
		Short: "Query all participation fulfilled bets for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participation fulfilled bets on an individual order book.

Example:
$ %s query orderbook participation-bets %s %d
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

			orderBookUID := args[0]
			participationIndex, err := cast.ToUint64E(args[1])
			if err != nil || participationIndex < 1 {
				return fmt.Errorf(
					"particiapnt index argument provided must be a non-negative-integer: %v",
					err,
				)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryParticipationFulfilledBetsRequest{
				OrderBookUid:       orderBookUID,
				ParticipationIndex: participationIndex,
				Pagination:         pageReq,
			}

			res, err := queryClient.ParticipationFulfilledBets(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participation fulfilled bets")

	return cmd
}
