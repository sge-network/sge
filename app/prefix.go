package app

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/app/params"
)

var (
	// AccountPubKeyPrefix used for generating public key
	AccountPubKeyPrefix = keepers.AccountAddressPrefix + "pub"
	// ValidatorAddressPrefix used for generating validator address
	ValidatorAddressPrefix = keepers.AccountAddressPrefix + "valoper"
	// ValidatorPubKeyPrefix used for generating validator public key
	ValidatorPubKeyPrefix = keepers.AccountAddressPrefix + "valoperpub"
	// ConsNodeAddressPrefix used for generating consensus node address
	ConsNodeAddressPrefix = keepers.AccountAddressPrefix + "valcons"
	// ConsNodePubKeyPrefix used for generating consensus node public key
	ConsNodePubKeyPrefix = keepers.AccountAddressPrefix + "valconspub"
)

// SetConfig sets prefixes configuration
func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(keepers.AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)

	err := sdk.RegisterDenom(params.HumanCoinUnit, sdkmath.LegacyOneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(params.BaseCoinUnit, sdkmath.LegacyNewDecWithPrec(1, params.SGEExponent))
	if err != nil {
		panic(err)
	}

	config.Seal()
}
