package types

import (
	"encoding/binary"

	"github.com/sge-network/sge/utils"
)

var _ binary.ByteOrder

// CampaignKeyPrefix is the prefix to retrieve all Campaign
var CampaignKeyPrefix = []byte{0x00}

// GetCampaignKey returns the store key to retrieve a Campaign from the index fields
func GetCampaignKey(uid string) []byte {
	return utils.StrBytes(uid)
}
