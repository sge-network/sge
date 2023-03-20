package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sge-network/sge/x/bet/types"
)

// Keeper is the type for module properties
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        sdk.StoreKey
	memKey          sdk.StoreKey
	paramstore      paramtypes.Subspace
	marketKeeper    types.MarketKeeper
	orderbookKeeper types.OrderBookKeeper
	dvmKeeper       types.DVMKeeper
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
		paramstore: ps,
	}
}

// SetMarketKeeper sets market keeper to the bet keeper.
func (k *Keeper) SetMarketKeeper(marketKeeper types.MarketKeeper) {
	k.marketKeeper = marketKeeper
}

// SetOrderBookKeeper sets order book keeper to the bet keeper.
func (k *Keeper) SetOrderBookKeeper(orderBookKeeper types.OrderBookKeeper) {
	k.orderbookKeeper = orderBookKeeper
}

// SetDVMKeeper sets dvm keeper to the bet keeper.
func (k *Keeper) SetDVMKeeper(dvmKeeper types.DVMKeeper) {
	k.dvmKeeper = dvmKeeper
}

// Logger returns the logger of the keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
