package keeper

import (
	"github.com/sge-network/sge/x/house/types"
)

var _ types.QueryServer = Keeper{}
