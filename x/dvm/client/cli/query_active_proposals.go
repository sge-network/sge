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

// CmdActivePubkeysChangeProposal implements a command to return a specific active pubkeys change proposal based on its id
func CmdActivePubkeysChangeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-pubkeys-change-proposal [id]",
		Short: "query active public keys change proposal",
		Long:  "query active public keys change proposal by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argID, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryActivePublicKeysChangeProposalRequest{
				Id: argID,
			}

			res, err := queryClient.ActivePublicKeysChangeProposal(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdActivePubkeysChangeProposalList returns a command object instance for querying the active public keys
func CmdActivePubkeysChangeProposalList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-pubkeys-change-proposals",
		Short: "query active public keys change proposal list",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryActivePublicKeysChangeProposalsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ActivePublicKeysChangeProposals(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
