package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.Codec
	storeKey sdk.StoreKey
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc}
}
