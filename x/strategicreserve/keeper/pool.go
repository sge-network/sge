package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sge-network/sge/x/strategicreserve/types"
)

// GetBetReserveAcc returns the `bet_reserve` module account
func (k Keeper) GetBetReserveAcc(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.BetReserveName)
}

// GetWinningsCollectorAcc returns the `winnings_collector` module account
func (k Keeper) GetWinningsCollectorAcc(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.WinningsCollectorName)
}

// GetSRPoolAcc returns the `sr_pool` module account
func (k Keeper) GetSRPoolAcc(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.SRPoolName)
}

// GetSRPoolBalance returns the balance of `sr_pool` module account
func (k Keeper) GetSRPoolBalance(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, k.GetSRPoolAcc(ctx).GetAddress(), k.DepositDenom(ctx))
}
