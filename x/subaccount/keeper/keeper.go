package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

type BetKeeper interface {
	GetBetID(ctx sdk.Context, uid string) (bettypes.UID2ID, bool)
	PlaceBet(ctx sdk.Context, bet *bettypes.Bet) error
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type Keeper struct {
	cdc        codec.Codec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace
	bankKeeper BankKeeper

	ovmKeeper bettypes.OVMKeeper
	betKeeper BetKeeper
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, ps paramtypes.Subspace, bankKeeper BankKeeper, ovmKeeper bettypes.OVMKeeper, betKeeper BetKeeper) Keeper {
	// set KeyTable if it is not already set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
		bankKeeper: bankKeeper,
		ovmKeeper:  ovmKeeper,
		betKeeper:  betKeeper,
	}
}
