package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/house/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
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
