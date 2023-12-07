package types

import (
	"encoding/binary"

	"github.com/sge-network/sge/utils"
)

var _ binary.ByteOrder

var (
	// RewardKeyPrefix is the prefix to retrieve all applied rewards
	RewardKeyPrefix = []byte{0x01}

	// RewardByReceiverAndCategoryKeyPrefix is the prefix to retrieve all applied rewards for a certain receiver account.
	RewardByReceiverAndCategoryKeyPrefix = []byte{0x02}

	// RewardByCampaignKeyPrefix is the prefix to retrieve all applied rewards for a certain campaign.
	RewardByCampaignKeyPrefix = []byte{0x03}
)

// GetRewardsByAccPrefix returns the store key to retrieve list of all applied rewards of a certain campaign
// this should be used with RewardByReceiverKeyPrefix
func GetRewardsByAccPrefix(receiverAcc string) []byte {
	return utils.StrBytes(receiverAcc)
}

// GetRewardsByCategoryPrefix returns the store key to retrieve list of all applied rewards of certain address and category
func GetRewardsByCategoryPrefix(receiverAcc string, rewardCategory RewardCategory) []byte {
	return append(GetRewardsByAccPrefix(receiverAcc), utils.Int32ToBytes(int32(rewardCategory))...)
}

// GetRewardsByCategoryKey returns the store key to retrieve list of applied reward of certain address and category
func GetRewardsByCategoryKey(receiverAcc string, rewardCategory RewardCategory, uid string) []byte {
	return append(GetRewardsByCategoryPrefix(receiverAcc, rewardCategory), utils.StrBytes(uid)...)
}

// GetRewardKey returns the store key to retrieve a certain reward.
func GetRewardKey(uid string) []byte {
	return utils.StrBytes(uid)
}

// GetRewardsByCampaignPrefix returns the store key to retrieve list of all applied rewards of a certain campaign
// this should be used with RewardKeyPrefix
func GetRewardsByCampaignPrefix(campaignUID string) []byte {
	return utils.StrBytes(campaignUID)
}

// GetRewardsByCampaignKey returns the store key to retrieve applied reward of a certain campaign
// this should be used with RewardKeyPrefix
func GetRewardsByCampaignKey(campaignUID, uid string) []byte {
	return append(utils.StrBytes(campaignUID), utils.StrBytes(uid)...)
}
