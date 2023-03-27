package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/x/house/types"
	"github.com/stretchr/testify/require"
)

const (
	testAddress = "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
)

func TestGenesisState_Validate(t *testing.T) {
	MarketUID := uuid.NewString()
	validState := types.GenesisState{
		DepositList: []types.Deposit{
			{
				Creator:               testAddress,
				MarketUID:             MarketUID,
				ParticipationIndex:    1,
				Amount:                sdk.NewInt(10),
				Fee:                   sdk.NewInt(1),
				Liquidity:             sdk.NewInt(9),
				WithdrawalCount:       1,
				TotalWithdrawalAmount: sdk.NewInt(10),
			},
		},
		WithdrawalList: []types.Withdrawal{
			{
				ID:                 1,
				Address:            testAddress,
				MarketUID:          MarketUID,
				ParticipationIndex: 1,
				Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
				Amount:             sdk.NewInt(10),
			},
		},
		Params: types.DefaultParams(),
	}

	wrongIndex := validState
	wrongIndex.WithdrawalList = []types.Withdrawal{validState.WithdrawalList[0]}
	wrongIndex.WithdrawalList[0].ParticipationIndex = 2

	wrongMarket := validState
	wrongMarket.WithdrawalList = []types.Withdrawal{validState.WithdrawalList[0]}
	wrongMarket.WithdrawalList[0].MarketUID = uuid.NewString()

	wrongCreator := validState
	wrongCreator.WithdrawalList = []types.Withdrawal{validState.WithdrawalList[0]}
	wrongCreator.WithdrawalList[0].Address = "new address"

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
			desc:     "wrong participation index",
			genState: &wrongIndex,
			valid:    false,
		},
		{
			desc:     "wrong market",
			genState: &wrongMarket,
			valid:    false,
		},
		{
			desc:     "wrong creator",
			genState: &wrongCreator,
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
