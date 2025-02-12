package keeper

import (
	"github.com/sge-network/sge/x/legacy/house/types"
)

var _ types.QueryServer = Keeper{}
