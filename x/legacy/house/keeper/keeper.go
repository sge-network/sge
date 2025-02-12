package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/x/legacy/house/types"
)

// Keeper of the house store
type Keeper struct {
	storeKey        storetypes.StoreKey
	cdc             codec.BinaryCodec
	paramstore      paramtypes.Subspace
	authzKeeper     types.AuthzKeeper
	orderbookKeeper types.OrderbookKeeper
	ovmKeeper       types.OVMKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// SdkExpectedKeepers contains expected keepers parameter needed by NewKeeper
type SdkExpectedKeepers struct {
	AuthzKeeper types.AuthzKeeper
}

// NewKeeper returns an instance of the housekeeper
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	orderbookKeeper types.OrderbookKeeper,
	ovmKeeper types.OVMKeeper,
	ps paramtypes.Subspace,
	expectedKeepers SdkExpectedKeepers,
	authority string,
) *Keeper {
	// set KeyTable if it is not already set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:        key,
		cdc:             cdc,
		orderbookKeeper: orderbookKeeper,
		ovmKeeper:       ovmKeeper,
		paramstore:      ps,
		authzKeeper:     expectedKeepers.AuthzKeeper,
		authority:       authority,
	}
}

// GetAuthority returns the x/house module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
