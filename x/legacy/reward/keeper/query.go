package keeper

import (
	"github.com/sge-network/sge/x/legacy/reward/types"
)

var _ types.QueryServer = Keeper{}
