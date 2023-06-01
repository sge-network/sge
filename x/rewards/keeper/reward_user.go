package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/rewards/types"
)

func (k Keeper) SetReward(ctx sdk.Context, reward types.RewardK) {
	rewardKey := types.GetRewardKey(reward.IncentiveUID)
	store := k.getRewardStore(ctx)
	b := k.cdc.MustMarshal(&reward)
	store.Set(rewardKey, b)
}

func (k Keeper) RewardUsers(ctx sdk.Context, msg *types.MsgRewardUser) error {
	if k.IsIncentiveIdPresent(ctx, msg.Reward.IncentiveId) {
		return sdkerrors.Wrapf(sdkerrors.ErrConflict, "IncentiveId already present")
	}
	storeRewards, err := types.NewRewardK(ctx, msg)
	if err != nil {
		return err
	}
	for _, awardee := range msg.Reward.Awardees {
		err := k.RewardUser(ctx, msg.Creator, msg.Reward.RewardType.String(), awardee.Amount, awardee.Address)
		if err != nil {
			return sdkerrors.Wrap(err, "Something failed")
		}
	}
	k.SetReward(ctx, storeRewards)
	return nil
}

func (k Keeper) RewardUser(ctx sdk.Context, creator string, rewardType string, amount uint64, awardee string) error {
	awardeeAddress, err := sdk.AccAddressFromBech32(awardee)
	if err != nil {
		return err
	}

	fmt.Println("Sending amount: ", amount, "to: ", awardee)
	err = k.srKeeper.RewardUser(ctx, awardeeAddress, sdk.NewIntFromUint64(amount), rewardType)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) IsIncentiveIdPresent(ctx sdk.Context, incentiveId string) bool {
	rewardStore := k.getRewardStore(ctx)
	rewardKey := types.GetRewardKey(incentiveId)
	if rewardStore.Get(rewardKey) == nil {
		return false
	}
	return true
}
