package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/dvm module sentinel errors
var (
	// ErrInvalidSignature in the case of signature do not match with any of the trusted public keys.
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 4101, "invalid signature")

	// ErrInsufficientSignatures in the case of not enough signature. 2/3 of trusted public keys must sign the ticket.
	ErrInsufficientSignatures = sdkerrors.Register(ModuleName, 4102, "insufficient signatures for consensus")

	// ErrInvalidTicketFormat if the format of ticket is not correct. like XXX.YYY.ZZZ
	ErrInvalidTicketFormat = sdkerrors.Register(ModuleName, 4103, "invalid ticket format")

	// ErrExpirationRequired if the expiration field is not included in the signed ticket.
	ErrExpirationRequired = sdkerrors.Register(ModuleName, 4104, "ticket expiry must not be empty")

	// ErrTicketExpired if the ticket is expired.
	ErrTicketExpired = sdkerrors.Register(ModuleName, 4105, "ticket already expired")

	// ErrMismatchKeyType if the type of key which is going to be added is not a valid public key.
	ErrMismatchKeyType = sdkerrors.Register(ModuleName, 4106, "invalid public key")

	// ErrShortKeyLength if the length of public key is not long enough.
	ErrShortKeyLength = sdkerrors.Register(ModuleName, 4107, "insufficient key length")

	// ErrNoPublicKeysFound if public key is not found
	ErrNoPublicKeysFound = sdkerrors.Register(ModuleName, 4108, "key not found")
)

const (
	// ErrTextInvalidRequest is the error message for an invalid request
	ErrTextInvalidRequest = "invalid signature"
)
