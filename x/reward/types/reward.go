package types

import (
	context "context"

	"github.com/mrz1836/go-sanitize"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
	yaml "gopkg.in/yaml.v2"

	sdkerrors "cosmossdk.io/errors"
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

// RewardFactoryKeepers holds the keeper objects usable by reward types methods.
type RewardFactoryKeepers struct {
	OVMKeeper
	BetKeeper
	SubAccountKeeper
	RewardKeeper
	AccountKeeper
}

func (keepers *RewardFactoryKeepers) getSubAccAddr(ctx sdk.Context, creator, receiver string) (string, error) {
	var (
		subAccountAddressString string
		err                     error
	)
	subAccountAddress, found := keepers.SubAccountKeeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(receiver))
	if !found {
		subAccountAddressString, err = keepers.SubAccountKeeper.CreateSubAccount(ctx, creator, receiver, []subaccounttypes.LockedBalance{})
		if err != nil {
			return "", sdkerrors.Wrapf(ErrSubAccountCreationFailed, "%s", receiver)
		}
	} else {
		subAccountAddressString = subAccountAddress.String()
	}
	return subAccountAddressString, nil
}

// RewardFactoryData holds the data usable by reward types methods.
type RewardFactoryData struct {
	Receiver Receiver
	Common   RewardPayloadCommon
}

func NewRewardFactoryData(receiver Receiver, common RewardPayloadCommon) RewardFactoryData {
	return RewardFactoryData{
		Receiver: receiver,
		Common:   common,
	}
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
	) (RewardFactoryData, error)
}

// NewReceiver creates receiver object.
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
