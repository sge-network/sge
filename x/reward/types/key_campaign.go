package types

import (
	"encoding/binary"

	"github.com/sge-network/sge/utils"
)

var _ binary.ByteOrder

// CampaignKeyPrefix is the prefix to retrieve all Campaign
var CampaignKeyPrefix = []byte{0x00}

// GetHouseCampaign returns the store key to retrieve list of Campaigns of a house
func GetHouseCampaign(houseUID string) []byte {
	return utils.StrBytes(houseUID)
}

// GetCampaignKey returns the store key to retrieve a Campaign from the index fields
func GetCampaignKey(houseUID, uid string) []byte {
	return append(GetHouseCampaign(houseUID), utils.StrBytes(uid)...)
}
