package client_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	client "github.com/sge-network/sge/tests/e2e/market"
	"github.com/sge-network/sge/testutil/network"
)

func TestE2ETestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	suite.Run(t, client.NewE2ETestSuite(cfg))
}
