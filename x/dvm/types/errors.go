package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/dvm module sentinel errors
var (
	// ErrInvalidSignature in the case of signature do not match with any of the trusted public keys.
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 1101, "invalid signature")

	// ErrInsufficientSignatures in the case of not enough signature. 2/3 of trusted public keys must sign the ticket.
	ErrInsufficientSignatures = sdkerrors.Register(ModuleName, 1102, "insufficient signatures, two-third of the list must sign this request")

	// ErrInvalidTicketFormat if the format of ticket is not correct. like XXX.YYY
	ErrInvalidTicketFormat = sdkerrors.Register(ModuleName, 1103, "invalid ticket format")

	// ErrExpirationRequired if the expiration field did'nt included in the signed ticket.
	ErrExpirationRequired = sdkerrors.Register(ModuleName, 1104, "expiration is required for tickets")

	// ErrTicketExpired if the ticket is expired.
	ErrTicketExpired = sdkerrors.Register(ModuleName, 1105, "ticket already expired")

	// ErrMismatchKeyType if the type of key which is going to be added is not a valid public key.
	ErrMismatchKeyType = sdkerrors.Register(ModuleName, 1106, "mismatch key type")

	// ErrShortKeyLength if the length of public key is not long enough.
	ErrShortKeyLength = sdkerrors.Register(ModuleName, 1107, "short key length")

	// ErrShortKeyLength if the length of public key is not long enough.
	ErrNoPublicKeysFound = sdkerrors.Register(ModuleName, 1108, "not public keys found")
)

const (
	// ErrTextInvalidRequest is the error message for an invalid request
	ErrTextInvalidRequest = "invalid signature"
)
