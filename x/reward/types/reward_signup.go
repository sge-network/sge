package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// SignUpReward is the type for signup rewards calculations
type SignUpReward struct{}

// NewSignUpReward create new object of signup reward calculator type.
func NewSignUpReward() SignUpReward { return SignUpReward{} }

// VaidateCampaign validates campaign definitions.
func (sur SignUpReward) ValidateCampaign(campaign Campaign) error {
	if campaign.RewardCategory != RewardCategory_REWARD_CATEGORY_SIGNUP {
		return sdkerrors.Wrapf(ErrWrongRewardCategory, "signup rewards can only have single definition")
	}
	if campaign.RewardAmount.SubaccountAmount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrWrongAmountForType, "signup rewards for subaccount should be positive")
	}
	if campaign.RewardAmountType != RewardAmountType_REWARD_AMOUNT_TYPE_FIXED {
		return sdkerrors.Wrapf(ErrWrongRewardAmountType, "reward amount type not supported for given reward type.")
	}

	return nil
}

// Calculate parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpReward) Calculate(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	campaign Campaign, ticket, creator string,
) (Receiver, RewardPayloadCommon, error) {
	var payload GrantSignupRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return Receiver{}, payload.Common, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	var (
		subAccountAddressString string
		err                     error
	)

	if err = payload.Common.Validate(); err != nil {
		return Receiver{}, payload.Common, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "%s", err)
	}

	addr, err := sdk.AccAddressFromBech32(payload.Common.Receiver)
	if err != nil {
		return Receiver{}, payload.Common, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	if keepers.SubAccountKeeper.IsSubAccount(ctx, addr) {
		return Receiver{}, payload.Common, ErrReceiverAddrCanNotBeSubAcc
	}

	subAccountAddress, found := keepers.SubAccountKeeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(payload.Common.Receiver))
	if !found {
		subAccountAddressString, err = keepers.SubAccountKeeper.CreateSubAccount(ctx, creator, payload.Common.Receiver, []subaccounttypes.LockedBalance{})
		if err != nil {
			return Receiver{}, payload.Common, sdkerrors.Wrapf(ErrSubAccountCreationFailed, "%s", payload.Common.Receiver)
		}
	} else {
		subAccountAddressString = subAccountAddress.String()
	}

	return NewReceiver(
		subAccountAddressString,
		payload.Common.Receiver,
		campaign.RewardAmount.SubaccountAmount,
		campaign.RewardAmount.MainAccountAmount,
		campaign.RewardAmount.UnlockPeriod,
	), payload.Common, nil
}
