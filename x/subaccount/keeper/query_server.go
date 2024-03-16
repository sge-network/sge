package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/subaccount/types"
)

func NewQueryServer(keeper Keeper) types.QueryServer {
	return queryServer{keeper: keeper}
}

type queryServer struct {
	keeper Keeper
}

func (q queryServer) Subaccount(goCtx context.Context, request *types.QuerySubaccountRequest) (*types.QuerySubaccountResponse, error) {
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	subaccountAddr, exists := q.keeper.GetSubAccountByOwner(ctx, addr)
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	accSummary, exists := q.keeper.GetAccountSummary(ctx, subaccountAddr)
	if !exists {
		panic("subaccount exists but balance not found")
	}

	lockedBalances, _ := q.keeper.GetBalances(ctx, subaccountAddr, types.BalanceType_BALANCE_TYPE_LOCKED)
	unlockedBalances, unlockedAmount := q.keeper.GetBalances(ctx, subaccountAddr, types.BalanceType_BALANCE_TYPE_UNLOCKED)

	bankBalance := q.keeper.bankKeeper.GetBalance(ctx, subaccountAddr, params.DefaultBondDenom)

	return &types.QuerySubaccountResponse{
		Address:                     subaccountAddr.String(),
		Balance:                     accSummary,
		LockedBalance:               lockedBalances,
		UnlockedBalance:             unlockedBalances,
		WithdrawableUnlockedBalance: accSummary.WithdrawableUnlockedBalance(unlockedAmount, bankBalance.Amount),
	}, nil
}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.keeper.GetParams(sdk.UnwrapSDKContext(ctx))
	return &types.QueryParamsResponse{Params: params}, nil
}
