package subaccount

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/sge-network/sge/x/subaccount/keeper"
)

// NewHandler initialize a new sdk.handler instance for registered messages
func NewHandler(k keeper.Keeper, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper) sdk.Handler {
	keeper.NewMsgServerImpl(k, accountKeeper, bankKeeper)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// ctx = ctx.WithEventManager(sdk.NewEventManager())

		// switch msg := msg.(type) {
		// default:
		//  	errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		//  	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		// }

		return nil, nil
	}
}
