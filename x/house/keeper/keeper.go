package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/x/house/types"
)

// keeper of the house store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	srKeeper   types.SRKeeper
	paramstore paramtypes.Subspace
}

// NewKeeper creates a new house Keeper instance
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, srKeeper types.SRKeeper, ps paramtypes.Subspace) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:   key,
		cdc:        cdc,
		srKeeper:   srKeeper,
		paramstore: ps,
	}
}
