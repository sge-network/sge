package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/reward module sentinel errors
var (
	ErrInTicketVerification          = sdkerrors.Register(ModuleName, 7100, "ticket verification failed")
	ErrInTicketValidation            = sdkerrors.Register(ModuleName, 7101, "ticket validation failed")
	ErrAuthorizationNotFound         = sdkerrors.Register(ModuleName, 7102, "no authorization found")
	ErrAuthorizationNotAccepted      = sdkerrors.Register(ModuleName, 7103, "authorization not accepted")
	BankError                        = sdkerrors.Register(ModuleName, 7104, "bank error")
	ErrExpiredCampaign               = sdkerrors.Register(ModuleName, 7105, "campaign is expired")
	ErrCampaignPoolBalance           = sdkerrors.Register(ModuleName, 7106, "not enough campaign pool balance")
	ErrUnknownRewardType             = sdkerrors.Register(ModuleName, 7107, "unknown reward type")
	ErrInFundingCampaignPool         = sdkerrors.Register(ModuleName, 7108, "error in funding the campaign pool")
	ErrUnknownAccType                = sdkerrors.Register(ModuleName, 7109, "unknown account type")
	ErrAllocationExpired             = sdkerrors.Register(ModuleName, 7110, "allocation expired")
	ErrInsufficientPoolBalance       = sdkerrors.Register(ModuleName, 7111, "insufficient campaign pool balance")
	ErrInDistributionOfRewards       = sdkerrors.Register(ModuleName, 7112, "reward distribution failed")
	ErrInvalidReceiverType           = sdkerrors.Register(ModuleName, 7113, "inappropriate receiver account type")
	ErrWrongDefinitionsCount         = sdkerrors.Register(ModuleName, 7114, "wrong reward definitions")
	ErrMissingDefinition             = sdkerrors.Register(ModuleName, 7115, "missing reward definition")
	ErrSubAccountNotfoundForTheOwner = sdkerrors.Register(ModuleName, 7116, "subaccount not found for the owner address")
)
