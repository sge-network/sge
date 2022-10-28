package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// SporteventKeeper defines the expected interface needed to get sportEvents from KVStore
type SporteventKeeper interface {
	GetSportEvent(ctx sdk.Context, sportEventUID string) (sporteventtypes.SportEvent, bool)
	AddExtraPayoutToEvent(ctx sdk.Context, sportEventUID string, amount sdk.Int) error
}

// StrategicreserveKeeper defines the expected interface needed to unlock fund and pay out
type StrategicreserveKeeper interface {
	ProcessBetPlacement(ctx sdk.Context, bettorAddress sdk.AccAddress,
		betFee sdk.Coin, betAmount sdk.Int, extraPayout sdk.Int,
		uniqueLock string) error

	BettorWins(ctx sdk.Context, bettorAddress sdk.AccAddress,
		betAmount sdk.Int, extraPayout sdk.Int, uniqueLock string) error

	BettorLoses(ctx sdk.Context, address sdk.AccAddress,
		betAmount sdk.Int, extraPayout sdk.Int, uniqueLock string) error

	RefundBettor(ctx sdk.Context, bettorAddress sdk.AccAddress,
		betAmount sdk.Int, extraPayout sdk.Int, uniqueLock string) error
}

// DVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type DVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}
