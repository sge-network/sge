package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/spf13/cobra"
)

// GetCmdQueryReserver initiates the query for querying the reserver
func GetCmdQueryReserver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reserver",
		Short: "Query the current reserver object",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryReserverRequest{}
			res, err := queryClient.Reserver(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Reserver)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
