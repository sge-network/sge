package types

import (
	"errors"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/legacy/bet/types"
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
		return sdkerrors.Wrap(sdkerrtypes.ErrInvalidRequest, err.Error())
	}

	switch payload.RewardAmountType {
	case RewardAmountType_REWARD_AMOUNT_TYPE_FIXED:
		if (!payload.RewardAmount.MainAccountPercentage.IsNil() &&
			payload.RewardAmount.MainAccountPercentage.GT(sdkmath.LegacyZeroDec())) ||
			(!payload.RewardAmount.SubaccountPercentage.IsNil() &&
				payload.RewardAmount.SubaccountPercentage.GT(sdkmath.LegacyZeroDec())) {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward percentage is not allowed for fixed amount campaign")
		}
		if (payload.RewardAmount.MainAccountAmount.IsNil() ||
			payload.RewardAmount.MainAccountAmount.LTE(sdkmath.ZeroInt())) &&
			(payload.RewardAmount.SubaccountAmount.IsNil() ||
				payload.RewardAmount.SubaccountAmount.LTE(sdkmath.ZeroInt())) {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward amount should be set for at least one of main account or subaccount")
		}
	case RewardAmountType_REWARD_AMOUNT_TYPE_PERCENTAGE:
		if (!payload.RewardAmount.MainAccountAmount.IsNil() &&
			payload.RewardAmount.MainAccountAmount.GT(sdkmath.ZeroInt())) ||
			(!payload.RewardAmount.SubaccountAmount.IsNil() &&
				payload.RewardAmount.SubaccountAmount.GT(sdkmath.ZeroInt())) {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward amount is not allowed for percentage campaign")
		}
		if (payload.RewardAmount.MainAccountPercentage.IsNil() ||
			payload.RewardAmount.MainAccountPercentage.LTE(sdkmath.LegacyZeroDec())) &&
			(payload.RewardAmount.SubaccountPercentage.IsNil() ||
				payload.RewardAmount.SubaccountPercentage.LTE(sdkmath.LegacyZeroDec())) {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward percentage should be set for at least one of main account or subaccount")
		}
	default:
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "unsupported reward amount type")
	}

	if ((!payload.RewardAmount.SubaccountAmount.IsNil() &&
		payload.RewardAmount.SubaccountAmount.GT(sdkmath.ZeroInt())) ||
		(!payload.RewardAmount.SubaccountPercentage.IsNil() &&
			payload.RewardAmount.SubaccountPercentage.GT(sdkmath.LegacyZeroDec()))) &&
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
func (payload *CreatePromoterPayload) Validate() error {
	if !utils.IsValidUID(payload.UID) {
		return types.ErrInvalidBetUID
	}
	if err := payload.Conf.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates promoter config set ticket payload.
func (payload *SetPromoterConfPayload) Validate() error {
	if err := payload.Conf.Validate(); err != nil {
		return err
	}

	return nil
}
