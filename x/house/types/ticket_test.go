package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/house/types"
)

func TestDepositTicketPayloadValidation(t *testing.T) {
	depositor := sample.AccAddress()
	tests := []struct {
		name    string
		payload types.DepositTicketPayload
		err     error
	}{
		{
			name: "valid",
			payload: types.DepositTicketPayload{
				DepositorAddress: depositor,
				KycData: sgetypes.KycDataPayload{
					Ignore: true,
				},
			},
		},
		{
			name: "valid kyc",
			payload: types.DepositTicketPayload{
				DepositorAddress: depositor,
				KycData: sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       depositor,
				},
			},
		},
		{
			name: "invalid address",
			payload: types.DepositTicketPayload{
				DepositorAddress: "invalid addr",
				KycData: sgetypes.KycDataPayload{
					Ignore: true,
				},
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "invalid valid kyc",
			payload: types.DepositTicketPayload{
				DepositorAddress: depositor,
				KycData: sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       sample.AccAddress(),
				},
			},
			err: types.ErrUserKycFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(tt.payload.DepositorAddress)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestWithdrawTicketPayloadValidation(t *testing.T) {
	depositor := sample.AccAddress()
	tests := []struct {
		name    string
		payload types.WithdrawTicketPayload
		err     error
	}{
		{
			name: "valid",
			payload: types.WithdrawTicketPayload{
				DepositorAddress: depositor,
				KycData: sgetypes.KycDataPayload{
					Ignore: true,
				},
			},
		},
		{
			name: "valid kyc",
			payload: types.WithdrawTicketPayload{
				DepositorAddress: depositor,
				KycData: sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       depositor,
				},
			},
		},
		{
			name: "invalid address",
			payload: types.WithdrawTicketPayload{
				DepositorAddress: "invalid addr",
				KycData: sgetypes.KycDataPayload{
					Ignore: true,
				},
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "invalid valid kyc",
			payload: types.WithdrawTicketPayload{
				DepositorAddress: depositor,
				KycData: sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       sample.AccAddress(),
				},
			},
			err: types.ErrUserKycFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(tt.payload.DepositorAddress)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
