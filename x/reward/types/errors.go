package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/reward module sentinel errors
var (
	ErrExpiredCampaign     = sdkerrors.Register(ModuleName, 5100, "campaign is expired")
	ErrCampaignPoolBalance = sdkerrors.Register(ModuleName, 5101, "not enough campaign pool balance")
	ErrUnknownRewardType   = sdkerrors.Register(ModuleName, 5102, "unknown reward type")
)
