package simulation_test

import (
	//#nosec
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/market/simulation"
)

func TestParamChangest(t *testing.T) {
	s := rand.NewSource(1)
	//#nosec
	r := rand.New(s)

	expected := []struct {
		composedKey string
		key         string
		subspace    string
	}{
		{"market/MinBetAmount", "MinBetAmount", "market"},
		{"market/MinBetFee", "MinBetFee", "market"},
		{"market/MaxBetFee", "MaxBetFee", "market"},
	}

	paramChanges := simulation.ParamChanges(r)
	require.Len(t, paramChanges, 4)

	for i, p := range paramChanges {
		require.Equal(t, expected[i].composedKey, p.ComposedKey())
		require.Equal(t, expected[i].key, p.Key())
		require.Equal(t, expected[i].subspace, p.Subspace())
	}
}
