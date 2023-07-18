package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

func TestMsgPlaceValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgPlace
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgPlace{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid bet message",
			msg: types.MsgPlace{
				Creator: sample.AccAddress(),
				Bet: &types.PlaceBetFields{
					UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
					Amount: sdk.NewInt(int64(10)),
					Ticket: "Ticket",
				},
			},
		},
		{
			name: "invalid bet UID",
			msg: types.MsgPlace{
				Creator: sample.AccAddress(),
				Bet: &types.PlaceBetFields{
					UID: "Invalid UID",
				},
			},
			err: types.ErrInvalidBetUID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewBet(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		inputBet := &types.PlaceBetFields{
			UID:    "betUid",
			Ticket: "ticket",
			Amount: sdk.NewInt(int64(10)),
		}
		creator := "creator"
		inputBetOdds := &types.BetOdds{
			UID:       "Oddsuid",
			MarketUID: "marketUID",
			Value:     "1000",
		}

		expectedBet := &types.Bet{
			UID:       inputBet.UID,
			Creator:   creator,
			MarketUID: inputBetOdds.MarketUID,
			OddsUID:   inputBetOdds.UID,
			OddsValue: inputBetOdds.Value,
			Amount:    inputBet.Amount,
			OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
		}
		res := types.NewBet(creator, inputBet, types.OddsType_ODDS_TYPE_DECIMAL, inputBetOdds)
		require.Equal(t, expectedBet, res)
	})
}
