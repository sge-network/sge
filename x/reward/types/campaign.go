package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

func NewCampaign(
	creator, promoter, uID string,
	startTS, endTS, claimsPerCategory uint64,
	rewardType RewardType,
	rewardCategory RewardCategory,
	rewardAmountType RewardAmountType,
	rewardAmount *RewardAmount,
	isActive bool,
	meta string,
	pool Pool,
) Campaign {
	return Campaign{
		Creator:           creator,
		Promoter:          promoter,
		UID:               uID,
		StartTS:           startTS,
		EndTS:             endTS,
		RewardCategory:    rewardCategory,
		RewardType:        rewardType,
		RewardAmountType:  rewardAmountType,
		RewardAmount:      rewardAmount,
		IsActive:          isActive,
		Meta:              meta,
		Pool:              pool,
		ClaimsPerCategory: claimsPerCategory,
	}
}

// GetRewardsFactory returns reward factory according to the campaign.
func (c *Campaign) GetRewardsFactory() (IRewardFactory, error) {
	switch c.RewardType {
	case RewardType_REWARD_TYPE_SIGNUP:
		return NewSignUpReward(), nil
	default:
		return nil, sdkerrors.Wrapf(ErrUnknownRewardType, "%d", c.RewardType)
	}
}

func (c *Campaign) CheckExpiration(blockTime uint64) error {
	if blockTime > c.EndTS {
		return sdkerrors.Wrapf(ErrCampaignEnded, "end timestamp %d, block time %d", c.EndTS, blockTime)
	}
	return nil
}

// CheckPoolBalance checks if a pool balance of a capaign has enough fund to pay the rewards.
func (c *Campaign) CheckPoolBalance(rewardAmount sdkmath.Int) error {
	if err := c.Pool.CheckBalance(rewardAmount); err != nil {
		return err
	}
	return nil
}
