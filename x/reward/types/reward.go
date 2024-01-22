package types

import (
	context "context"

	"github.com/mrz1836/go-sanitize"
	yaml "gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Receiver struct {
	SubAccountAddr    string
	SubAccountAmount  sdkmath.Int
	UnlockPeriod      uint64
	MainAccountAddr   string
	MainAccountAmount sdkmath.Int
}

// RewardFactoryKeepers holds the keeper objectes usable by reward types methods.
type RewardFactoryKeepers struct {
	OVMKeeper
	BetKeeper
	SubAccountKeeper
}

func NewReward(
	uid, creator, receiver string,
	campaignUID string,
	rewardAmount *RewardAmount,
	sourceUID, meta string,
) Reward {
	return Reward{
		UID:          uid,
		Creator:      creator,
		Receiver:     receiver,
		CampaignUID:  campaignUID,
		RewardAmount: rewardAmount,
		SourceUID:    sourceUID,
		Meta:         sanitize.XSS(meta),
	}
}

func NewRewardByCategory(uid, addr string, rewardCategory RewardCategory) RewardByCategory {
	return RewardByCategory{
		UID:            uid,
		RewardCategory: rewardCategory,
		Addr:           addr,
	}
}

func NewRewardByCampaign(uid, campaignUID string) RewardByCampaign {
	return RewardByCampaign{
		UID:         uid,
		CampaignUID: campaignUID,
	}
}

// IRewardFactory defines the methods that should be implemented for all reward types.
type IRewardFactory interface {
	ValidateCampaign(campaign Campaign) error
	Calculate(
		goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers, campaign Campaign, ticket, creator string,
	) (Receiver, RewardPayloadCommon, error)
}

// NewReceiver creates reveiver object.
func NewReceiver(saAddr, maAddr string, saAmount, maAmount sdkmath.Int, unlockPeriod uint64) Receiver {
	return Receiver{
		SubAccountAddr:    saAddr,
		SubAccountAmount:  saAmount,
		UnlockPeriod:      unlockPeriod,
		MainAccountAddr:   maAddr,
		MainAccountAmount: maAmount,
	}
}

func (ds Receiver) String() string {
	out, err := yaml.Marshal(ds)
	if err != nil {
		panic(err)
	}
	return string(out)
}
