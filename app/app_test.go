package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"

	"github.com/sge-network/sge/app"
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
			sdksimapp.EmptyAppOptions{},
		)
	}
	require.NotPanics(t, panicFunc)
}
