package types

import (
	"encoding/binary"

	"github.com/sge-network/sge/utils"
)

var _ binary.ByteOrder

var (
	// PromoterKeyPrefix is the prefix to retrieve all promoter
	PromoterKeyPrefix = []byte{0x04}

	// PromoterAddressKeyPrefix is the prefix to retrieve all addresses promoter
	PromoterAddressKeyPrefix = []byte{0x05}
)

// GetPromoterKey returns the store key to retrieve a promoter
func GetPromoterKey(uid string) []byte {
	return utils.StrBytes(uid)
}

// GetPromoterByAddressKey returns the promoter of a certain promoter address.
func GetPromoterByAddressKey(accAddr string) []byte {
	return utils.StrBytes(accAddr)
}
