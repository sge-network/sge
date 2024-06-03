package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	ovmKeeper     bettypes.OVMKeeper
	betKeeper     types.BetKeeper
	houseKeeper   types.HouseKeeper
	obKeeper      types.OrderBookKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	ovmKeeper bettypes.OVMKeeper,
	betKeeper types.BetKeeper,
	obKeeper types.OrderBookKeeper,
	hk types.HouseKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it is not already set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		ovmKeeper:     ovmKeeper,
		betKeeper:     betKeeper,
		houseKeeper:   hk,
		obKeeper:      obKeeper,
		authority:     authority,
	}
	return k
}

// GetAuthority returns the x/subaccount module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
