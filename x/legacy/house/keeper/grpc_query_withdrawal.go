package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/legacy/house/types"
)

// WithdrawalsByAccount returns all withdrawals of a given account address
func (k Keeper) WithdrawalsByAccount(c context.Context,
	req *types.QueryWithdrawalsByAccountRequest,
) (*types.QueryWithdrawalsByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var withdrawals []types.Withdrawal
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(k.getWithdrawalStore(ctx), types.GetWithdrawalListPrefix(req.Address))
	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var withdrawal types.Withdrawal
		if err := k.cdc.Unmarshal(value, &withdrawal); err != nil {
			return err
		}

		withdrawals = append(withdrawals, withdrawal)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWithdrawalsByAccountResponse{Withdrawals: withdrawals, Pagination: pageRes}, nil
}

// Withdrawal returns specific withdrawal.
func (k Keeper) Withdrawal(c context.Context,
	req *types.QueryWithdrawalRequest,
) (*types.QueryWithdrawalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetWithdraw(
		ctx,
		req.DepositorAddress,
		req.MarketUid,
		req.ParticipationIndex,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryWithdrawalResponse{Withdrawal: val}, nil
}
