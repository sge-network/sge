package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Definitions []Definition

// RewardFactoryKeepers holds the keeper objectes usable by reward types methods.
type RewardFactoryKeepers struct {
	OVMKeeper
	BetKeeper
}

// IRewardFactory defines the methods that should be implemented for all of reward types.
type IRewardFactory interface {
	VaidateDefinitions(campaign Campaign) error
	CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
		definitions Definitions, ticket string) ([]Distribution, error)
}

// NewDistribution created new distribution object.
func NewDistribution(accAddr string, allocation Allocation) Distribution {
	return Distribution{
		AccAddr:    accAddr,
		Allocation: allocation,
	}
}

// NewAllocation created allocation object.
func NewAllocation(amount sdkmath.Int, receiverAccType ReceiverAccType, unlockTS uint64) Allocation {
	return Allocation{
		Amount:     amount,
		RecAccType: receiverAccType,
		UnlockTS:   unlockTS,
	}
}

// ValidateBasic validates the basic properties of a reward definition.
func (d *Definition) ValidateBasic(blockTime uint64) error {
	if d.RecAccType != ReceiverAccType_RECEIVER_ACC_TYPE_SUB {
		if d.UnlockTS != 0 {
			return sdkerrors.Wrapf(ErrUnlockTSIsSubAccOnly, "%d", d.UnlockTS)
		}
	} else if d.UnlockTS <= blockTime {
		return sdkerrors.Wrapf(ErrUnlockTSDefBeforeBlockTime, "%d", d.UnlockTS)
	}
	return nil
}
