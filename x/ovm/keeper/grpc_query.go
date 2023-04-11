package keeper

import (
	"github.com/sge-network/sge/x/ovm/types"
)

var _ types.QueryServer = Keeper{}
