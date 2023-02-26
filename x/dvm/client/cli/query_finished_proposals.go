package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// CmdFinishedPubkeysChangeProposal implements a command to return a specific finished pubkeys change proposal based on its id
func CmdFinishedPubkeysChangeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finished-pubkeys-change-proposal [id]",
		Short: "query finished public keys change proposal",
		Long:  "query finished public keys change proposal by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argID, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryFinishedPublicKeysChangeProposalRequest{
				Id: argID,
			}

			res, err := queryClient.FinishedPublicKeysChangeProposal(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdFinishedPubkeysChangeProposalList returns a command object instance for querying the finished public keys
func CmdFinishedPubkeysChangeProposalList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finished-pubkeys-change-proposals",
		Short: "Query finished public keys change proposal list",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFinishedPublicKeysChangeProposalsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.FinishedPublicKeysChangeProposals(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
