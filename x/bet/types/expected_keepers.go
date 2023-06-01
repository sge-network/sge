package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	markettypes "github.com/sge-network/sge/x/market/types"
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

// MarketKeeper defines the expected interface needed to get markets from KVStore
type MarketKeeper interface {
	GetMarket(ctx sdk.Context, marketUID string) (markettypes.Market, bool)
	GetFirstUnsettledResolvedMarket(ctx sdk.Context) (string, bool)
	GetDefaultBetConstraints(ctx sdk.Context) (params *markettypes.MarketBetConstraints)
	RemoveUnsettledResolvedMarket(ctx sdk.Context, marketUID string)
}

// OVMKeeper defines the expected interface needed to verify ticket and unmarshal it
type OVMKeeper interface {
	VerifyTicketUnmarshal(goCtx context.Context, ticket string, clm interface{}) error
}

// SRKeeper defines the expected interface needed to process bet placement
type SRKeeper interface {
	ProcessBetPlacement(ctx sdk.Context, betUID, bookUID, oddsUID string, maxLossMultiplier sdk.Dec, betAmount sdk.Int, payoutProfit sdk.Dec, bettorAddress sdk.AccAddress, betFee sdk.Int, oddsType OddsType, oddsVal string, betID uint64) ([]*BetFulfillment, error)
	RefundBettor(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, betFee, payout sdk.Int, uniqueLock string) error
	BettorWins(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string, fulfillment []*BetFulfillment, bookUID string) error
	BettorLoses(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string, fulfillment []*BetFulfillment, bookUID string) error
	SetOrderBookAsUnsettledResolved(ctx sdk.Context, orderBookUID string) error
	WithdrawBetFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee sdk.Int) error
}
