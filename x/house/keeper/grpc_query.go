package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sge-network/sge/x/house/types"
)

var _ types.QueryServer = Keeper{}

// Deposits queries all deposits
func (k Keeper) Deposits(c context.Context, req *types.QueryDepositsRequest) (*types.QueryDepositsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var deposits []types.Deposit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	depositStore := prefix.NewStore(store, types.DepositKeyPrefix)

	pageRes, err := query.FilteredPaginate(depositStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		deposit, err := types.UnmarshalDeposit(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			deposits = append(deposits, deposit)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDepositsResponse{Deposits: deposits, Pagination: pageRes}, nil
}

// DepositorDeposits queries all deposits of a give depositor address
func (k Keeper) DepositorDeposits(c context.Context, req *types.QueryDepositorDepositsRequest) (*types.QueryDepositorDepositsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DepositorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "depositor address cannot be empty")
	}

	var deposits []types.Deposit
	ctx := sdk.UnwrapSDKContext(c)

	depAddr, err := sdk.AccAddressFromBech32(req.DepositorAddress)
	if err != nil {
		return nil, err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DepositKeyPrefix)
	delStore := prefix.NewStore(store, types.GetDepositsKey(depAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		deposit, err := types.UnmarshalDeposit(k.cdc, value)
		if err != nil {
			return err
		}
		deposits = append(deposits, deposit)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDepositorDepositsResponse{Deposits: deposits, Pagination: pageRes}, nil
}

// DepositorWithdrawals queries all withdrawals of a give depositor address
func (k Keeper) DepositorWithdrawals(c context.Context, req *types.QueryDepositorWithdrawalsRequest) (*types.QueryDepositorWithdrawalsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DepositorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "depositor address cannot be empty")
	}

	var withdrawals []types.Withdrawal
	ctx := sdk.UnwrapSDKContext(c)

	depAddr, err := sdk.AccAddressFromBech32(req.DepositorAddress)
	if err != nil {
		return nil, err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.WithdrawalKeyPrefix)
	witStore := prefix.NewStore(store, types.GetWithdrawalsKey(depAddr))
	pageRes, err := query.Paginate(witStore, req.Pagination, func(key []byte, value []byte) error {
		withdrawal, err := types.UnmarshalWithdrawal(k.cdc, value)
		if err != nil {
			return err
		}
		withdrawals = append(withdrawals, withdrawal)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDepositorWithdrawalsResponse{Withdrawals: withdrawals, Pagination: pageRes}, nil
}
