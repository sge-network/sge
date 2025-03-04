package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates deposit ticket payload.
func (payload *DepositTicketPayload) Validate(depositor string) error {
	_, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}

	if !payload.KycData.Validate(depositor) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", depositor)
	}

	return nil
}

// Validate validates withdrawal payload.
func (payload *WithdrawTicketPayload) Validate(depositor string) error {
	_, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}

	if !payload.KycData.Validate(depositor) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", depositor)
	}

	return nil
}
