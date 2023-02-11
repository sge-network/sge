package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestMinterQuery(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	k.SetParams(ctx, params)

	inflation, err := k.Inflation(wctx, &types.QueryInflationRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryInflationResponse{Inflation: k.GetMinter(ctx).Inflation}, inflation)

	phaseStep, err := k.PhaseStep(wctx, &types.QueryPhaseStepRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryPhaseStepResponse{PhaseStep: k.GetMinter(ctx).PhaseStep}, phaseStep)

	phaseProvisions, err := k.PhaseProvisions(wctx, &types.QueryPhaseProvisionsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryPhaseProvisionsResponse{PhaseProvisions: k.GetMinter(ctx).PhaseProvisions}, phaseProvisions)

	endPhaseStatus, err := k.EndPhaseStatus(wctx, &types.QueryEndPhaseStatusRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryEndPhaseStatusResponse{IsInEndPhase: false}, endPhaseStatus)
}
