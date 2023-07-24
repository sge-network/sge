package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	orderbookmodulekeeper "github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

type BetKeeper interface {
	GetBetID(ctx sdk.Context, uid string) (bettypes.UID2ID, bool)
	PlaceBet(ctx sdk.Context, bet *bettypes.Bet) error
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type HouseKeeper interface {
	GetParams(ctx sdk.Context) housetypes.Params
	Deposit(ctx sdk.Context, creator, depositor, marketUID string, amount sdk.Int) (participationIndex uint64, err error)
}

type OrderBookKeeper interface {
	RegisterHook(hooks orderbookmodulekeeper.Hook)
}

type Keeper struct {
	cdc        codec.Codec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace
	bankKeeper BankKeeper

	ovmKeeper   bettypes.OVMKeeper
	betKeeper   BetKeeper
	houseKeeper HouseKeeper
}

func (k Keeper) AfterHouseWin(ctx sdk.Context, house sdk.AccAddress, originalAmount, profit sdk.Int, fee *sdk.Int) {
	// update balance
	balance, exists := k.GetBalance(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	if fee != nil {
		err = balance.Unspend(*fee)
		if err != nil {
			panic(err)
		}
	}
	k.SetBalance(ctx, house, balance)

	// send profits
	subAccountOwner, exists := k.GetSubAccountOwner(ctx, house)
	if !exists {
		panic("data corruption: subaccount owner not found")
	}
	err = k.bankKeeper.SendCoins(ctx, house, subAccountOwner, sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, profit)))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) AfterHouseLoss(ctx sdk.Context, house sdk.AccAddress, originalAmount sdk.Int, lostAmt sdk.Int, fee *sdk.Int) {
	balance, exists := k.GetBalance(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	err = balance.AddLoss(lostAmt)
	if err != nil {
		panic(err)
	}
	if fee != nil {
		err = balance.Unspend(*fee)
		if err != nil {
			panic(err)
		}
	}

	k.SetBalance(ctx, house, balance)

	// send profits
	profits := originalAmount.Sub(lostAmt)
	if !profits.IsPositive() {
		return
	}
	subAccountOwner, exists := k.GetSubAccountOwner(ctx, house)
	if !exists {
		panic("data corruption: subaccount owner not found")
	}
	err = k.bankKeeper.SendCoins(ctx, house, subAccountOwner, sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, profits)))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) AfterHouseRefund(ctx sdk.Context, house sdk.AccAddress, originalAmount, fee sdk.Int) {
	balance, exists := k.GetBalance(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	err = balance.Unspend(fee)
	if err != nil {
		panic(err)
	}

	k.SetBalance(ctx, house, balance)
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, ps paramtypes.Subspace, bankKeeper BankKeeper, ovmKeeper bettypes.OVMKeeper, betKeeper BetKeeper, obKeeper OrderBookKeeper) Keeper {
	// set KeyTable if it is not already set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
		bankKeeper: bankKeeper,
		ovmKeeper:  ovmKeeper,
		betKeeper:  betKeeper,
	}
	obKeeper.RegisterHook(k)
	return k
}
