package keeper

import (
	"errors"

	"github.com/sge-network/sge/x/reward/types"
)

// GetRewardsFactory returns reward factory according to the campaign.
func GetRewardsFactory(rewardType types.RewardType) (types.IRewardFactory, error) {
	switch rewardType {
	case types.RewardType_REWARD_TYPE_SIGNUP:
		return types.NewSignUpReward(), nil
	case types.RewardType_REWARD_TYPE_REFERRAL:
		return types.NewReferralReward(), nil
	case types.RewardType_REWARD_TYPE_AFFILIATION:
		return types.NewAffiliationReward(), nil
	case types.RewardType_REWARD_TYPE_NOLOSS_BETS:
		return types.NewNoLossBetsReward(), nil
	default:
		return nil, errors.New("unknown reward")
	}

}
