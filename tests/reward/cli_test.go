package client

import (
	"testing"

	"github.com/sge-network/sge/testutil/network"
	"github.com/stretchr/testify/suite"
)

func TestE2ETestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	suite.Run(t, NewE2ETestSuite(cfg))
}
