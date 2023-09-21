package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/bet/types"
)

func TestMsgWagerValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgWager
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgWager{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid bet message",
			msg: types.MsgWager{
				Creator: sample.AccAddress(),
				Props: &types.WagerProps{
					UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
					Amount: sdkmath.NewInt(int64(10)),
					Ticket: "Ticket",
				},
			},
		},
		{
			name: "invalid bet UID",
			msg: types.MsgWager{
				Creator: sample.AccAddress(),
				Props: &types.WagerProps{
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
		inputBet := &types.WagerProps{
			UID:    "betUid",
			Ticket: "ticket",
			Amount: sdkmath.NewInt(int64(10)),
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
