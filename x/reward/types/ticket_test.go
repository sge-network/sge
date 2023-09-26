package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/reward/types"
)

func TestCreateCampaignPayloadValidation(t *testing.T) {
	funderAddr := sample.AccAddress()
	poolBalance := sdkmath.NewInt(1000)
	blockTime := uint64(time.Now().Unix())

	tests := []struct {
		name    string
		payload types.CreateCampaignPayload
		err     error
	}{
		{
			name: "start after end",
			payload: types.CreateCampaignPayload{
				FunderAddress: funderAddr,
				StartTs:       blockTime + 10,
				EndTs:         blockTime + 5,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "expired",
			payload: types.CreateCampaignPayload{
				FunderAddress: funderAddr,
				StartTs:       blockTime - 1,
				EndTs:         blockTime,
			},
			err: types.ErrExpiredCampaign,
		},
		{
			name: "not enough defs",
			payload: types.CreateCampaignPayload{
				FunderAddress: funderAddr,
				StartTs:       blockTime + 1,
				EndTs:         blockTime + 2,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "nil pool amount",
			payload: types.CreateCampaignPayload{
				FunderAddress: funderAddr,
				StartTs:       blockTime + 1,
				EndTs:         blockTime + 2,
				RewardDefs: []types.Definition{
					{
						RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
						DstAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
						Amount:     sdkmath.NewInt(1000),
						UnlockTS:   0,
					},
				},
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid pool amount",
			payload: types.CreateCampaignPayload{
				FunderAddress: funderAddr,
				StartTs:       blockTime + 1,
				EndTs:         blockTime + 2,
				RewardDefs: []types.Definition{
					{
						RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
						DstAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
						Amount:     sdkmath.NewInt(1000),
						UnlockTS:   0,
					},
				},
				PoolAmount: sdkmath.NewInt(0),
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "valid",
			payload: types.CreateCampaignPayload{
				FunderAddress: funderAddr,
				StartTs:       blockTime + 1,
				EndTs:         blockTime + 2,
				Type:          types.RewardType_REWARD_TYPE_SIGNUP,
				RewardDefs: []types.Definition{
					{
						RecType:    types.ReceiverType_RECEIVER_TYPE_SINGLE,
						DstAccType: types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN,
						Amount:     sdkmath.NewInt(1000),
						UnlockTS:   0,
					},
				},
				PoolAmount: poolBalance,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(blockTime)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateCampaignPayloadValidation(t *testing.T) {
	blockTime := uint64(time.Now().Unix())

	tests := []struct {
		name    string
		payload types.UpdateCampaignPayload
		err     error
	}{
		{
			name: "expired",
			payload: types.UpdateCampaignPayload{
				EndTs: blockTime - 1,
			},
			err: types.ErrExpiredCampaign,
		},
		{
			name: "valid now",
			payload: types.UpdateCampaignPayload{
				EndTs: blockTime,
			},
		},
		{
			name: "valid",
			payload: types.UpdateCampaignPayload{
				EndTs: blockTime + 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(blockTime)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
