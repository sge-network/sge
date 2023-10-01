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

	return gs.Params.Validate()
}
