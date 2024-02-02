package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/bet/types"
)

// Keeper is the type for module properties
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        storetypes.StoreKey
	memKey          storetypes.StoreKey
	paramstore      paramtypes.Subspace
	modFunder       *utils.ModuleAccFunder
	marketKeeper    types.MarketKeeper
	orderbookKeeper types.OrderbookKeeper
	ovmKeeper       types.OVMKeeper
}

// SdkExpectedKeepers contains expected keepers parameter needed by NewKeeper
type SdkExpectedKeepers struct {
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
}

// NewKeeper creates new keeper object
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	expectedKeepers SdkExpectedKeepers,
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
		modFunder: utils.NewModuleAccFunder(
			expectedKeepers.BankKeeper,
			expectedKeepers.AccountKeeper,
			types.ErrorBank,
		),
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
