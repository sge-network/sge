package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	orderbookmodulekeeper "github.com/sge-network/sge/x/orderbook/keeper"
)

// BetKeeper defines the expected interface needed to retrieve or set bets.
type BetKeeper interface {
	GetBetID(ctx sdk.Context, uid string) (bettypes.UID2ID, bool)
	Wager(ctx sdk.Context, bet *bettypes.Bet) error
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// HouseKeeper defines the expected interface needed to deposit or withdraw.
type HouseKeeper interface {
	GetParams(ctx sdk.Context) housetypes.Params
	Deposit(ctx sdk.Context, creator, depositor, marketUID string, amount math.Int) (participationIndex uint64, err error)
	GetDeposit(ctx sdk.Context, depositorAddr, marketUID string, participationIndex uint64) (housetypes.Deposit, bool)
	Withdraw(ctx sdk.Context, deposit housetypes.Deposit, creator, depositorAddr, marketUID string, participationIndex uint64, mode housetypes.WithdrawalMode, withdrawableAmount math.Int) (uint64, error)
}

// OrderbookKeeper defines the expected interface needed to initiate an order book for a market
type OrderBookKeeper interface {
	RegisterHook(hooks orderbookmodulekeeper.Hook)
	CalcWithdrawalAmount(
		ctx sdk.Context,
		depositorAddress string,
		marketUID string,
		participationIndex uint64,
		mode housetypes.WithdrawalMode,
		totalWithdrawnAmount math.Int,
		amount math.Int,
	) (math.Int, error)
}
