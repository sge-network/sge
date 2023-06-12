package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

// AccountKeeper defines the expected account keeper methods.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper methods.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(
		ctx sdk.Context,
		senderAddr sdk.AccAddress,
		recipientModule string,
		amt sdk.Coins,
	) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToModule(
		ctx sdk.Context,
		senderModule, recipientModule string,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToAccount(
		ctx sdk.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
}

// BetKeeper defines the expected bet keeper methods.
type BetKeeper interface {
	GetBetID(ctx sdk.Context, uid string) (val bettypes.UID2ID, found bool)
}

// MarketKeeper defines the expected market keeper methods.
type MarketKeeper interface {
	GetMarket(ctx sdk.Context, marketUID string) (val markettypes.Market, found bool)
}

// HouseKeeper defines the expected market keeper methods.
type HouseKeeper interface {
	GetDeposit(
		ctx sdk.Context,
		depositorAddress, marketUID string,
		participationIndex uint64,
	) (val housetypes.Deposit, found bool)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// FeeGrantKeeper defines the expected interface needed for the fee grant.
type FeeGrantKeeper interface {
	GrantAllowance(
		ctx sdk.Context,
		granter, grantee sdk.AccAddress,
		feeAllowance sdkfeegrant.FeeAllowanceI,
	) error
	GetAllowance(ctx sdk.Context, granter, grantee sdk.AccAddress) (sdkfeegrant.FeeAllowanceI, error)
}
