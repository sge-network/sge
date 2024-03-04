package types

type RewardPoolFunder struct{}

func (RewardPoolFunder) GetModuleAcc() string {
	return rewardPool
}
