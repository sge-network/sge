package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/house/keeper"
	"github.com/sge-network/sge/x/house/types"
	"github.com/stretchr/testify/require"
)

func createNWithdrawals(
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.Withdrawal {
	items := make([]types.Withdrawal, n)

	for i := range items {
		items[i].ID = uint64(i)
		items[i].Creator = sample.AccAddress()
		items[i].Address = testDepositorAddress
		items[i].MarketUID = testMarketUID
		items[i].ParticipationIndex = uint64(i + 1)
		items[i].Mode = types.WithdrawalMode_WITHDRAWAL_MODE_FULL
		items[i].Amount = sdk.NewInt(100)

		keeper.SetWithdrawal(ctx, items[i])
	}
	return items
}

func TestWithdrawalGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNWithdrawals(k, ctx, 10)

	deposits, err := k.GetAllWithdrawals(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(deposits),
	)
}
