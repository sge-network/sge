package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/stretchr/testify/require"
)

const (
	testAddress = "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
)

func TestGenesisState_Validate(t *testing.T) {
	marketUID := uuid.NewString()
	oddsUID := uuid.NewString()
	validState := types.GenesisState{
		OrderBookList: []types.OrderBook{
			{
				UID:                marketUID,
				ParticipationCount: 1,
				OddsCount:          1,
				Status:             types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE,
			},
		},
		OrderBookExposureList: []types.OrderBookOddsExposure{
			{
				OrderBookUID:     marketUID,
				OddsUID:          oddsUID,
				FulfillmentQueue: []uint64{},
			},
		},
		OrderBookParticipationList: []types.OrderBookParticipation{
			{
				OrderBookUID:       marketUID,
				Index:              1,
				ParticipantAddress: testAddress,
			},
		},
		ParticipationExposureList: []types.ParticipationExposure{
			{
				OrderBookUID:       marketUID,
				OddsUID:            oddsUID,
				ParticipationIndex: 1,
			},
		},
		ParticipationExposureByIndexList: []types.ParticipationExposure{
			{
				OrderBookUID:       marketUID,
				OddsUID:            oddsUID,
				ParticipationIndex: 1,
			},
		},
		HistoricalParticipationExposureList: []types.ParticipationExposure{
			{
				OrderBookUID:       marketUID,
				OddsUID:            oddsUID,
				ParticipationIndex: 1,
			},
		},
		ParticipationBetPairExposureList: []types.ParticipationBetPair{
			{
				OrderBookUID:       marketUID,
				ParticipationIndex: 1,
			},
		},
		Stats: types.OrderBookStats{
			ResolvedUnsettled: []string{marketUID},
		},
		Params: types.DefaultParams(),
	}

	invalidParticipantAddress := validState
	invalidParticipantAddress.OrderBookParticipationList = []types.OrderBookParticipation{
		validState.OrderBookParticipationList[0],
	}
	invalidParticipantAddress.OrderBookParticipationList[0].ParticipantAddress = "wrong"

	notEqualOrderBookCount := validState
	notEqualOrderBookCount.OrderBookList = []types.OrderBook{}

	notEqualOrderBookExposureCount := validState
	notEqualOrderBookExposureCount.OrderBookExposureList = []types.OrderBookOddsExposure{}

	notEqualParticipationExposureCount := validState
	notEqualParticipationExposureCount.ParticipationExposureList = []types.ParticipationExposure{}

	notEqualParticipationExposureIndexCount := validState
	notEqualParticipationExposureIndexCount.ParticipationExposureByIndexList = []types.ParticipationExposure{}

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
			genState: &validState,
			valid:    true,
		},
		{
			desc:     "invalid participant address",
			genState: &invalidParticipantAddress,
			valid:    false,
		},
		{
			desc:     "not equal book",
			genState: &notEqualOrderBookCount,
			valid:    false,
		},
		{
			desc:     "not equal book exposure",
			genState: &notEqualOrderBookExposureCount,
			valid:    false,
		},
		{
			desc:     "not equal participation exposure",
			genState: &notEqualParticipationExposureCount,
			valid:    false,
		},
		{
			desc:     "not equal participation exposure index",
			genState: &notEqualParticipationExposureIndexCount,
			valid:    false,
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
