package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// BetKeeper defines the expected bet keeper.
type BetKeeper interface {
	GetBetID(ctx sdk.Context, uid string) (val bettypes.UID2ID, found bool)
}
