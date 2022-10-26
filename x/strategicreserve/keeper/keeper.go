package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// Keeper is the type for module properties
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	paramstore    paramtypes.Subspace
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

// ExpectedKeepers contains expected keepers parameter needed by NewKeeper
type ExpectedKeepers struct {
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
}

// NewKeeper returns a new keeper object
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	expectedKeepers ExpectedKeepers,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// Ensures that the `sr_pool`,  `bet_reserve` and `winnings_collector`
	// module accounts are set
	if addr := expectedKeepers.AccountKeeper.GetModuleAddress(types.SRPoolName); addr == nil {
		panic(fmt.Sprintf(consts.ErrModuleAccountHasNotBeenSet, types.SRPoolName))
	}

	if addr := expectedKeepers.AccountKeeper.GetModuleAddress(types.BetReserveName); addr == nil {
		panic(fmt.Sprintf(consts.ErrModuleAccountHasNotBeenSet, types.BetReserveName))
	}

	if addr := expectedKeepers.AccountKeeper.GetModuleAddress(types.WinningsCollectorName); addr == nil {
		panic(fmt.Sprintf(consts.ErrModuleAccountHasNotBeenSet, types.WinningsCollectorName))
	}

	return &Keeper{

		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		bankKeeper:    expectedKeepers.BankKeeper,
		accountKeeper: expectedKeepers.AccountKeeper,
	}
}

// Logger returns a logger for logging error/debug/info logs
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
