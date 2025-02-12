package keeper

import (
	"context"

	"github.com/sge-network/sge/x/mint/types"
)

var _ types.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// Params returns the params of the module
func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// Inflation returns current inflation.
func (q queryServer) Inflation(ctx context.Context, _ *types.QueryInflationRequest) (*types.QueryInflationResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryInflationResponse{Inflation: minter.Inflation}, nil
}

// PhaseStep returns phase step.
func (q queryServer) PhaseStep(ctx context.Context, _ *types.QueryPhaseStepRequest) (*types.QueryPhaseStepResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryPhaseStepResponse{PhaseStep: minter.PhaseStep}, nil
}

// PhaseProvisions returns current phase provision.
func (q queryServer) PhaseProvisions(ctx context.Context, _ *types.QueryPhaseProvisionsRequest) (*types.QueryPhaseProvisionsResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryPhaseProvisionsResponse{PhaseProvisions: minter.PhaseProvisions}, nil
}

// EndPhaseStatus returns end phase status.
func (q queryServer) EndPhaseStatus(ctx context.Context, _ *types.QueryEndPhaseStatusRequest) (*types.QueryEndPhaseStatusResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}

	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryEndPhaseStatusResponse{
		IsInEndPhase: params.IsEndPhaseByStep(int(minter.PhaseStep)),
	}, nil
}
