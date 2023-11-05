package types

import (
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

	if len(payload.RewardDefs) == 0 {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "at least one reward definition should be defined")
	}

	if payload.TotalFunds.IsNil() || !payload.TotalFunds.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "pool amount should be positive")
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
