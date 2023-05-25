package types

type msgRewardToStoreReward struct{}

// Convert converts msgAwardeeToStoreAwardee
func (c msgRewardToStoreReward) Convert(msg *MsgRewardUser) (RewardK, error) {
	reward := RewardK{}

	awardees, err := MsgAwardeesToStoreAwardees.Convert(msg.Reward.Awardees)
	if err != nil {
		return reward, err
	}

	reward.RewardType = RewardK_RewardTypeK(RewardK_RewardTypeK_value[msg.Reward.RewardType.String()])
	reward.IncentiveUID = msg.Reward.IncentiveId
	reward.Awardees = awardees
	reward.Meta = msg.Reward.Meta

	return reward, nil
}

var MsgRewardToStoreReward = msgRewardToStoreReward{}
