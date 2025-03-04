package keeper

import (
	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

var _ types.QueryServer = Keeper{}
