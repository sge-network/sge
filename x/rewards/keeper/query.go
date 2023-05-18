package keeper

import (
	"github.com/sge-network/sge/x/rewards/types"
)

var _ types.QueryServer = Keeper{}
