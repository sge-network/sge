package keeper

import (
	"github.com/sge-network/sge/x/dvm/types"
)

var _ types.QueryServer = Keeper{}
