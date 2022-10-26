package keeper

import (
	"github.com/sge-network/sge/x/sportevent/types"
)

var _ types.QueryServer = Keeper{}
