package cli

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cobra"
)

// CmdListBet implements a command to return all bets
func CmdListBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bets",
		Short: "get list of bets",
		Long:  "Get list of bets in paginated response.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBetsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Bets(context.Background(), params)
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

// CmdListBetByUIDs returns command object for querying bets by uid list
func CmdListBetByUIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bets-by-uids [creator:uid]",
		Short: "Query bets list of bettor by UIDs",
		Long:  "Get list of bets by list of the bettor of uids.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqItems := strings.Split(args[0], listSeparator)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var items []*types.QueryBetRequest
			for _, val := range reqItems {
				pair := strings.Split(val, mapSeparator)
				items = append(items, &types.QueryBetRequest{
					Creator: pair[0],
					Uid:     pair[1],
				})
			}

			params := &types.QueryBetsByUIDsRequest{
				Items: items,
			}

			res, err := queryClient.BetsByUIDs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdShowBet implements a command to return a specific bet
func CmdShowBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bet [creator] [uid]",
		Short: "bet details of bettor by uid",
		Long:  "Get bet details by bettor address adn uid.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argCreator := args[0]
			argUID := args[1]

			params := &types.QueryBetRequest{
				Creator: argCreator,
				Uid:     argUID,
			}

			res, err := queryClient.Bet(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
