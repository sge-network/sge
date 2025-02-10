package types

import (
	context "context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/authz"

	bettypes "github.com/sge-network/sge/x/legacy/bet/types"
	markettypes "github.com/sge-network/sge/x/legacy/market/types"
	subaccounttypes "github.com/sge-network/sge/x/legacy/subaccount/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, ecipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// BetKeeper defines the expected interface needed to access bet state.
type BetKeeper interface {
	GetBet(ctx sdk.Context, creator string, id uint64) (val bettypes.Bet, found bool)
	GetBetID(ctx sdk.Context, uid string) (val bettypes.UID2ID, found bool)
}

// BetKeeper defines the expected interface needed to access market state.
type MarketKeeper interface {
	GetMarket(ctx sdk.Context, marketUID string) (markettypes.Market, bool)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// AuthzKeeper defines the expected authz keeper.
type AuthzKeeper interface {
	GetAuthorization(
		ctx context.Context,
		grantee sdk.AccAddress,
		granter sdk.AccAddress,
		msgType string,
	) (authz.Authorization, *time.Time)
	SaveGrant(
		ctx context.Context,
		grantee, granter sdk.AccAddress,
		authorization authz.Authorization,
		expiration *time.Time,
	) error
	DeleteGrant(
		ctx context.Context,
		grantee, granter sdk.AccAddress,
		msgType string,
	) error
}

// SubaccountKeeper defines the expected interface needed to get/create/topup a subaccount.
type SubaccountKeeper interface {
	TopUp(ctx sdk.Context, creator, subAccOwnerAddr string, lockedBalance []subaccounttypes.LockedBalance) (string, error)
	GetSubaccountByOwner(ctx sdk.Context, mainAccountAddress sdk.AccAddress) (sdk.AccAddress, bool)
	CreateSubaccount(ctx sdk.Context, creator, owner string, lockedBalances []subaccounttypes.LockedBalance) (string, error)
	IsSubaccount(ctx sdk.Context, subAccAddr sdk.AccAddress) bool
}

// RewardKeeper defines the expected interface needed to get and filter the rewards.
type RewardKeeper interface {
	HasRewardOfReceiverByPromoter(ctx sdk.Context, promoterUID, addr string, category RewardCategory) bool
	GetPromoterByAddress(ctx sdk.Context, address string) (val PromoterByAddress, found bool)
}
