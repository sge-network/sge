package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		CampaignList: []Campaign{},
		Params:       DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in campaign
	campaignIndexMap := make(map[string]struct{})

	for _, elem := range gs.CampaignList {
		index := string(GetCampaignKey(elem.UID))
		if _, ok := campaignIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for campaign")
		}
		campaignIndexMap[index] = struct{}{}
	}

	rewardIndexMap := make(map[string]struct{})
	for _, elem := range gs.RewardList {
		index := string(GetRewardKey(elem.UID))
		if _, ok := rewardIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for reward")
		}
		rewardIndexMap[index] = struct{}{}
	}

	rewardByRecCatIndexMap := make(map[string]struct{})
	for _, elem := range gs.RewardByCategoryList {
		index := string(GetRewardKey(elem.UID))
		if _, ok := rewardByRecCatIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for reward by receiver and category")
		}
		rewardByRecCatIndexMap[index] = struct{}{}
	}

	rewardByCampaignIndexMap := make(map[string]struct{})
	for _, elem := range gs.RewardByCampaignList {
		index := string(GetRewardKey(elem.UID))
		if _, ok := rewardByCampaignIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for reward by campaign")
		}
		rewardByCampaignIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
