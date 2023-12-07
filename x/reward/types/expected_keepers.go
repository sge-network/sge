package types

import (
	context "context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	bettypes "github.com/sge-network/sge/x/bet/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, ecipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// BetKeeper defines the expected interface needed to access bet state.
type BetKeeper interface {
	GetBet(ctx sdk.Context, creator string, id uint64) (val bettypes.Bet, found bool)
	GetBetID(ctx sdk.Context, uid string) (val bettypes.UID2ID, found bool)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// AuthzKeeper defines the expected authz keeper.
type AuthzKeeper interface {
	GetAuthorization(
		ctx sdk.Context,
		grantee sdk.AccAddress,
		granter sdk.AccAddress,
		msgType string,
	) (authz.Authorization, *time.Time)
	SaveGrant(
		ctx sdk.Context,
		grantee, granter sdk.AccAddress,
		authorization authz.Authorization,
		expiration *time.Time,
	) error
	DeleteGrant(
		ctx sdk.Context,
		grantee, granter sdk.AccAddress,
		msgType string,
	) error
}

// SubAccountKeeper defines the expected interface needed to get/create/topup a subaccount.
type SubAccountKeeper interface {
	TopUp(ctx sdk.Context, creator, subAccOwnerAddr string, lockedBalance []subaccounttypes.LockedBalance) (string, error)
	GetSubAccountByOwner(ctx sdk.Context, mainAccountAddress sdk.AccAddress) (sdk.AccAddress, bool)
	CreateSubAccount(ctx sdk.Context, creator, owner string, lockedBalances []subaccounttypes.LockedBalance) (string, error)
	IsSubAccount(ctx sdk.Context, subAccAddr sdk.AccAddress) bool
}

// RewardKeeper defines the expected interface needed to get and filter the rewards.
type RewardKeeper interface {
	HasRewardByReceiver(ctx sdk.Context, addr string, category RewardCategory) bool
}
