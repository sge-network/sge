package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/x/legacy/bet/types"
)

// Keeper is the type for module properties
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        storetypes.StoreKey
	memKey          storetypes.StoreKey
	paramstore      paramtypes.Subspace
	marketKeeper    types.MarketKeeper
	orderbookKeeper types.OrderbookKeeper
	ovmKeeper       types.OVMKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates new keeper object
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authority string,
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
		authority:  authority,
	}
}

// SetMarketKeeper sets market keeper to the bet keeper.
func (k *Keeper) SetMarketKeeper(marketKeeper types.MarketKeeper) {
	k.marketKeeper = marketKeeper
}

// SetOrderbookKeeper sets orderbook keeper to the bet keeper.
func (k *Keeper) SetOrderbookKeeper(orderbookKeeper types.OrderbookKeeper) {
	k.orderbookKeeper = orderbookKeeper
}

// SetOVMKeeper sets ovm keeper to the bet keeper.
func (k *Keeper) SetOVMKeeper(ovmKeeper types.OVMKeeper) {
	k.ovmKeeper = ovmKeeper
}

// Logger returns the logger of the keeper
func (Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the x/bet module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
