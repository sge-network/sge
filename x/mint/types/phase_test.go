package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestPhaseYML(t *testing.T) {
	phase := types.Phase{
		Inflation:       sdk.NewDec(10),
		YearCoefficient: sdk.NewDec(1),
	}

	ymlStr := phase.String()
	require.Equal(t, "inflation: \"10.000000000000000000\"\nyear_coefficient: \"1.000000000000000000\"\n", ymlStr)
}
