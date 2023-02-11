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
	GetFirstUnsettledResovedSportEvent(ctx sdk.Context) (string, bool)
	RemoveUnsettledResolvedSportEvent(ctx sdk.Context, sportEventUID string)
}

// DVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type DVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// OrderBookKeeper defines the expected interface needed to process bet placement
type OrderBookKeeper interface {
	ProcessBetPlacement(ctx sdk.Context, uniqueLock, bookId, oddId string, maxLossMultiplier sdk.Dec, payoutProfit sdk.Int, bettorAddress sdk.AccAddress, betFee sdk.Coin, betAmount sdk.Int, oddsType OddsType, oddsVal string) (error, []*BetFullfillment)
	RefundBettor(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string) error
	BettorWins(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string, fullfillment []*BetFullfillment, bookId string) error
	BettorLoses(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string, fullfillment []*BetFullfillment, bookId string) error
	AddBookSettlement(ctx sdk.Context, orderBookId string) error
}
