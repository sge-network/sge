package cli

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/rewards/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRewardByIncentiveId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-by-incentive-id incentiveID",
		Short: "Query RewardByIncentiveId incentiveID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			incentiveID := args[0]

			if incentiveID == "" {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid incentiveId (%s)", err)
			}

			params := &types.QueryRewardByIncentiveIdRequest{
				IncentiveId: incentiveID,
			}
			res, err := queryClient.RewardByIncentiveId(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
