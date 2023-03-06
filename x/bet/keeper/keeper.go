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
	cdc              codec.BinaryCodec
	storeKey         sdk.StoreKey
	memKey           sdk.StoreKey
	paramstore       paramtypes.Subspace
	sporteventKeeper types.SporteventKeeper
	orderbookKeeper  types.OrderBookKeeper
	dvmKeeper        types.DVMKeeper
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

func (k *Keeper) SetSportEventKeeper(sportEventKeeper types.SporteventKeeper) {
	k.sporteventKeeper = sportEventKeeper
}

func (k *Keeper) SetOrderBookKeeper(orderBookKeeper types.OrderBookKeeper) {
	k.orderbookKeeper = orderBookKeeper
}

func (k *Keeper) SetDVMKeeper(dvmKeeper types.DVMKeeper) {
	k.dvmKeeper = dvmKeeper
}

// Logger returns the logger of the keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
