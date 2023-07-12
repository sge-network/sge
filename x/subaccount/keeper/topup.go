package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

func (m msgServer) TopUp(ctx context.Context, msg *types.MsgTopUp) (*types.MsgTopUpResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	_, err := sumBalanceUnlocks(sdkContext, msg.LockedBalances)
	if err != nil {
		return nil, err
	}

	subaccountOwner := sdk.MustAccAddressFromBech32(msg.SubAccount)
	if !m.keeper.HasSubAccount(sdkContext, subaccountOwner) {
		return nil, types.ErrSubaccountDoesNotExist
	}

	//subaccountID := m.keeper.GetSubAccountByOwner(sdkContext, subaccountOwner)
	//_ := m.keeper.GetBalance(sdkContext, subaccountID)

	return &types.MsgTopUpResponse{}, nil
}
