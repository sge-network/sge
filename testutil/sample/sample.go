package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccAddressAsString returns a sample account address as string
func AccAddressAsString() string {
	return AccAddress().String()
}

func AccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

// ValAddress returns a sample validator address
func ValAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr, _ := sdk.Bech32ifyAddressBytes("cosmosvaloper", pk.Address().Bytes())
	return addr
}
