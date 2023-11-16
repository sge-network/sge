package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/reward module sentinel errors
var (
	ErrInTicketVerification       = sdkerrors.Register(ModuleName, 7100, "ticket verification failed")
	ErrInTicketValidation         = sdkerrors.Register(ModuleName, 7101, "ticket validation failed")
	ErrAuthorizationNotFound      = sdkerrors.Register(ModuleName, 7102, "no authorization found")
	ErrAuthorizationNotAccepted   = sdkerrors.Register(ModuleName, 7103, "authorization not accepted")
	ErrorBank                     = sdkerrors.Register(ModuleName, 7104, "bank error")
	ErrExpiredCampaign            = sdkerrors.Register(ModuleName, 7105, "campaign is expired")
	ErrCampaignPoolBalance        = sdkerrors.Register(ModuleName, 7106, "not enough campaign pool balance")
	ErrUnknownRewardType          = sdkerrors.Register(ModuleName, 7107, "unknown reward type")
	ErrInFundingCampaignPool      = sdkerrors.Register(ModuleName, 7108, "error in funding the campaign pool")
	ErrWithdrawFromCampaignPool   = sdkerrors.Register(ModuleName, 7109, "error in withdrawing from the campaign pool")
	ErrCampaignEnded              = sdkerrors.Register(ModuleName, 7110, "campaign validity period is ended")
	ErrInsufficientPoolBalance    = sdkerrors.Register(ModuleName, 7111, "insufficient campaign pool balance")
	ErrInDistributionOfRewards    = sdkerrors.Register(ModuleName, 7112, "reward distribution failed")
	ErrInvalidGranteeType         = sdkerrors.Register(ModuleName, 7113, "inappropriate reward grantee account type")
	ErrWrongRewardCategory        = sdkerrors.Register(ModuleName, 7114, "wrong reward category")
	ErrMissingDefinition          = sdkerrors.Register(ModuleName, 7115, "missing reward definition")
	ErrSubAccRewardTopUp          = sdkerrors.Register(ModuleName, 7116, "subaccount reward topup failed")
	ErrUnlockTSIsSubAccOnly       = sdkerrors.Register(ModuleName, 7117, "unlock timestamp is allowed for subaccount only")
	ErrUnlockTSDefBeforeBlockTime = sdkerrors.Register(ModuleName, 7118, "unlock timestamp should not be before the current block time")
	ErrInvalidNoLossBetUID        = sdkerrors.Register(ModuleName, 7120, "invalid no loss bet uid")
	ErrWrongAmountForType         = sdkerrors.Register(ModuleName, 7121, "wrong amount for account type")
	ErrSubAccountCreationFailed   = sdkerrors.Register(ModuleName, 7122, "subaccount creation failed")
	ErrWrongRewardAmountType      = sdkerrors.Register(ModuleName, 7123, "wrong reward amount type")
	ErrUnknownAccType             = sdkerrors.Register(ModuleName, 7124, "unknown account type")
	ErrReceiverAddrCanNotBeSubAcc = sdkerrors.Register(ModuleName, 7125, "receiver account can not be sub account address")
	ErrInvalidFunds               = sdkerrors.Register(ModuleName, 7126, "invalid funds")
)
