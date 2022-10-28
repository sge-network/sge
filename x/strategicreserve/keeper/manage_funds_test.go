package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/stretchr/testify/require"
)

func TestProcessBetPlacement(t *testing.T) {
	k, ctx := setupKeeper(t)
	user := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc          string
		bettorAddress sdk.AccAddress
		betFee        sdk.Coin
		betAmount     sdk.Int
		extraPayout   sdk.Int
		uniqueLock    string
		err           error
	}{
		{
			desc:          "Success!",
			bettorAddress: user.Address,
			betFee:        sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(1)),
			betAmount:     sdk.NewIntFromUint64(99),
			extraPayout:   sdk.NewIntFromUint64(198),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7f",
			err:           nil,
		},
		{
			desc:          "Failure! Lock already exists in payout store",
			bettorAddress: user.Address,
			betFee:        sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(1)),
			betAmount:     sdk.NewIntFromUint64(99),
			extraPayout:   sdk.NewIntFromUint64(198),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7x",
			err:           types.ErrLockAlreadyExists,
		},
		{
			desc:          "Failure! Insufficient user balance",
			bettorAddress: user.Address,
			betFee:        sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(1)),
			betAmount:     sdk.NewIntFromUint64(45000000000000),
			extraPayout:   sdk.NewIntFromUint64(9000000),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7l",
			err:           types.ErrInsufficientUserBalance,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			reserver := types.Reserver{
				SrPool: &types.SRPool{
					LockedAmount:   sdk.NewIntFromUint64(4500),
					UnlockedAmount: sdk.NewIntFromUint64(149999999995500),
				},
			}
			k.SetReserver(ctx, reserver)

			k.SetPayoutLock(ctx, "32932b20-8737-490b-b00b-8c16eccd8e7x")

			err := k.ProcessBetPlacement(ctx, tc.bettorAddress, tc.betFee, tc.betAmount,
				tc.extraPayout, tc.uniqueLock)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestBettorWins(t *testing.T) {
	k, ctx := setupKeeper(t)
	user := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc          string
		bettorAddress sdk.AccAddress
		betAmount     sdk.Int
		extraPayout   sdk.Int
		uniqueLock    string
		err           error
	}{
		{
			desc:          "Success! Payout done",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(45),
			extraPayout:   sdk.NewIntFromUint64(90),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7p",
			err:           nil,
		},
		{
			desc:          "Failure! Payout lock does not exist in the payout store",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(45),
			extraPayout:   sdk.NewIntFromUint64(90),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7x",
			err:           types.ErrPayoutLockDoesnotExist,
		},
		{
			desc:          "Failure! SR locked amount has insufficient balance",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(50),
			extraPayout:   sdk.NewIntFromUint64(5000),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7p",
			err:           types.ErrInsufficientLockedAmountInSrPool,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			k.SetPayoutLock(ctx, "32932b20-8737-490b-b00b-8c16eccd8e7p")

			reserver := types.Reserver{
				SrPool: &types.SRPool{
					LockedAmount:   sdk.NewIntFromUint64(4500),
					UnlockedAmount: sdk.NewIntFromUint64(150000000000000),
				},
			}
			k.SetReserver(ctx, reserver)

			err := k.BettorWins(ctx, tc.bettorAddress, tc.betAmount,
				tc.extraPayout, tc.uniqueLock)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestBettorLoses(t *testing.T) {
	k, ctx := setupKeeper(t)
	user := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc        string
		address     sdk.AccAddress
		betAmount   sdk.Int
		extraPayout sdk.Int
		uniqueLock  string
		err         error
	}{
		{
			desc:        "Success! Payout done",
			address:     user.Address,
			betAmount:   sdk.NewIntFromUint64(45),
			extraPayout: sdk.NewIntFromUint64(90),
			uniqueLock:  "32932b20-8737-490b-b00b-8c16eccd8e7f",
			err:         nil,
		},
		{
			desc:        "Failure! Payout lock does not exist",
			address:     user.Address,
			betAmount:   sdk.NewIntFromUint64(45),
			extraPayout: sdk.NewIntFromUint64(90),
			uniqueLock:  "32932b20-8737-490b-b00b-8c16eccd8e7x",
			err:         sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, "32932b20-8737-490b-b00b-8c16eccd8e7x"),
		},
		{
			desc:        "Failure! Insufficient balance in Bet Reserve Account",
			address:     user.Address,
			betAmount:   sdk.NewIntFromUint64(5000),
			extraPayout: sdk.NewIntFromUint64(5000),
			uniqueLock:  "32932b20-8737-490b-b00b-8c16eccd8e7f",
			err:         types.ErrInsufficientBalanceInModuleAccount,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			k.SetPayoutLock(ctx, "32932b20-8737-490b-b00b-8c16eccd8e7f")

			reserver := types.Reserver{
				SrPool: &types.SRPool{
					LockedAmount:   sdk.NewIntFromUint64(6000),
					UnlockedAmount: sdk.NewIntFromUint64(150000000000000),
				},
			}
			k.SetReserver(ctx, reserver)

			err := k.BettorLoses(ctx, tc.address, tc.betAmount,
				tc.extraPayout, tc.uniqueLock)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRefundBettor(t *testing.T) {
	k, ctx := setupKeeper(t)
	user := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc          string
		bettorAddress sdk.AccAddress
		betAmount     sdk.Int
		extraPayout   sdk.Int
		uniqueLock    string
		err           error
	}{
		{
			desc:          "Success! Bettor is refunded",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(45),
			extraPayout:   sdk.NewIntFromUint64(90),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7f",
			err:           nil,
		},
		{
			desc:          "Failure! Payout lock does not exist",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(45),
			extraPayout:   sdk.NewIntFromUint64(90),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7x",
			err:           sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, "32932b20-8737-490b-b00b-8c16eccd8e7x"),
		},
		{
			desc:          "Failure! Insufficient balance in Bet Reserve Account",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(5000),
			extraPayout:   sdk.NewIntFromUint64(5000),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7f",
			err:           sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount, types.BetReserveName),
		},
		{
			desc:          "Failure! Insufficient balance in SR Locked Amount",
			bettorAddress: user.Address,
			betAmount:     sdk.NewIntFromUint64(5000),
			extraPayout:   sdk.NewIntFromUint64(6001),
			uniqueLock:    "32932b20-8737-490b-b00b-8c16eccd8e7f",
			err:           types.ErrInsufficientLockedAmountInSrPool,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			k.SetPayoutLock(ctx, "32932b20-8737-490b-b00b-8c16eccd8e7f")

			reserver := types.Reserver{
				SrPool: &types.SRPool{
					LockedAmount:   sdk.NewIntFromUint64(6000),
					UnlockedAmount: sdk.NewIntFromUint64(150000000000000),
				},
			}
			k.SetReserver(ctx, reserver)

			err := k.RefundBettor(ctx, tc.bettorAddress, tc.betAmount,
				tc.extraPayout, tc.uniqueLock)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUpdateSrPool(t *testing.T) {
	k, ctx := setupKeeper(t)

	tc := []struct {
		desc              string
		newLockedAmount   sdk.Int
		newUnlockedAmount sdk.Int
	}{
		{
			desc:              "Success! SR_Pool updated",
			newLockedAmount:   sdk.NewInt(50000000000000),
			newUnlockedAmount: sdk.NewInt(100000000000000),
		},
	}

	t.Run(tc[0].desc, func(t *testing.T) {
		reserver := types.Reserver{
			SrPool: &types.SRPool{
				LockedAmount:   sdk.NewIntFromUint64(0),
				UnlockedAmount: sdk.NewIntFromUint64(150000000000000),
			},
		}
		k.SetReserver(ctx, reserver)

		k.UpdateSrPool(ctx, tc[0].newLockedAmount, tc[0].newUnlockedAmount)

		updatedReserver := k.GetReserver(ctx)
		require.Equal(t, tc[0].newLockedAmount, updatedReserver.SrPool.LockedAmount)
		require.Equal(t, tc[0].newUnlockedAmount, updatedReserver.SrPool.UnlockedAmount)
	})
}

func TestTransferFundsFromUserToModule(t *testing.T) {
	k, ctx := setupKeeper(t)

	user := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc          string
		address       sdk.AccAddress
		moduleAccName string
		amount        sdk.Int
		err           error
	}{
		{
			desc:          "Success! Funds transferred from user to module account",
			address:       user.Address,
			amount:        sdk.NewInt(450),
			moduleAccName: types.SRPoolName,
			err:           nil,
		},
		{
			desc:          "Failure! Insufficient user balance",
			address:       user.Address,
			amount:        sdk.NewInt(45000000000000000),
			moduleAccName: types.SRPoolName,
			err:           sdkerrors.Wrapf(types.ErrInsufficientUserBalance, user.Address.String()),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := k.TransferFundsFromUserToModule(ctx,
				tc.address, tc.moduleAccName, tc.amount)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTransferFundsFromModuleToUser(t *testing.T) {
	k, ctx := setupKeeper(t)

	user := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc          string
		moduleAccName string
		address       sdk.AccAddress
		amount        sdk.Int
		err           error
	}{
		{
			desc:          "Success! Funds transferred from module account to user",
			moduleAccName: types.SRPoolName,
			address:       user.Address,
			amount:        sdk.NewInt(450),
			err:           nil,
		},
		{
			desc:          "Failure! Insufficient balance in module account",
			moduleAccName: types.BetReserveName,
			address:       user.Address,
			amount:        sdk.NewInt(4500),
			err:           sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount, types.BetReserveName),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			reserver := types.Reserver{
				SrPool: &types.SRPool{
					LockedAmount:   sdk.ZeroInt(),
					UnlockedAmount: sdk.NewIntFromUint64(150000000000000)},
			}
			k.SetReserver(ctx, reserver)

			err := k.TransferFundsFromModuleToUser(ctx, tc.moduleAccName,
				tc.address, tc.amount)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTransferFundsFromModuleToModule(t *testing.T) {
	k, ctx := setupKeeper(t)

	tcs := []struct {
		desc            string
		senderModule    string
		recipientModule string
		amount          sdk.Int
		err             error
	}{
		{
			desc:            "Success! Funds transferred from module account to module",
			senderModule:    types.BetReserveName,
			recipientModule: types.SRPoolName,
			amount:          sdk.NewInt(450),
			err:             nil,
		},
		{
			desc:            "Failure! Sender & recipient module names are same",
			senderModule:    types.BetReserveName,
			recipientModule: types.BetReserveName,
			amount:          sdk.NewInt(450),
			err:             types.ErrDuplicateSenderAndRecipientModule,
		},
		{
			desc:            "Failure! Insufficient balance in module account",
			senderModule:    types.BetReserveName,
			recipientModule: types.SRPoolName,
			amount:          sdk.NewInt(4500),
			err:             sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount, types.BetReserveName),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			reserver := types.Reserver{
				SrPool: &types.SRPool{
					LockedAmount:   sdk.ZeroInt(),
					UnlockedAmount: sdk.NewIntFromUint64(150000000000000)},
			}
			k.SetReserver(ctx, reserver)

			err := k.TransferFundsFromModuleToModule(ctx, tc.senderModule,
				tc.recipientModule, tc.amount)
			if tc.err != nil {
				require.True(t, errors.Is(tc.err, err))
				return
			}
			require.NoError(t, err)
		})
	}
}
