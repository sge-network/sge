package keeper

import (
	"github.com/sge-network/sge/x/mint/types"
)

var _ types.QueryServer = Keeper{}
