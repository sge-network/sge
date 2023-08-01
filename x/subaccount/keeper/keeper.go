package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

type Keeper struct {
	cdc        codec.Codec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace
	bankKeeper types.BankKeeper

	ovmKeeper   bettypes.OVMKeeper
	betKeeper   types.BetKeeper
	houseKeeper types.HouseKeeper
	obKeeper    types.OrderBookKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	ovmKeeper bettypes.OVMKeeper,
	betKeeper types.BetKeeper,
	obKeeper types.OrderBookKeeper,
	hk types.HouseKeeper,
) Keeper {
	// set KeyTable if it is not already set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		paramstore:  ps,
		bankKeeper:  bankKeeper,
		ovmKeeper:   ovmKeeper,
		betKeeper:   betKeeper,
		houseKeeper: hk,
		obKeeper:    obKeeper,
	}
	obKeeper.RegisterHook(k)
	return k
}
