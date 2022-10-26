package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// CmdPubKeysList returns a command object instance for querying the public keys
func CmdPubKeysList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-pubkeys",
		Short: "Query public keys list",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryListPubKeyAllRequest{}

			res, err := queryClient.ListPubKeys(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
