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
	blockTime := uint64(time.Now().Unix())

	tests := []struct {
		name    string
		payload types.CreateCampaignPayload
		err     error
	}{
		{
			name: "start after end",
			payload: types.CreateCampaignPayload{
				Promoter: funderAddr,
				StartTs:  blockTime + 10,
				EndTs:    blockTime + 5,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "expired",
			payload: types.CreateCampaignPayload{
				Promoter: funderAddr,
				StartTs:  blockTime - 1,
				EndTs:    blockTime,
			},
			err: types.ErrExpiredCampaign,
		},
		{
			name: "not enough defs",
			payload: types.CreateCampaignPayload{
				Promoter: funderAddr,
				StartTs:  blockTime + 1,
				EndTs:    blockTime + 2,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "valid",
			payload: types.CreateCampaignPayload{
				Promoter:         funderAddr,
				StartTs:          blockTime + 1,
				EndTs:            blockTime + 2,
				Category:         types.RewardCategory_REWARD_CATEGORY_SIGNUP,
				RewardType:       types.RewardType_REWARD_TYPE_AFFILIATE_SIGNUP,
				RewardAmountType: types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED,
				RewardAmount: &types.RewardAmount{
					SubaccountAmount: sdkmath.NewInt(1000),
					UnlockPeriod:     uint64(time.Now().Add(10 * time.Minute).Unix()),
				},
				ClaimsPerCategory: 1,
				IsActive:          true,
				Meta:              "sample campaign",
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
