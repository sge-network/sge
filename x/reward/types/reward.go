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
	SubaccountAddr  string
	MainAccountAddr string
	RewardAmount    RewardAmount
}

// RewardFactoryKeepers holds the keeper objects usable by reward types methods.
type RewardFactoryKeepers struct {
	OVMKeeper
	BetKeeper
	SubaccountKeeper
	RewardKeeper
	AccountKeeper
}

func (keepers *RewardFactoryKeepers) getSubaccountAddr(ctx sdk.Context, creator, receiver string) (string, error) {
	var (
		subAccountAddressString string
		err                     error
	)
	subAccountAddress, found := keepers.SubaccountKeeper.GetSubaccountByOwner(ctx, sdk.MustAccAddressFromBech32(receiver))
	if !found {
		subAccountAddressString, err = keepers.SubaccountKeeper.CreateSubaccount(ctx, creator, receiver, []subaccounttypes.LockedBalance{})
		if err != nil {
			return "", sdkerrors.Wrapf(ErrSubaccountCreationFailed, "%s", receiver)
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

func newRewardAmount(mainAccountAmount, subaccountAmount sdkmath.Int,
	mainAccountPercentage, subaccountPercentage sdk.Dec,
	unlockPeriod uint64) RewardAmount {
	return RewardAmount{
		MainAccountAmount:     mainAccountAmount,
		SubaccountAmount:      subaccountAmount,
		MainAccountPercentage: mainAccountPercentage,
		SubaccountPercentage:  subaccountPercentage,
		UnlockPeriod:          unlockPeriod,
	}
}

func NewReward(
	uid, creator string, receiver Receiver,
	campaignUID string,
	sourceUID, meta string,
) Reward {
	return Reward{
		UID:          uid,
		Creator:      creator,
		Receiver:     receiver.MainAccountAddr,
		CampaignUID:  campaignUID,
		RewardAmount: &receiver.RewardAmount,
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
func NewReceiver(maAddr, saAddr string,
	maAmount, saAmount sdkmath.Int,
	maPercentage, saPercentage sdk.Dec,
	unlockPeriod uint64) Receiver {
	return Receiver{
		SubaccountAddr:  saAddr,
		MainAccountAddr: maAddr,
		RewardAmount: newRewardAmount(
			maAmount,
			saAmount,
			maPercentage,
			saPercentage,
			unlockPeriod,
		),
	}
}

func (ds Receiver) String() string {
	out, err := yaml.Marshal(ds)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (ds Receiver) TotalAmount() sdkmath.Int {
	return ds.RewardAmount.MainAccountAmount.Add(ds.RewardAmount.SubaccountAmount)
}
