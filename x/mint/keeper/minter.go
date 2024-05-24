package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/mint/types"
)

// GetMinter returns the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic(types.ErrTextNilMinter)
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// SetMinter sets the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdkmath.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// TokenSupply implements an alias call to the underlying bank keeper's
// to be used in BeginBlocker.
func (k Keeper) TokenSupply(ctx sdk.Context, denom string) sdkmath.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}

// BondedRatio implements an alias call to the underlying staking keeper's
// to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) sdkmath.LegacyDec {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
