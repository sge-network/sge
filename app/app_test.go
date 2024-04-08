package app_test

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/sge-network/sge/app"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	panicFunc := func() {
		db := tmdb.NewMemDB()
		encCdc := app.MakeEncodingConfig()
		app.NewSgeApp(
			log.NewNopLogger(),
			db,
			nil,
			true,
			map[int64]bool{},
			"",
			0,
			encCdc,
			simtestutil.EmptyAppOptions{},
		)
	}
	require.NotPanics(t, panicFunc)
}
