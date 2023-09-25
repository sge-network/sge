package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

func NewCampaign(
	creator, funderAddres, uID string,
	startTS, endTS uint64,
	rewardType RewardType,
	rewardDefs []Definition,
	pool Pool,
) Campaign {
	return Campaign{
		Creator:       creator,
		FunderAddress: funderAddres,
		UID:           uID,
		StartTS:       startTS,
		EndTS:         endTS,
		RewardType:    rewardType,
		RewardDefs:    rewardDefs,
		Pool:          pool,
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
func (c *Campaign) CheckPoolBalance(distributions []Distribution) error {
	totalAmount := sdkmath.ZeroInt()
	for _, d := range distributions {
		totalAmount.Add(d.Allocation.Amount)
	}
	if err := c.Pool.CheckBalance(totalAmount); err != nil {
		return err
	}
	return nil
}
