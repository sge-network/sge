package types

import (
	context "context"

	yaml "gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Allocation struct {
	MainAcc Receiver
	SubAcc  Receiver
}

type Receiver struct {
	Addr     string
	Amount   sdkmath.Int
	UnlockTS uint64
}

// RewardFactoryKeepers holds the keeper objectes usable by reward types methods.
type RewardFactoryKeepers struct {
	OVMKeeper
	BetKeeper
	SubAccountKeeper
}

// IRewardFactory defines the methods that should be implemented for all of reward types.
type IRewardFactory interface {
	VaidateCampaign(campaign Campaign) error
	Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers, campaign Campaign, ticket string) (Allocation, error)
}

// NewReceiver creates reveiver object.
func NewReceiver(addr string, amount sdkmath.Int, unlockTS uint64) Receiver {
	return Receiver{
		Amount:   amount,
		Addr:     addr,
		UnlockTS: unlockTS,
	}
}

// // ValidateBasic validates the basic properties of a reward definition.
// // TODO: move logic to the new design
// func (d *Definition) ValidateBasic(blockTime uint64) error {
// 	if d.RecAccType != ReceiverAccType_RECEIVER_ACC_TYPE_SUB {
// 		if d.UnlockTS != 0 {
// 			return sdkerrors.Wrapf(ErrUnlockTSIsSubAccOnly, "%d", d.UnlockTS)
// 		}
// 	} else if d.UnlockTS <= blockTime {
// 		return sdkerrors.Wrapf(ErrUnlockTSDefBeforeBlockTime, "%d", d.UnlockTS)
// 	}
// 	return nil
// }

func (ds Allocation) String() string {
	out, err := yaml.Marshal(ds)
	if err != nil {
		panic(err)
	}
	return string(out)
}
