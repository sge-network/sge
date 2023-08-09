package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	orderbookmodulekeeper "github.com/sge-network/sge/x/orderbook/keeper"
)

type BetKeeper interface {
	GetBetID(ctx sdk.Context, uid string) (bettypes.UID2ID, bool)
	Wager(ctx sdk.Context, bet *bettypes.Bet) error
}

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type HouseKeeper interface {
	GetParams(ctx sdk.Context) housetypes.Params
	Deposit(ctx sdk.Context, creator, depositor, marketUID string, amount sdk.Int) (participationIndex uint64, err error)
	GetDeposit(ctx sdk.Context, depositorAddr, marketUID string, participationIndex uint64) (housetypes.Deposit, bool)
	Withdraw(ctx sdk.Context, deposit housetypes.Deposit, creator, depositorAddr string, marketUID string, participationIndex uint64, mode housetypes.WithdrawalMode, withdrawableAmount sdk.Int) (uint64, error)
}

type OrderBookKeeper interface {
	RegisterHook(hooks orderbookmodulekeeper.Hook)
	CalcWithdrawalAmount(
		ctx sdk.Context,
		depositorAddress string,
		marketUID string,
		participationIndex uint64,
		mode housetypes.WithdrawalMode,
		totalWithdrawnAmount sdk.Int,
		amount sdk.Int,
	) (sdk.Int, error)
}
