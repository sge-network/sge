package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

type MsgRouter interface {
	Handler(msg sdk.Msg) baseapp.MsgServiceHandler
}

type Keeper struct {
	cdc        codec.Codec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace
	bankKeeper bankkeeper.Keeper

	msgRouter MsgRouter
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, ps paramtypes.Subspace, router MsgRouter) Keeper {
	// set KeyTable if it is not already set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{storeKey: storeKey, cdc: cdc, paramstore: ps, msgRouter: router}
}
