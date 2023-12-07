package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/simapp"
)

func TestExport(t *testing.T) {
	testCases := []struct {
		name          string
		forZeroHeight bool
	}{
		{
			"for zero height",
			true,
		},
		{
			"for non-zero height",
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp, _, err := simapp.GetTestObjectsWithOptions(
				simapp.Options{
					CreateGenesisValidators: true,
				},
			)
			require.NoError(t, err)

			tApp.Commit()
			_, err = tApp.ExportAppStateAndValidators(tc.forZeroHeight, []string{})
			require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
		})
	}
}
