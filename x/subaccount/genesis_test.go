package subaccount_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/subaccount"
	"github.com/sge-network/sge/x/subaccount/types"
)

var (
	subAccOwner = sample.NativeAccAddress()
	micro       = sdkmath.NewInt(1_000_000)
	subAccFunds = sdkmath.NewInt(10_000).Mul(micro)
	subAccAddr  = types.NewAddressFromSubaccount(1)
)

func TestGenesis(t *testing.T) {
	app, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	k := app.SubaccountKeeper

	wantGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		Subaccounts: []types.GenesisSubaccount{
			{
				Address: subAccAddr.String(),
				Owner:   subAccOwner.String(),
				Balance: types.Balance{
					DepositedAmount: subAccFunds,
					SpentAmount:     sdk.ZeroInt(),
					WithdrawmAmount: sdk.ZeroInt(),
					LostAmount:      sdk.ZeroInt(),
				},
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: uint64(time.Now().Add(24 * time.Hour).Unix()),
						Amount:   subAccFunds,
					},
				},
			},
		},
		SubaccountId: 2, // next subaccount id
	}

	subaccount.InitGenesis(ctx, *k, wantGenesis)

	require.Equal(t, wantGenesis.SubaccountId, k.Peek(ctx))
	require.Equal(t, wantGenesis.Params, k.GetParams(ctx))
	require.Len(t, wantGenesis.Subaccounts, 1)
	require.Equal(t, wantGenesis.Subaccounts[0], k.GetAllSubaccounts(ctx)[0])

	exportedGenesis := subaccount.ExportGenesis(ctx, *k)
	require.Equal(t, wantGenesis, *exportedGenesis)
}
