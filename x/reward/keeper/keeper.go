package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/cometbft/cometbft/libs/log"

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
		accountKeeper    types.AccountKeeper
		authzKeeper      types.AuthzKeeper
		betKeeper        types.BetKeeper
		ovmKeeper        types.OVMKeeper
		subaccountKeeper types.SubaccountKeeper
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
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
	betKeeper types.BetKeeper,
	ovmKeeper types.OVMKeeper,
	subaccountKeeper types.SubaccountKeeper,
	expectedKeepers SdkExpectedKeepers,
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
		modFunder: utils.NewModuleAccFunder(
			expectedKeepers.BankKeeper,
			expectedKeepers.AccountKeeper,
			types.ErrorBank,
		),
		accountKeeper:    expectedKeepers.AccountKeeper,
		betKeeper:        betKeeper,
		ovmKeeper:        ovmKeeper,
		subaccountKeeper: subaccountKeeper,
		authzKeeper:      expectedKeepers.AuthzKeeper,
		authority:        authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the x/reward module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
