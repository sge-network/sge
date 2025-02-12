package types

import (
	context "context"

	sdkmath "cosmossdk.io/math"
	sdkfeegrant "cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/legacy/bet/types"
	housetypes "github.com/sge-network/sge/x/legacy/house/types"
	markettypes "github.com/sge-network/sge/x/legacy/market/types"
)

// AccountKeeper defines the expected account keeper methods.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper methods.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(
		ctx context.Context,
		senderAddr sdk.AccAddress,
		recipientModule string,
		amt sdk.Coins,
	) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToModule(
		ctx context.Context,
		senderModule, recipientModule string,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToAccount(
		ctx context.Context,
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
		ctx context.Context,
		granter, grantee sdk.AccAddress,
		feeAllowance sdkfeegrant.FeeAllowanceI,
	) error
	GetAllowance(ctx context.Context, granter, grantee sdk.AccAddress) (sdkfeegrant.FeeAllowanceI, error)
}

// Event Hooks
// These can be utilized to communicate between a orderbook keeper and another
// keepers.

// OrderBookHooks event hooks for orderbook methods.
type OrderBookHooks interface {
	AfterHouseWin(ctx sdk.Context, house sdk.AccAddress, originalAmount, profit sdkmath.Int)
	AfterHouseLoss(ctx sdk.Context, house sdk.AccAddress, originalAmount, lostAmt sdkmath.Int)
	AfterHouseRefund(ctx sdk.Context, house sdk.AccAddress, originalAmount sdkmath.Int)
	AfterHouseFeeRefund(ctx sdk.Context, house sdk.AccAddress, fee sdkmath.Int)
}
