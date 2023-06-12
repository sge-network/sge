package types

import (
	"fmt"

	"github.com/sge-network/sge/utils"
)

// DefaultUID is the default  global uid
const DefaultUID uint64 = 1

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		BetList:        []Bet{},
		PendingBetList: []PendingBet{},
		SettledBetList: []SettledBet{},
		Uid2IdList:     []UID2ID{},
		Stats:          BetStats{},
		Params:         DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	betCount := uint64(len(gs.BetList))
	if betCount != gs.Stats.Count {
		return fmt.Errorf(
			"%s: %d <> %d",
			ErrTextInitGenesisFailedBecauseOfNotEqualStats,
			betCount,
			gs.Stats.Count,
		)
	}

	activeAndSettledCount := uint64(len(gs.PendingBetList)) + uint64(len(gs.SettledBetList))
	if activeAndSettledCount != betCount {
		return fmt.Errorf(
			"%s: %d <> %d",
			ErrTextInitGenesisFailedBetCountNotEqualActiveAndSettled,
			activeAndSettledCount,
			betCount,
		)
	}

	// Check for duplicated uid in bet
	betUIDMap := make(map[string]struct{})

	for _, elem := range gs.BetList {
		uid := string(utils.StrBytes(elem.UID))
		if _, ok := betUIDMap[uid]; ok {
			return fmt.Errorf("duplicated uid for bet")
		}
		betUIDMap[uid] = struct{}{}
	}

	// Set all the bets
	for _, bet := range gs.BetList {
		var id uint64
		for _, uid2ID := range gs.Uid2IdList {
			if uid2ID.UID == bet.UID {
				id = uid2ID.ID
			}
		}

		if id == 0 {
			// this means the imported genesis is broken because there is no corresponding
			// id mapped to the uid
			return fmt.Errorf("%s: %s", ErrTextInitGenesisFailedBecauseOfMissingBetID, bet.UID)
		}

		if bet.SettlementHeight == 0 && bet.Status == Bet_STATUS_SETTLED {
			return fmt.Errorf("%s: %s", ErrTextInitGenesisFailedSettlementHeightIsZero, bet.UID)
		}

		isActive := false
		for _, active := range gs.PendingBetList {
			if active.UID == bet.UID {
				if bet.SettlementHeight != 0 {
					return fmt.Errorf(
						"%s: %s",
						ErrTextInitGenesisFailedSettlementHeightIsNotZero,
						bet.UID,
					)
				}
				isActive = true
			}
		}

		isSettled := false
		for _, settled := range gs.SettledBetList {
			if settled.UID == bet.UID {
				if bet.SettlementHeight == 0 {
					return fmt.Errorf(
						"%s: %s",
						ErrTextInitGenesisFailedSettlementHeightIsZeroForList,
						bet.UID,
					)
				}
				isSettled = true
			}
		}

		if !isActive && !isSettled {
			return fmt.Errorf("%s: %s", ErrTextInitGenesisFailedNotActiveOrSettled, bet.UID)
		}
	}

	return gs.Params.Validate()
}
