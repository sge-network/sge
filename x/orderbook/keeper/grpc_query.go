package keeper

import (
	"github.com/sge-network/sge/x/orderbook/types"
)

var _ types.QueryServer = Keeper{}
