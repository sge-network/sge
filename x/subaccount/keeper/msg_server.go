package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

type msgServer struct {
	keeper Keeper

	accountKeeper keeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, accountKeeper keeper.AccountKeeper, bankKeeper bankkeeper.Keeper) types.MsgServer {
	return &msgServer{
		keeper:        keeper,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

var _ types.MsgServer = msgServer{}
