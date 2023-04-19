package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/sge-network/sge/x/mint/keeper"
	sgeMintTypes "github.com/sge-network/sge/x/mint/types"
	"github.com/spf13/cast"
)

var gaugeKeys = []string{"minted_tokens"}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter
	minter := k.GetMinter(ctx)

	// fetch stored params
	params := k.GetParams(ctx)
	currentBlock := ctx.BlockHeight()

	// detect current phase according to current block
	currentPhase, currentPhaseStep := minter.CurrentPhase(params, currentBlock)

	// set the new minter properties if the phase has changed or inflation has changed
	if currentPhaseStep != cast.ToInt(minter.PhaseStep) || !minter.Inflation.Equal(currentPhase.Inflation) {
		// set new inflation rate
		newInflation := currentPhase.Inflation
		minter.Inflation = newInflation

		// set new phase step
		minter.PhaseStep = cast.ToInt32(currentPhaseStep)

		// set phase provisions of new phase step
		totalSupply := k.TokenSupply(ctx, params.MintDenom)
		minter.PhaseProvisions = minter.NextPhaseProvisions(totalSupply, params.ExcludeAmount, currentPhase)

		// store minter
		k.SetMinter(ctx, minter)
	}

	// if the inflation rate is zero, means that we have no minting, so the rest of the code should not be called
	if minter.Inflation.Equal(sdk.ZeroDec()) {
		return
	}

	// mint coins, update supply
	mintedCoin, truncatedTokens := minter.BlockProvisions(params, currentPhaseStep)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// set truncated value in this block to be added to provision calculation in the next block
	minter.TruncatedTokens = truncatedTokens
	k.SetMinter(ctx, minter)

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), gaugeKeys...)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(sgeMintTypes.AttributeKeyPhaseProvisions, minter.PhaseProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
