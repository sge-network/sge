package simapp

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// PKs is a slice of public keys for test
var PKs = createTestPubKeys(500)

func CreateJwtTicket(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claim)
	return token.SignedString(TestOVMPrivateKeys[0])
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := cast.ToString(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
		bech := res.String()
		addr, _ := testAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

// testAddr returns sample account address
func testAddr(addr, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHexUnsafe(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// createTestPubKeys returns a total of numPubKeys public keys in ascending order.
func createTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	var publicKeys []cryptotypes.PubKey
	var buffer bytes.Buffer

	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := cast.ToString(i)
		buffer.WriteString(
			"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF",
		) // base pubkey string
		buffer.WriteString(
			numString,
		) // adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, newPubKeyFromHex(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

// newPubKeyFromHex returns a PubKey from a hex string.
func newPubKeyFromHex(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	if len(pkBytes) != ed25519.PubKeySize {
		panic(sdkerrors.Wrap(sdkerrtypes.ErrInvalidPubKey, "invalid pubkey size"))
	}
	return &ed25519.PubKey{Key: pkBytes}
}

func RandomString(length int) string {
	buff := make([]byte, int(math.Ceil(float64(length)/float64(1.33333333333))))
	_, err := rand.Read(buff)
	if err != nil {
		panic(err)
	}
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:length] // strip 1 extra character we get from odd length results
}

func WriteKeyringFile(name string, dir string, contents []byte) error {
	file := filepath.Join(dir, name)

	//#nosec
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("could not create directory %q: %w", dir, err)
	}

	//#nosec
	if err := os.WriteFile(file, contents, 0o644); err != nil { //nolint: gosec
		return err
	}

	return nil
}
