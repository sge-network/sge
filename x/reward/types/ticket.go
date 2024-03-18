package types

import (
	"errors"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates campaign creation ticket payload.
func (payload *CreateCampaignPayload) Validate(blockTime uint64) error {
	_, err := sdk.AccAddressFromBech32(payload.Promoter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if payload.StartTs >= payload.EndTs {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "start timestamp can not be after end time")
	}

	if payload.EndTs <= blockTime {
		return sdkerrors.Wrapf(ErrExpiredCampaign, "%d", payload.EndTs)
	}

	if err := payload.validateRewardCategory(); err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, err.Error())
	}

	if (payload.RewardAmount.MainAccountAmount.IsNil() ||
		payload.RewardAmount.MainAccountAmount.Equal(sdkmath.ZeroInt())) &&
		(payload.RewardAmount.SubaccountAmount.IsNil() ||
			payload.RewardAmount.SubaccountAmount.Equal(sdkmath.ZeroInt())) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward amount should be set for atleast one of main account or subaccount")
	}

	if payload.RewardAmount.SubaccountAmount.GT(sdkmath.ZeroInt()) &&
		payload.RewardAmount.UnlockPeriod == 0 {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "sub account should have unlock period")
	}

	return nil
}

func (payload *CreateCampaignPayload) validateRewardCategory() error {
	err := errors.New("reward category is not compatible with reward type")
	switch payload.Category {
	case RewardCategory_REWARD_CATEGORY_SIGNUP:
		if payload.RewardType != RewardType_REWARD_TYPE_SIGNUP &&
			payload.RewardType != RewardType_REWARD_TYPE_REFERRAL_SIGNUP &&
			payload.RewardType != RewardType_REWARD_TYPE_AFFILIATE_SIGNUP {
			return err
		}
	case RewardCategory_REWARD_CATEGORY_BET_DISCOUNT:
		if payload.RewardType != RewardType_REWARD_TYPE_BET_DISCOUNT {
			return err
		}
	case RewardCategory_REWARD_CATEGORY_AFFILIATE:
		if payload.RewardType != RewardType_REWARD_TYPE_AFFILIATE {
			return err
		}
	case RewardCategory_REWARD_CATEGORY_MILESTONE:
		if payload.RewardType != RewardType_REWARD_TYPE_MILESTONE {
			return err
		}
	case RewardCategory_REWARD_CATEGORY_REFERRAL:
		if payload.RewardType != RewardType_REWARD_TYPE_REFERRAL {
			return err
		}
	default:
		return errors.New("unknown category reward")
	}
	return nil
}

// Validate validates campaign update ticket payload.
func (payload *UpdateCampaignPayload) Validate(blockTime uint64) error {
	if payload.EndTs < blockTime {
		return sdkerrors.Wrapf(ErrExpiredCampaign, "%d", payload.EndTs)
	}

	return nil
}

// Validate validates campaign withdraw funds ticket payload.
func (payload *WithdrawFundsPayload) Validate() error {
	_, err := sdk.AccAddressFromBech32(payload.Promoter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid promoter address (%s)", err)
	}
	return nil
}

// Validate validates common reward ticket payload.
func (payload *RewardPayloadCommon) Validate() error {
	if payload.KycData == nil {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", payload.Receiver)
	}
	if !payload.KycData.Validate(payload.Receiver) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", payload.Receiver)
	}
	return nil
}

// Validate validates promoter config set ticket payload.
func (payload *SetPromoterConfPayload) Validate() error {
	catMap := make(map[RewardCategory]struct{})
	for _, v := range payload.Conf.CategoryCap {
		_, ok := catMap[v.Category]
		if ok {
			return sdkerrors.Wrapf(ErrDuplicateCategoryInConf, "%s", v.Category)
		}
		if v.CapPerAcc <= 0 {
			return sdkerrors.Wrapf(ErrCategoryCapShouldBePos, "%s", v.Category)
		}

		catMap[v.Category] = struct{}{}
	}

	return nil
}
