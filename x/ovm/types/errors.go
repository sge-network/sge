package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/ovm module sentinel errors
var (
	ErrInvalidSignature       = sdkerrors.Register(ModuleName, 4101, "invalid signature")
	ErrInsufficientSignatures = sdkerrors.Register(ModuleName, 4102, "insufficient signatures for consensus")
	ErrInvalidTicketFormat    = sdkerrors.Register(ModuleName, 4103, "invalid ticket format")
	ErrExpirationRequired     = sdkerrors.Register(ModuleName, 4104, "ticket expiry must not be empty")
	ErrTicketExpired          = sdkerrors.Register(ModuleName, 4105, "ticket already expired")
	ErrMismatchKeyType        = sdkerrors.Register(ModuleName, 4106, "invalid public key")
	ErrShortKeyLength         = sdkerrors.Register(ModuleName, 4107, "insufficient key length")
	ErrNoPublicKeysFound      = sdkerrors.Register(ModuleName, 4108, "key not found")
	ErrKeyVaultNotFound       = sdkerrors.Register(ModuleName, 4109, "key vault not found")
)

const (
	// ErrTextInvalidRequest is the error message for an invalid request
	ErrTextInvalidRequest = "invalid signature"
)
