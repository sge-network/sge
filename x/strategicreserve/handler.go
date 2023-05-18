package strategicreserve

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/sge-network/sge/x/strategicreserve/keeper"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

func NewDataFeeCollectorFeedProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.DataFeeCollectorFeedProposal:
			return keeper.HandleDataFeeCollectorFeedProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized sr proposal content type: %T", c)
		}
	}
}
