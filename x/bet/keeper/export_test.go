package keeper

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

// func (k KeeperTest) ResolveBetResult(bet *types.Bet, sportEvent sporteventtypes.SportEvent) error {
// 	return resolveBetResult(bet, sportEvent)
// }

// func (k KeeperTest) CheckBetStatus(bet *types.Bet) error {
// 	return checkBetStatus(bet.Status)
// }
