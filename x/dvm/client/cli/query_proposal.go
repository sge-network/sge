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

// CmdPubkeysChangeProposal implements a command to return a specific pubkeys change proposal based on its id
func CmdPubkeysChangeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pubkeys-change-proposal [status] [id]",
		Short: "query public keys change proposal",
		Long:  "query public keys change proposal by id.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argStatus, err := cast.ToInt32E(args[0])
			if err != nil {
				return err
			}

			argID, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryPublicKeysChangeProposalRequest{
				Id:     argID,
				Status: types.ProposalStatus(argStatus),
			}

			res, err := queryClient.PublicKeysChangeProposal(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdPubkeysChangeProposalList returns a command object instance for querying the public keys
func CmdPubkeysChangeProposalList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pubkeys-change-proposals [status]",
		Short: "query public keys change proposal list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			argStatus, err := cast.ToInt32E(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPublicKeysChangeProposalsRequest{
				Pagination: pageReq,
				Status:     types.ProposalStatus(argStatus),
			}

			res, err := queryClient.PublicKeysChangeProposals(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
