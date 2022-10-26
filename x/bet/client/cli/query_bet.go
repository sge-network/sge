package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cobra"
)

// CmdListBet implements a command to return all bets
func CmdListBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-all-bet",
		Short: "list all bet",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryListBetAllRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListAllBet(context.Background(), params)
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

// CmdShowBet implements a command to return a specific bet
func CmdShowBet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bet [uid]",
		Short: "shows a bet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argUID := args[0]

			params := &types.QueryBetRequest{
				Uid: argUID,
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
