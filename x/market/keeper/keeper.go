package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sge-network/sge/x/market/types"
)

// Keeper is the type for module properties
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        sdk.StoreKey
	memKey          sdk.StoreKey
	paramStore      paramtypes.Subspace
	dvmKeeper       types.DVMKeeper
	orderBookKeeper types.OrderBookKeeper
}

// NewKeeper creates new keeper object
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramStore: ps,
	}
}

// SetOrderBookKeeper sets the order book module keeper to the market keeper.
func (k *Keeper) SetOrderBookKeeper(orderBookKeeper types.OrderBookKeeper) {
	k.orderBookKeeper = orderBookKeeper
}

// SetDVMKeeper sets the dvm module keeper to the market keeper.
func (k *Keeper) SetDVMKeeper(dvmKeeper types.DVMKeeper) {
	k.dvmKeeper = dvmKeeper
}

// Logger returns the logger of the keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
