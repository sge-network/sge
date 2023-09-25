package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/reward/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		memKey           storetypes.StoreKey
		paramstore       paramtypes.Subspace
		modFunder        *utils.ModuleAccFunder
		authzKeeper      types.AuthzKeeper
		ovmKeeper        types.OVMKeeper
		subaccountKeeper types.SubAccountKeeper
	}
)

// SdkExpectedKeepers contains expected keepers parameter needed by NewKeeper
type SdkExpectedKeepers struct {
	AuthzKeeper   types.AuthzKeeper
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ovmKeeper types.OVMKeeper,
	subaccountKeeper types.SubAccountKeeper,
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
			types.BankError,
		),
		ovmKeeper:        ovmKeeper,
		subaccountKeeper: subaccountKeeper,
		authzKeeper:      expectedKeepers.AuthzKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
