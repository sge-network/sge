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

func (m msgServer) PlaceBet(ctx context.Context, msg *types.MsgPlaceBet) (*types.MsgPlaceBetResponse, error) {
	handler := m.keeper.msgRouter.Handler(msg.Msg)
	if handler == nil {
		return nil, status.Error(codes.Unavailable, "the message type is not supported by the subaccount router")
	}
	sdkResult, err := handler(sdk.UnwrapSDKContext(ctx), msg.Msg)
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
