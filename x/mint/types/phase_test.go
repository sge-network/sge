package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/mint/types"
)

func TestPhaseYML(t *testing.T) {
	phase := types.Phase{
		Inflation:       sdkmath.LegacyNewDec(10),
		YearCoefficient: sdkmath.LegacyNewDec(1),
	}

	ymlStr := phase.String()
	require.Equal(
		t,
		"inflation: \"10.000000000000000000\"\nyear_coefficient: \"1.000000000000000000\"\n",
		ymlStr,
	)
}
