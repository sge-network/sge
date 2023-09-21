package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/bet/types"
)

const (
	testAddress = "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
)

func TestGenesisState_Validate(t *testing.T) {
	marketUID := uuid.NewString()
	betUID1, betUID2 := uuid.NewString(), uuid.NewString()
	betID1, betID2 := uint64(1), uint64(2)

	validState := &types.GenesisState{
		BetList: []types.Bet{
			{
				UID:       betUID1,
				Creator:   testAddress,
				MarketUID: marketUID,
			},
			{
				UID:              betUID2,
				Creator:          testAddress,
				MarketUID:        marketUID,
				SettlementHeight: 1,
			},
		},
		PendingBetList: []types.PendingBet{
			{
				UID:     betUID1,
				Creator: testAddress,
			},
		},
		SettledBetList: []types.SettledBet{
			{
				UID:           betUID2,
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
		Params: types.DefaultParams(),
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
				PendingBetList: validState.PendingBetList,
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
