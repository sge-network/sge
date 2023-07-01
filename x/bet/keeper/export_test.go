package keeper

import (
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

func (KeeperTest) ProcessBetResultAndStatus(bet *types.Bet, market markettypes.Market) error {
	return processBetResultAndStatus(bet, market)
}

func (KeeperTest) CheckBetStatus(bet *types.Bet) error {
	return checkBetStatus(bet.Status)
}
