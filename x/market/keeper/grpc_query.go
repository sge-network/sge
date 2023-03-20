package keeper

import (
	"github.com/sge-network/sge/x/market/types"
)

var _ types.QueryServer = Keeper{}
