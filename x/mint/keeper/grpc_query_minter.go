package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/mint/types"
)

// Inflation returns current inflation.
func (k Keeper) Inflation(
	c context.Context,
	_ *types.QueryInflationRequest,
) (*types.QueryInflationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)

	return &types.QueryInflationResponse{Inflation: minter.Inflation}, nil
}

// PhaseStep returns phase step.
func (k Keeper) PhaseStep(
	c context.Context,
	_ *types.QueryPhaseStepRequest,
) (*types.QueryPhaseStepResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)

	return &types.QueryPhaseStepResponse{PhaseStep: minter.PhaseStep}, nil
}

// PhaseProvisions returns current phase provision.
func (k Keeper) PhaseProvisions(
	c context.Context,
	_ *types.QueryPhaseProvisionsRequest,
) (*types.QueryPhaseProvisionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)

	return &types.QueryPhaseProvisionsResponse{PhaseProvisions: minter.PhaseProvisions}, nil
}

// EndPhaseStatus returns end phase status.
func (k Keeper) EndPhaseStatus(
	c context.Context,
	_ *types.QueryEndPhaseStatusRequest,
) (*types.QueryEndPhaseStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	return &types.QueryEndPhaseStatusResponse{
		IsInEndPhase: params.IsEndPhaseByStep(int(minter.PhaseStep)),
	}, nil
}
