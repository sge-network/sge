package types

import (
	context "context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// OrderbookKeeper defines the expected orderbook keeper.
type OrderbookKeeper interface {
	InitiateOrderBookParticipation(ctx sdk.Context, addr sdk.AccAddress, bookUID string,
		liquidity, fee sdk.Int,
	) (uint64, error)
	CalcWithdrawalAmount(ctx sdk.Context, depositorAddress string, marketUID string,
		participationIndex uint64, mode WithdrawalMode, totalWithdrawnAmount, amount sdk.Int,
	) (sdk.Int, error)
	WithdrawOrderBookParticipation(ctx sdk.Context, marketUID string,
		participationIndex uint64, amount sdk.Int,
	) error
	PublishOrderBookEvent(ctx sdk.Context, orderBookUID string)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// AuthzKeeper defines the expected authz keeper.
type AuthzKeeper interface {
	GetCleanAuthorization(
		ctx sdk.Context,
		grantee, granter sdk.AccAddress,
		msgType string,
	) (cap authz.Authorization, expiration time.Time)
	SaveGrant(
		ctx sdk.Context,
		grantee, granter sdk.AccAddress,
		authorization authz.Authorization,
		expiration time.Time,
	) error
	DeleteGrant(
		ctx sdk.Context,
		grantee, granter sdk.AccAddress,
		msgType string,
	) error
}
