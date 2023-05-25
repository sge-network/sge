package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	yaml "gopkg.in/yaml.v2"
)

func NewRewardK(ctx sdk.Context, msg *MsgRewardUser) (RewardK, error) {
	storeRewards, err := MsgRewardToStoreReward.Convert(msg)
	if err != nil {
		return RewardK{}, sdkerrors.Wrap(err, "Unable to convert MsgRewardUser to store object")
	}
	return storeRewards, nil
}

// String returns a human-readable string representation of a Deposit.
func (d *RewardK) String() string {
	out, err := yaml.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(out)
}
