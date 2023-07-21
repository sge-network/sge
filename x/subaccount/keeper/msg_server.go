package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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

// sumBalanceUnlocks sums all the balances to unlock and returns the total amount. It
// returns an error if any of the unlock times is expired.
func sumBalanceUnlocks(ctx sdk.Context, balanceUnlocks []*types.LockedBalance) (sdk.Int, error) {
	moneyToSend := sdk.NewInt(0)

	for _, balanceUnlock := range balanceUnlocks {
		if balanceUnlock.UnlockTime.Unix() < ctx.BlockTime().Unix() {
			return sdk.Int{}, types.ErrUnlockTokenTimeExpired
		}

		moneyToSend = moneyToSend.Add(balanceUnlock.Amount)
	}

	return moneyToSend, nil
}

// sendCoinsToSubaccount sends the coins to the subaccount.
func (m msgServer) sendCoinsToSubaccount(ctx sdk.Context, senderAccount sdk.AccAddress, subAccountAddress sdk.AccAddress, moneyToSend sdk.Int) error {

	denom := m.keeper.GetParams(ctx).LockedBalanceDenom
	err := m.bankKeeper.SendCoins(ctx, senderAccount, subAccountAddress, sdk.NewCoins(sdk.NewCoin(denom, moneyToSend)))
	if err != nil {
		return errors.Wrap(err, "unable to send coins")
	}

	return nil
}
