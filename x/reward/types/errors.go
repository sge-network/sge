package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/reward module sentinel errors
var (
	ErrInTicketVerification     = sdkerrors.Register(ModuleName, 7100, "ticket verification failed")
	ErrInTicketValidation       = sdkerrors.Register(ModuleName, 7101, "ticket validation failed")
	ErrAuthorizationNotFound    = sdkerrors.Register(ModuleName, 7102, "no authorization found")
	ErrAuthorizationNotAccepted = sdkerrors.Register(ModuleName, 7103, "authorization not accepted")
	ErrExpiredCampaign          = sdkerrors.Register(ModuleName, 7104, "campaign is expired")
	ErrCampaignPoolBalance      = sdkerrors.Register(ModuleName, 7105, "not enough campaign pool balance")
	ErrUnknownRewardType        = sdkerrors.Register(ModuleName, 7106, "unknown reward type")
)
