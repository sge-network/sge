package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// GetDeposit returns a specific deposit.
func (k Keeper) GetDeposit(ctx sdk.Context, depositIDBytes []byte) (deposit types.Deposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DepositKeyPrefix)
	value := store.Get(depositIDBytes)
	if value == nil {
		return deposit, false
	}

	deposit = types.MustUnmarshalDeposit(k.cdc, value)

	return deposit, true
}

// IterateAllDeposits iterates through all of the deposits.
func (k Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.DepositKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDeposit(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDeposits returns all deposits used during genesis dump.
func (k Keeper) GetAllDeposits(ctx sdk.Context) (deposits []types.Deposit) {
	k.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})

	return deposits
}

// SetDeposit sets a deposit.
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) error {
	depAddress := sdk.MustAccAddressFromBech32(deposit.DepositorAddress)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DepositKeyPrefix)
	b := types.MustMarshalDeposit(k.cdc, deposit)

	depoistKey := types.GetDepositKey(depAddress, deposit.SportEventUid, deposit.ParticipantId)
	store.Set(depoistKey, b)
	return nil
}

// Deposit performs a deposit, set/update everything necessary within the store.
func (k Keeper) Deposit(ctx sdk.Context, depAddr sdk.AccAddress, sportEventUid string, depAmt sdk.Int) (error, uint64) {
	// Create the deposit object
	deposit := types.NewDeposit(depAddr, sportEventUid, depAmt, sdk.ZeroInt(), 0)

	// Set the house participation fee
	deposit.SetHouseParticipationFee(k.HouseParticipationFee(ctx))

	participantId, err := k.orderBookKeeper.AddBookParticipant(
		ctx, depAddr, sportEventUid, deposit.Liquidity, deposit.Fee, types.HouseParticipationFeeName,
	)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrOrderBookDepositProcessing, "%s", err), participantId
	}

	deposit.ParticipantId = participantId

	if err = k.SetDeposit(ctx, deposit); err != nil {
		return err, participantId
	}

	return nil, participantId
}
