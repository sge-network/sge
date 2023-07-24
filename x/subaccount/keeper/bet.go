package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/subaccount/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m msgServer) PlaceBet(goCtx context.Context, msg *types.MsgPlaceBet) (*types.MsgPlaceBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// address: cosmos1234
	// creator: cosmos-subaccount-1234

	subAccAddr, exists := m.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	// swap the creator address with the subaccount address
	msg.Msg.Creator = subAccAddr.String()

	handler := m.keeper.msgRouter.Handler(msg.Msg)
	if handler == nil {
		return nil, status.Error(codes.Unavailable, "the message type is not supported by the subaccount router")
	}

	sdkResult, err := handler(ctx, msg.Msg)
	if err != nil {
		return nil, err
	}

	resp := &bettypes.MsgPlaceBetResponse{}
	err = proto.Unmarshal(sdkResult.Data, resp)
	if err != nil {
		return nil, err
	}

	return &types.MsgPlaceBetResponse{Response: resp}, nil
}
