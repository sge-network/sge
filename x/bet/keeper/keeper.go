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
	cdc                    codec.BinaryCodec
	storeKey               sdk.StoreKey
	memKey                 sdk.StoreKey
	paramstore             paramtypes.Subspace
	sporteventKeeper       types.SporteventKeeper
	strategicreserveKeeper types.StrategicreserveKeeper
	dvmKeeper              types.DVMKeeper
}

// ExpectedKeepers contains expected keepers parameter needed by NewKeeper
type ExpectedKeepers struct {
	SporteventKeeper       types.SporteventKeeper
	StrategicreserveKeeper types.StrategicreserveKeeper
	DVMKeeper              types.DVMKeeper
}

// NewKeeper creates new keeper object
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

	return &Keeper{
		cdc:                    cdc,
		storeKey:               storeKey,
		memKey:                 memKey,
		paramstore:             ps,
		sporteventKeeper:       expectedKeepers.SporteventKeeper,
		strategicreserveKeeper: expectedKeepers.StrategicreserveKeeper,
		dvmKeeper:              expectedKeepers.DVMKeeper,
	}
}

// Logger retuns the logger of the keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
