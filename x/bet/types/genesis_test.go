package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

const (
	testAddress = "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
)

func TestGenesisState_Validate(t *testing.T) {
	sportEventUID := uuid.NewString()
	betUID1, betUID2 := uuid.NewString(), uuid.NewString()
	betID1, betID2 := uint64(1), uint64(2)

	validState := &types.GenesisState{
		BetList: []types.Bet{
			{
				UID:           betUID1,
				Creator:       testAddress,
				SportEventUID: sportEventUID,
			},
			{
				UID:              betUID2,
				Creator:          testAddress,
				SportEventUID:    sportEventUID,
				SettlementHeight: 1,
			},
		},
		ActiveBetList: []types.ActiveBet{
			{
				ID:      betID1,
				Creator: testAddress,
			},
		},
		SettledBetList: []types.SettledBet{
			{
				ID:            betID2,
				BettorAddress: testAddress,
			},
		},
		Uid2IdList: []types.UID2ID{
			{
				UID: betUID1,
				ID:  betID1,
			},
			{
				UID: betUID2,
				ID:  betID2,
			},
		},
		Stats: types.BetStats{
			Count: 2,
		},
	}

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: validState,
			valid:    true,
		},
		{
			desc: "duplicated bet",
			genState: &types.GenesisState{
				BetList: []types.Bet{
					{
						UID: betUID1,
					},
					{
						UID: betUID1,
					},
				},
				ActiveBetList:  validState.ActiveBetList,
				SettledBetList: validState.SettledBetList,
				Uid2IdList:     validState.Uid2IdList,
				Stats:          validState.Stats,
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
