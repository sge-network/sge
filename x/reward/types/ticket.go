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

	if payload.TotalFunds.IsNil() || !payload.TotalFunds.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "pool amount should be positive")
	}

	if payload.RewardAmount.MainAccountAmount.Equal(sdkmath.ZeroInt()) &&
		payload.RewardAmount.SubaccountAmount.Equal(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "reward amount should be set for atleast one of main account or subaccount")
	}

	if payload.RewardAmount.SubaccountAmount.GT(sdkmath.ZeroInt()) &&
		payload.RewardAmount.UnlockPeriod == 0 {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "sub account should have unlock timestamp")
	}

	return nil
}

func (payload *CreateCampaignPayload) validateRewardCategory() error {
	var err = errors.New("reward category is not compatible with reward type")
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
