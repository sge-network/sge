package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates market add ticket payload.
func (payload *DepositTicketPayload) Validate(creator string) error {
	if !payload.KycData.validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	return nil
}

// Validate validates market add ticket payload.
func (payload *WithdrawTicketPayload) Validate(creator string) error {
	if !payload.KycData.validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	return nil
}

// validate checks whether the kyc data is valid for a particular depositor
// if the kyc is required
func (payload KycDataPayload) validate(address string) bool {
	// ignore is true means that kyc check should be ignored
	if payload.Ignore {
		return true
	}

	if payload.Approved && payload.ID == address {
		return true
	}

	return false
}
