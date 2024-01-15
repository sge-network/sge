package simapp_test

import (
	"testing"

	"github.com/sge-network/sge/testutil/simapp"
	"github.com/stretchr/testify/require"
)

func TestGetTestObjects(t *testing.T) {
	_, _, err := simapp.GetTestObjects()
	require.NoError(t, err)
}
