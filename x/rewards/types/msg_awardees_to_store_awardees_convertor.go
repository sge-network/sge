package types

import (
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
)

type msgAwardeesToStoreAwardees struct{}

// Convert converts msgAwardeesToStoreAwardees
func (c msgAwardeesToStoreAwardees) Convert(msgAwardees []*Awardee) ([]*AwardeeK, error) {
	awardees := []*AwardeeK{}
	for _, msgAwardee := range msgAwardees {
		awardees = append(awardees, &AwardeeK{
			Address: msgAwardee.Address,
			Amount:  github_com_cosmos_cosmos_sdk_types.NewInt(int64(msgAwardee.Amount)),
		})
	}
	return awardees, nil
}

var MsgAwardeesToStoreAwardees = msgAwardeesToStoreAwardees{}
