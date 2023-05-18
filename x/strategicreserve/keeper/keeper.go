package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/sge-network/sge/x/strategicreserve/types"
)

// keeper of the strategicreserve store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramstore    paramtypes.Subspace
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	BetKeeper     types.BetKeeper
	marketKeeper  types.MarketKeeper
	houseKeeper   types.HouseKeeper
}

// SdkExpectedKeepers contains expected keepers parameter needed by NewKeeper
type SdkExpectedKeepers struct {
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
}

// NewKeeper creates a new strategicreserve Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	ps paramtypes.Subspace,
	expectedKeepers SdkExpectedKeepers,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramstore:    ps,
		bankKeeper:    expectedKeepers.BankKeeper,
		accountKeeper: expectedKeepers.AccountKeeper,
	}
}

// SetBetKeeper sets the bet module keeper to the order book keeper.
func (k *Keeper) SetBetKeeper(betKeeper types.BetKeeper) {
	k.BetKeeper = betKeeper
}

// SetMarketKeeper sets the market module keeper to the order book keeper.
func (k *Keeper) SetMarketKeeper(marketKeeper types.MarketKeeper) {
	k.marketKeeper = marketKeeper
}

// SetMarketKeeper sets the market module keeper to the order book keeper.
func (k *Keeper) SetHouseKeeper(houseKeeper types.HouseKeeper) {
	k.houseKeeper = houseKeeper
}

// Logger returns the logger of the keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
