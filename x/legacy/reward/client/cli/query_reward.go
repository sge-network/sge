package cli

import (
	"context"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sge-network/sge/x/legacy/reward/types"
)

func CmdListReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards",
		Short: "list all rewards",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRewardsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Rewards(context.Background(), params)
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

func CmdGetReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward [uid]",
		Short: "shows a reward",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argUID := args[0]

			params := &types.QueryRewardRequest{
				Uid: argUID,
			}

			res, err := queryClient.Reward(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetRewardsByUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards-by-address [promoter_uid] [address]",
		Short: "shows a list of rewards by promoter and address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPromoterUID := args[0]
			argAddress := args[1]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryRewardsByAddressRequest{
				PromoterUid: argPromoterUID,
				Address:     argAddress,
				Pagination:  pageReq,
			}

			res, err := queryClient.RewardsByAddress(context.Background(), params)
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

func CmdGetRewardsByCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards-by-campaign [campaign]",
		Short: "shows a list of rewards by campaign id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argCampaignID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryRewardsByCampaignRequest{
				CampaignUid: argCampaignID,
				Pagination:  pageReq,
			}

			res, err := queryClient.RewardsByCampaign(context.Background(), params)
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

func CmdGetRewardByUserAndCategory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards-by-user-category [promoter_uid] [address] [category]",
		Short: "shows a list of rewards by promoter, user, and category",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPromoterUID := args[0]
			argUser := args[1]
			argRewCategoryInt32, err := cast.ToInt32E(args[2])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryRewardsByAddressAndCategoryRequest{
				PromoterUid: argPromoterUID,
				Address:     argUser,
				Category:    types.RewardCategory(argRewCategoryInt32),
				Pagination:  pageReq,
			}

			res, err := queryClient.RewardsByAddressAndCategory(context.Background(), params)
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
