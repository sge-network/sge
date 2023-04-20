package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SRKeeper defines the expected strategicreserve keeper.
type SRKeeper interface {
	InitiateOrderBookParticipation(ctx sdk.Context, addr sdk.AccAddress,
		bookUID string, liquidity, fee sdk.Int) (uint64, error)
	WithdrawOrderBookParticipation(ctx sdk.Context, depAddr,
		bookUID string, bpNumber uint64, mode WithdrawalMode, amount sdk.Int) (sdk.Int, error)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}
