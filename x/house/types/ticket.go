package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates deposit ticket payload.
func (payload *DepositTicketPayload) Validate(depositor string) error {
	if !payload.KycData.Validate(depositor) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", depositor)
	}

	return nil
}

// Validate validates withdrawal payload.
func (payload *WithdrawTicketPayload) Validate(depositor string) error {
	if !payload.KycData.Validate(depositor) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", depositor)
	}

	return nil
}
