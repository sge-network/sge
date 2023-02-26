package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSportEventQuerySingle(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSportEvent(k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QuerySportEventRequest
		response *types.QuerySportEventResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QuerySportEventRequest{
				Uid: msgs[0].UID,
			},
			response: &types.QuerySportEventResponse{SportEvent: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QuerySportEventRequest{
				Uid: msgs[1].UID,
			},
			response: &types.QuerySportEventResponse{SportEvent: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QuerySportEventRequest{
				Uid: cast.ToString(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.SportEvent(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
