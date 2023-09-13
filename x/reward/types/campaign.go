package types

import (
	cosmerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

func NewCampaign(
	creator, uID string,
	startTS, endTS uint64,
	rewardType RewardType,
	rewardDefs []Definition,
	pool Pool,
) Campaign {
	return Campaign{
		Creator:    creator,
		UID:        uID,
		StartTS:    startTS,
		EndTS:      endTS,
		RewardType: rewardType,
		RewardDefs: rewardDefs,
		Pool:       pool,
	}
}

// GetRewardsFactory returns reward factory according to the campaign.
func (c *Campaign) GetRewardsFactory() (IRewardFactory, error) {
	switch c.RewardType {
	case RewardType_REWARD_TYPE_SIGNUP:
		return NewSignUpReward(), nil
	case RewardType_REWARD_TYPE_REFERRAL:
		return NewReferralReward(), nil
	case RewardType_REWARD_TYPE_AFFILIATION:
		return NewAffiliationReward(), nil
	case RewardType_REWARD_TYPE_NOLOSS_BETS:
		return NewNoLossBetsReward(), nil
	default:
		return nil, cosmerrors.Wrapf(ErrUnknownRewardType, "%d", c.RewardType)
	}
}

// CheckPoolBalance checks if a pool balance of a capaign has enough fund to pay the reward.
func (c *Campaign) CheckPoolBalance() error {
	totalAmount := sdkmath.ZeroInt()
	for _, d := range c.RewardDefs {
		totalAmount.Add(d.Amount)
	}
	if err := c.Pool.CheckBalance(totalAmount); err != nil {
		return err
	}
	return nil
}
