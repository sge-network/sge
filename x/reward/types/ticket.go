package types

import (
	cosmerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates campaign creation ticket payload.
func (payload *CreateCampaignPayload) Validate(ctx sdk.Context) error {
	_, err := sdk.AccAddressFromBech32(payload.FunderAddress)
	if err != nil {
		return cosmerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if payload.StartTs >= payload.EndTs {
		return cosmerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start timestamp can not be after end time")
	}

	if payload.EndTs > uint64(ctx.BlockTime().Unix()) {
		return cosmerrors.Wrapf(ErrExpiredCampaign, "%d", payload.EndTs)
	}

	if payload.StartTs >= payload.EndTs {
		return cosmerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start timestamp can not be after end time")
	}

	if len(payload.RewardDefs) == 0 {
		return cosmerrors.Wrapf(sdkerrors.ErrInvalidRequest, "at least one reward definition should be defined")
	}

	if !payload.PoolAmount.GT(sdkmath.ZeroInt()) {
		return cosmerrors.Wrapf(sdkerrors.ErrInvalidRequest, "pool amount should be positive")
	}

	return nil
}

// Validate validates campaign update ticket payload.
func (payload *UpdateCampaignPayload) Validate(ctx sdk.Context) error {
	if payload.EndTs > uint64(ctx.BlockTime().Unix()) {
		return cosmerrors.Wrapf(ErrExpiredCampaign, "%d", payload.EndTs)
	}

	return nil
}
