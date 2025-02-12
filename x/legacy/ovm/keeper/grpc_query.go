package keeper

import (
	"github.com/sge-network/sge/x/legacy/ovm/types"
)

var _ types.QueryServer = Keeper{}
