package v3

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the v1.2.0 upgrade.
const UpgradeName = "v1.2.0"

// Expedite governance params
var (
	// DefaultMinExpeditedDepositTokens is the default minimum deposit required for expedited proposals.
	DefaultMinExpeditedDepositTokens = sdkmath.NewInt(50000000000)
	// DefaultExpeditedQuorum is the default quorum percentage required for expedited proposals.
	DefaultExpeditedQuorum = sdk.NewDecWithPrec(750, 3)
	// DefaultExpeditedThreshold is the default voting threshold percentage required for expedited proposals.
	DefaultExpeditedThreshold = sdk.NewDecWithPrec(750, 3)
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
}
