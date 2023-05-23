package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// SetFeeGrant sets a specific feegrant in the store.
func (k Keeper) SetFeeGrant(ctx sdk.Context, feeGrant types.FeeGrant) {
	store := k.getFeeGrantStore(ctx)
	b := k.cdc.MustMarshal(&feeGrant)
	store.Set(utils.StrBytes(feeGrant.Grantee), b)
}

// GetFeeGrant returns a specific fee grant by grantee address.
func (k Keeper) GetFeeGrant(ctx sdk.Context, grantee string) (val types.FeeGrant, found bool) {
	store := k.getFeeGrantStore(ctx)
	b := store.Get(utils.StrBytes(grantee))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// FeeGrantExists checks if a specific feegrant for the grantee exists or not.
func (k Keeper) FeeGrantExists(ctx sdk.Context, grantee string) bool {
	_, found := k.GetFeeGrant(ctx, grantee)
	return found
}

// RemoveFeeGrant removes a specific feegrant in the store.
func (k Keeper) RemoveFeeGrant(ctx sdk.Context, feeGrant types.FeeGrant) {
	store := k.getFeeGrantStore(ctx)
	store.Delete(utils.StrBytes(feeGrant.Grantee))
}

// GetFeeGrants returns all fee grants.
func (k Keeper) GetFeeGrants(ctx sdk.Context) (list []types.FeeGrant, err error) {
	store := k.getFeeGrantStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.FeeGrant
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetSrPoolSdkFeeGrant(ctx sdk.Context, creator string, grantee sdk.AccAddress) error {
	granter := k.accountKeeper.GetModuleAddress(types.DataFeeCollector)

	// Checking for duplicate entry
	if f, _ := k.feeGrantKeeper.GetAllowance(ctx, granter, grantee); f != nil {
		return types.ErrSDKFeeGrantExists
	}

	err := k.feeGrantKeeper.GrantAllowance(ctx, granter, grantee, types.DefaultFeeGrantAllowance(ctx))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInFeeGrantAllowance, "%s", err)
	}

	return nil
}

func (k Keeper) InvokeSdkFeeGrantIfNotExists(ctx sdk.Context) {
	allFeeGrants, err := k.GetFeeGrants(ctx)
	if err != nil {
		panic(err)
	}

	granter := k.accountKeeper.GetModuleAddress(types.DataFeeCollector)

	for _, fg := range allFeeGrants {
		grantee := sdk.MustAccAddressFromBech32(fg.Grantee)
		_, err := k.feeGrantKeeper.GetAllowance(ctx, granter, grantee)
		if err != nil {
			err = k.feeGrantKeeper.GrantAllowance(ctx, granter, grantee, types.DefaultFeeGrantAllowance(ctx))
			if err != nil {
				panic(err)
			}
		}
	}
}
