package keeper

import (
	"github.com/sge-network/sge/x/legacy/market/types"
)

var _ types.QueryServer = Keeper{}
