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

// CmdListBetByCreator implements a command to return all bets of a certain creator address
func CmdListBetByCreator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bets-by-creator [creator]",
		Short: "get list of bets for a creator-address",
		Long:  "Get list of bets for a creator address in paginated response.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			argCreator := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBetsByCreatorRequest{
				Creator:    argCreator,
				Pagination: pageReq,
			}

			res, err := queryClient.BetsByCreator(context.Background(), params)
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

// CmdListActiveBets implements a command to return all active bets of a market
func CmdListActiveBets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-bets [market-uid]",
		Short: "get list of active bets of a market",
		Long:  "Get list of active bets of a market in paginated response.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			argMarketUID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryActiveBetsRequest{
				MarketUid:  argMarketUID,
				Pagination: pageReq,
			}

			res, err := queryClient.ActiveBets(context.Background(), params)
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
		Short: "Query bets list of bettor by creator:UID list",
		Long:  "Get list of bets creator:UID comma separated list ex: \"address1:uid1,address2:uid2\" .",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var items []string
			reqItems := strings.Split(args[0], listSeparator)
			for _, val := range reqItems {
				items = append(items, val)
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
		Short: "bet details by creator-address and uid",
		Long:  "Get bet details by bet-creator-address address and uid.",
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
