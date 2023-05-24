package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RewardUser(ctx sdk.Context, creator string, rewardType string, amount uint64, awardee string) error {
	awardeeAddress, err := sdk.AccAddressFromBech32(awardee)
	if err != nil {
		return err
	}

	fmt.Println("Sending amount: ", amount, "to: ", awardee)
	err = k.srKeeper.RewardUser(ctx, awardeeAddress, sdk.NewIntFromUint64(amount), rewardType)
	if err != nil {
		return err
	}
	return nil
}
