package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// Keeper of the orderbook store
type Keeper struct {
	storeKey       storetypes.StoreKey
	cdc            codec.BinaryCodec
	paramstore     paramtypes.Subspace
	bankKeeper     types.BankKeeper
	accountKeeper  types.AccountKeeper
	BetKeeper      types.BetKeeper
	marketKeeper   types.MarketKeeper
	houseKeeper    types.HouseKeeper
	ovmKeeper      types.OVMKeeper
	feeGrantKeeper types.FeeGrantKeeper
	hooks          types.OrderBookHooks
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// SdkExpectedKeepers contains expected keepers parameter needed by NewKeeper
type SdkExpectedKeepers struct {
	BankKeeper     types.BankKeeper
	AccountKeeper  types.AccountKeeper
	FeeGrantKeeper types.FeeGrantKeeper
}

// NewKeeper creates a new orderbook Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ps paramtypes.Subspace,
	expectedKeepers SdkExpectedKeepers,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:       key,
		cdc:            cdc,
		paramstore:     ps,
		bankKeeper:     expectedKeepers.BankKeeper,
		accountKeeper:  expectedKeepers.AccountKeeper,
		feeGrantKeeper: expectedKeepers.FeeGrantKeeper,
		authority:      authority,
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

// SetHouseKeeper sets the market module keeper to the order book keeper.
func (k *Keeper) SetHouseKeeper(houseKeeper types.HouseKeeper) {
	k.houseKeeper = houseKeeper
}

// SetOVMKeeper sets the ovm module keeper to the market keeper.
func (k *Keeper) SetOVMKeeper(ovmKeeper types.OVMKeeper) {
	k.ovmKeeper = ovmKeeper
}

// SetHooks sets the hooks for governance
func (k *Keeper) SetHooks(gh types.OrderBookHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set orderbook hooks twice")
	}

	k.hooks = gh

	return k
}

// Logger returns the logger of the keeper
func (Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the x/orderbook module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
