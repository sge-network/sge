package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates deposit ticket payload.
func (payload *DepositTicketPayload) Validate(creator string) error {
	if !payload.KycData.Validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	return nil
}

// Validate validates withdrawal payload.
func (payload *WithdrawTicketPayload) Validate(creator string) error {
	if !payload.KycData.Validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	return nil
}
