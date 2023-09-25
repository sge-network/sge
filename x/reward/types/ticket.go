package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates campaign creation ticket payload.
func (payload *CreateCampaignPayload) Validate(ctx sdk.Context) error {
	_, err := sdk.AccAddressFromBech32(payload.FunderAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if payload.StartTs >= payload.EndTs {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "start timestamp can not be after end time")
	}

	if payload.EndTs > uint64(ctx.BlockTime().Unix()) {
		return sdkerrors.Wrapf(ErrExpiredCampaign, "%d", payload.EndTs)
	}

	if payload.StartTs >= payload.EndTs {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "start timestamp can not be after end time")
	}

	if len(payload.RewardDefs) == 0 {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "at least one reward definition should be defined")
	}

	if !payload.PoolAmount.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "pool amount should be positive")
	}

	return nil
}

// Validate validates campaign update ticket payload.
func (payload *UpdateCampaignPayload) Validate(ctx sdk.Context) error {
	if payload.EndTs > uint64(ctx.BlockTime().Unix()) {
		return sdkerrors.Wrapf(ErrExpiredCampaign, "%d", payload.EndTs)
	}

	return nil
}
