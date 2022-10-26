package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AccountAddressPrefix prefix used for generating account address
	AccountAddressPrefix = "sge"
)

var (
	// AccountPubKeyPrefix used for generating public key
	AccountPubKeyPrefix = AccountAddressPrefix + "pub"
	// ValidatorAddressPrefix used for generating validator address
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	// ValidatorPubKeyPrefix used for generating validator public key
	ValidatorPubKeyPrefix = AccountAddressPrefix + "valoperpub"
	// ConsNodeAddressPrefix used for generating consensus node address
	ConsNodeAddressPrefix = AccountAddressPrefix + "valcons"
	// ConsNodePubKeyPrefix used for generating consensus node public key
	ConsNodePubKeyPrefix = AccountAddressPrefix + "valconspub"
	// CoinType is the SGE coin type as defined in SLIP44 (https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType uint32 = 909
	// Purpose is the purpose of the BIP44
	Purpose uint32 = 44
	// FullFundraiserPath is the parts of the BIP44 HD path that are fixed by
	// what we used during the SGE fundraiser.
	// More info (https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)
	FullFundraiserPath = "m/44'/909'/0'/0/0"
)

// SetConfig sets prefixes configuration
func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.SetCoinType(CoinType)
	config.SetPurpose(Purpose)
	// nolint
	config.SetFullFundraiserPath(FullFundraiserPath)
	config.Seal()
}
