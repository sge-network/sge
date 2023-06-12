package types

import (
	context "context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// OrderbookKeeper defines the expected orderbook keeper.
type OrderbookKeeper interface {
	InitiateOrderBookParticipation(ctx sdk.Context, addr sdk.AccAddress,
		bookUID string, liquidity, fee sdk.Int) (uint64, error)
	WithdrawOrderBookParticipation(ctx sdk.Context, depAddr,
		bookUID string, bpNumber uint64, mode WithdrawalMode, amount sdk.Int) (sdk.Int, error)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// AuthzKeeper defines the expected authz keeper.
type AuthzKeeper interface {
	GetCleanAuthorization(
		ctx sdk.Context,
		grantee sdk.AccAddress,
		granter sdk.AccAddress,
		msgType string,
	) (cap authz.Authorization, expiration time.Time)
}
