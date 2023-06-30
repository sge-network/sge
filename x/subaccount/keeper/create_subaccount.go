package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) CreateSubAccount(
	ctx context.Context,
	request *types.MsgCreateSubAccountRequest,
) (*types.MsgCreateAccountResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	//_ := sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(0))

	for _, balanceUnlock := range request.LockedBalances {
		if balanceUnlock.UnlockTime.Unix() < sdkContext.BlockTime().Unix() {
			return nil, types.ErrUnlockTokenTimeExpired
		}
	}

	return &types.MsgCreateAccountResponse{}, nil
}
