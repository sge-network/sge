package types

import (
	context "context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
)

// BetKeeper defines the expected interface needed to retrieve or set bets.
type BetKeeper interface {
	PrepareBetObject(ctx sdk.Context, creator string, props *bettypes.WagerProps) (*bettypes.Bet, map[string]*bettypes.BetOddsCompact, error)
	Wager(ctx sdk.Context, bet *bettypes.Bet, betOdds map[string]*bettypes.BetOddsCompact) error
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
	ParseDepositTicketAndValidate(goCtx context.Context, ctx sdk.Context, msg *housetypes.MsgDeposit, authzAllowed bool) (string, error)
	Deposit(ctx sdk.Context, creator, depositor string, marketUID string, amount sdkmath.Int) (participationIndex uint64, feeAmount sdkmath.Int, err error)
	ParseWithdrawTicketAndValidate(goCtx context.Context, ctx sdk.Context, msg *housetypes.MsgWithdraw, authzAllowed bool) (string, bool, error)
	CalcAndWithdraw(ctx sdk.Context, msg *housetypes.MsgWithdraw, depositorAddr string, isOnBehalf bool) (uint64, error)
}

// OrderbookKeeper defines the expected interface needed to initiate an order book for a market
type OrderBookKeeper interface {
	CalcWithdrawalAmount(
		ctx sdk.Context,
		depositorAddress string,
		marketUID string,
		participationIndex uint64,
		mode housetypes.WithdrawalMode,
		totalWithdrawnAmount sdkmath.Int,
		amount sdkmath.Int,
	) (sdkmath.Int, error)
}
