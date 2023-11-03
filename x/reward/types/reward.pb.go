// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/reward/reward.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// RewardType defines supported types of rewards of reward module.
type RewardType int32

const (
	// the invalid or unknown
	RewardType_REWARD_TYPE_UNSPECIFIED RewardType = 0
	// signup reward
	RewardType_REWARD_TYPE_SIGNUP RewardType = 1
	// referral reward
	RewardType_REWARD_TYPE_REFERRAL RewardType = 2
	// affiliate reward
	RewardType_REWARD_TYPE_AFFILIATE RewardType = 3
	// bet refunds
	RewardType_REWARD_TYPE_BET_REFUND RewardType = 4
	// milestone reward
	RewardType_REWARD_TYPE_MILESTONE RewardType = 5
	// bet discounts
	RewardType_REWARD_TYPE_BET_DISCOUNT RewardType = 6
	// other rewards
	RewardType_REWARD_TYPE_OTHER RewardType = 7
)

var RewardType_name = map[int32]string{
	0: "REWARD_TYPE_UNSPECIFIED",
	1: "REWARD_TYPE_SIGNUP",
	2: "REWARD_TYPE_REFERRAL",
	3: "REWARD_TYPE_AFFILIATE",
	4: "REWARD_TYPE_BET_REFUND",
	5: "REWARD_TYPE_MILESTONE",
	6: "REWARD_TYPE_BET_DISCOUNT",
	7: "REWARD_TYPE_OTHER",
}

var RewardType_value = map[string]int32{
	"REWARD_TYPE_UNSPECIFIED":  0,
	"REWARD_TYPE_SIGNUP":       1,
	"REWARD_TYPE_REFERRAL":     2,
	"REWARD_TYPE_AFFILIATE":    3,
	"REWARD_TYPE_BET_REFUND":   4,
	"REWARD_TYPE_MILESTONE":    5,
	"REWARD_TYPE_BET_DISCOUNT": 6,
	"REWARD_TYPE_OTHER":        7,
}

func (x RewardType) String() string {
	return proto.EnumName(RewardType_name, int32(x))
}

func (RewardType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3412add70a4f177f, []int{0}
}

// RewardType defines supported types of rewards of reward module.
type RewardAmountType int32

const (
	// the invalid or unknown
	RewardAmountType_REWARD_AMOUNT_TYPE_UNSPECIFIED RewardAmountType = 0
	// Fixed amount
	RewardAmountType_REWARD_AMOUNT_TYPE_FIXED RewardAmountType = 1
	// Business logic defined amount
	RewardAmountType_REWARD_AMOUNT_TYPE_LOGIC RewardAmountType = 2
	// Percentage of bet amount
	RewardAmountType_REWARD_AMOUNT_TYPE_PERCENTAGE RewardAmountType = 3
)

var RewardAmountType_name = map[int32]string{
	0: "REWARD_AMOUNT_TYPE_UNSPECIFIED",
	1: "REWARD_AMOUNT_TYPE_FIXED",
	2: "REWARD_AMOUNT_TYPE_LOGIC",
	3: "REWARD_AMOUNT_TYPE_PERCENTAGE",
}

var RewardAmountType_value = map[string]int32{
	"REWARD_AMOUNT_TYPE_UNSPECIFIED": 0,
	"REWARD_AMOUNT_TYPE_FIXED":       1,
	"REWARD_AMOUNT_TYPE_LOGIC":       2,
	"REWARD_AMOUNT_TYPE_PERCENTAGE":  3,
}

func (x RewardAmountType) String() string {
	return proto.EnumName(RewardAmountType_name, int32(x))
}

func (RewardAmountType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3412add70a4f177f, []int{1}
}

// ReceiverType defines all of reward receiver types in the system.
type ReceiverType int32

const (
	// the invalid or unknown
	ReceiverType_RECEIVER_TYPE_UNSPECIFIED ReceiverType = 0
	// single receiver account
	ReceiverType_RECEIVER_TYPE_SINGLE ReceiverType = 1
	// referrer
	ReceiverType_RECEIVER_TYPE_REFERRER ReceiverType = 2
	// referee
	ReceiverType_RECEIVER_TYPE_REFEREE ReceiverType = 3
	// affiliate
	ReceiverType_RECEIVER_TYPE_AFFILIATE ReceiverType = 4
)

var ReceiverType_name = map[int32]string{
	0: "RECEIVER_TYPE_UNSPECIFIED",
	1: "RECEIVER_TYPE_SINGLE",
	2: "RECEIVER_TYPE_REFERRER",
	3: "RECEIVER_TYPE_REFEREE",
	4: "RECEIVER_TYPE_AFFILIATE",
}

var ReceiverType_value = map[string]int32{
	"RECEIVER_TYPE_UNSPECIFIED": 0,
	"RECEIVER_TYPE_SINGLE":      1,
	"RECEIVER_TYPE_REFERRER":    2,
	"RECEIVER_TYPE_REFEREE":     3,
	"RECEIVER_TYPE_AFFILIATE":   4,
}

func (x ReceiverType) String() string {
	return proto.EnumName(ReceiverType_name, int32(x))
}

func (ReceiverType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3412add70a4f177f, []int{2}
}

// Reward is the transaction made to reward a user
// for a specific action.
type Reward struct {
	// creator is the address of the account that invokes the reward transaction.
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	// receiver is the address of the account that receives the reward.
	Receiver string `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	// campaign_uid is the unique identifier of the campaign.
	CampaignUID string `protobuf:"bytes,3,opt,name=campaign_uid,proto3" json:"campaign_uid"`
	// reward_type is the type of the reward.
	RewardType RewardType `protobuf:"varint,4,opt,name=reward_type,proto3,enum=sgenetwork.sge.reward.RewardType" json:"reward_type"`
	// reward_amount is the amount of the reward.
	RewardAmount *RewardAmount `protobuf:"bytes,5,opt,name=reward_amount,proto3" json:"reward_amount"`
	// source is the source of the reward.
	// It is used to identify the source of the reward.
	// For example, the source of a referral signup is Type - referral.
	Source string `protobuf:"bytes,6,opt,name=source,proto3" json:"source"`
	// source_code is the source code of the reward.
	// It is used to identify the source of the reward.
	// For example, the source code of a referral signup is referral code of referer.
	SourceCode string `protobuf:"bytes,7,opt,name=source_code,proto3" json:"source_code"`
	// source_uid is the address of the source.
	// It is used to identify the source of the reward.
	// For example, the source uid of a referral signup reward is the address of the referer.
	SourceUID string `protobuf:"bytes,8,opt,name=source_uid,proto3" json:"source_uid"`
	// created_at is the time when the reward is created.
	CreatedAt uint64 `protobuf:"varint,9,opt,name=created_at,proto3" json:"created_at"`
}

func (m *Reward) Reset()         { *m = Reward{} }
func (m *Reward) String() string { return proto.CompactTextString(m) }
func (*Reward) ProtoMessage()    {}
func (*Reward) Descriptor() ([]byte, []int) {
	return fileDescriptor_3412add70a4f177f, []int{0}
}
func (m *Reward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Reward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Reward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Reward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Reward.Merge(m, src)
}
func (m *Reward) XXX_Size() int {
	return m.Size()
}
func (m *Reward) XXX_DiscardUnknown() {
	xxx_messageInfo_Reward.DiscardUnknown(m)
}

var xxx_messageInfo_Reward proto.InternalMessageInfo

func (m *Reward) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *Reward) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *Reward) GetCampaignUID() string {
	if m != nil {
		return m.CampaignUID
	}
	return ""
}

func (m *Reward) GetRewardType() RewardType {
	if m != nil {
		return m.RewardType
	}
	return RewardType_REWARD_TYPE_UNSPECIFIED
}

func (m *Reward) GetRewardAmount() *RewardAmount {
	if m != nil {
		return m.RewardAmount
	}
	return nil
}

func (m *Reward) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Reward) GetSourceCode() string {
	if m != nil {
		return m.SourceCode
	}
	return ""
}

func (m *Reward) GetSourceUID() string {
	if m != nil {
		return m.SourceUID
	}
	return ""
}

func (m *Reward) GetCreatedAt() uint64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

type RewardAmount struct {
	// main account reward amount
	MainAccountAmount cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=main_account_amount,json=mainAccountAmount,proto3,customtype=cosmossdk.io/math.Int" json:"main_account_amount" yaml:"main_account_amount"`
	// sub account reward amount
	SubaccountAmount cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=subaccount_amount,json=subaccountAmount,proto3,customtype=cosmossdk.io/math.Int" json:"subaccount_amount" yaml:"subaccount_amount"`
	// unlock timestamp
	UnlockTs cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=unlock_ts,json=unlockTs,proto3,customtype=cosmossdk.io/math.Int" json:"unlock_ts" yaml:"unlock_ts"`
}

func (m *RewardAmount) Reset()         { *m = RewardAmount{} }
func (m *RewardAmount) String() string { return proto.CompactTextString(m) }
func (*RewardAmount) ProtoMessage()    {}
func (*RewardAmount) Descriptor() ([]byte, []int) {
	return fileDescriptor_3412add70a4f177f, []int{1}
}
func (m *RewardAmount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RewardAmount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RewardAmount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RewardAmount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardAmount.Merge(m, src)
}
func (m *RewardAmount) XXX_Size() int {
	return m.Size()
}
func (m *RewardAmount) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardAmount.DiscardUnknown(m)
}

var xxx_messageInfo_RewardAmount proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("sgenetwork.sge.reward.RewardType", RewardType_name, RewardType_value)
	proto.RegisterEnum("sgenetwork.sge.reward.RewardAmountType", RewardAmountType_name, RewardAmountType_value)
	proto.RegisterEnum("sgenetwork.sge.reward.ReceiverType", ReceiverType_name, ReceiverType_value)
	proto.RegisterType((*Reward)(nil), "sgenetwork.sge.reward.Reward")
	proto.RegisterType((*RewardAmount)(nil), "sgenetwork.sge.reward.RewardAmount")
}

func init() { proto.RegisterFile("sge/reward/reward.proto", fileDescriptor_3412add70a4f177f) }

var fileDescriptor_3412add70a4f177f = []byte{
	// 768 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0xcd, 0x6e, 0xea, 0x46,
	0x14, 0xc6, 0x84, 0x90, 0x30, 0xa4, 0x95, 0x33, 0x0d, 0x89, 0x43, 0x1a, 0x4c, 0xe8, 0x86, 0x46,
	0xad, 0x51, 0xd3, 0x55, 0x9b, 0x45, 0x65, 0xcc, 0x40, 0x2d, 0x11, 0x83, 0x06, 0xd3, 0xbf, 0x8d,
	0x65, 0x8c, 0xeb, 0x20, 0x02, 0x46, 0xb6, 0x69, 0x9a, 0xb7, 0xe8, 0xa6, 0xab, 0x3e, 0x46, 0xf7,
	0x5d, 0x67, 0x99, 0x65, 0x75, 0x17, 0xd6, 0x95, 0xb3, 0x63, 0x79, 0x9f, 0xe0, 0xca, 0x33, 0x26,
	0x0c, 0x09, 0xf7, 0x66, 0xe5, 0x73, 0xce, 0xf7, 0x7d, 0x73, 0xe6, 0xfc, 0x68, 0x0c, 0x8e, 0x7c,
	0xc7, 0xae, 0x79, 0xf6, 0xad, 0xe9, 0x0d, 0x93, 0x8f, 0x34, 0xf3, 0xdc, 0xc0, 0x85, 0x05, 0xdf,
	0xb1, 0xa7, 0x76, 0x70, 0xeb, 0x7a, 0x63, 0xc9, 0x77, 0x6c, 0x89, 0x82, 0xc5, 0x03, 0xc7, 0x75,
	0x5c, 0xc2, 0xa8, 0xc5, 0x16, 0x25, 0x57, 0xfe, 0xcb, 0x80, 0x2c, 0x26, 0x04, 0x28, 0x80, 0x1d,
	0xcb, 0xb3, 0xcd, 0xc0, 0xf5, 0x04, 0xae, 0xcc, 0x55, 0x73, 0x78, 0xe9, 0xc2, 0x22, 0xd8, 0xf5,
	0x6c, 0xcb, 0x1e, 0xfd, 0x61, 0x7b, 0x42, 0x9a, 0x40, 0x4f, 0x3e, 0x54, 0xc0, 0x9e, 0x65, 0x4e,
	0x66, 0xe6, 0xc8, 0x99, 0x1a, 0xf3, 0xd1, 0x50, 0xd8, 0x8a, 0xf1, 0xba, 0x18, 0x85, 0x62, 0x5e,
	0x49, 0xe2, 0x7d, 0xb5, 0xb1, 0x08, 0xc5, 0x35, 0x1a, 0x5e, 0xf3, 0xe0, 0x00, 0xe4, 0xe9, 0x2d,
	0x8d, 0xe0, 0x6e, 0x66, 0x0b, 0x99, 0x32, 0x57, 0xfd, 0xf4, 0xe2, 0x4c, 0xda, 0x58, 0x88, 0x44,
	0xaf, 0xab, 0xdf, 0xcd, 0xec, 0xfa, 0x69, 0x14, 0x8a, 0x60, 0xe5, 0x2f, 0x42, 0x91, 0x3d, 0x07,
	0xb3, 0x0e, 0xbc, 0x01, 0x9f, 0x24, 0xae, 0x39, 0x71, 0xe7, 0xd3, 0x40, 0xd8, 0x2e, 0x73, 0xd5,
	0xfc, 0xc5, 0x17, 0x1f, 0xcd, 0x22, 0x13, 0x6a, 0xfd, 0x2c, 0x0a, 0xc5, 0x3d, 0x36, 0xb2, 0x08,
	0xc5, 0xf5, 0xd3, 0xf0, 0xba, 0x0b, 0xbf, 0x02, 0x59, 0xdf, 0x9d, 0x7b, 0x96, 0x2d, 0x64, 0x49,
	0x43, 0x0e, 0xa2, 0x50, 0xcc, 0xf6, 0x48, 0x64, 0x11, 0x8a, 0x09, 0x86, 0x93, 0x2f, 0xfc, 0x01,
	0xe4, 0xa9, 0x65, 0x58, 0xee, 0xd0, 0x16, 0x76, 0x88, 0x84, 0x14, 0x47, 0x25, 0x8a, 0x3b, 0x24,
	0xc5, 0x31, 0x24, 0xcc, 0x3a, 0xf0, 0x12, 0x80, 0xc4, 0x8d, 0x67, 0xb0, 0x4b, 0xf4, 0x27, 0x51,
	0x28, 0xe6, 0xa8, 0x9e, 0x4e, 0x80, 0xa1, 0x60, 0xc6, 0x8e, 0xc5, 0x64, 0xd2, 0xf6, 0xd0, 0x30,
	0x03, 0x21, 0x57, 0xe6, 0xaa, 0x19, 0x2a, 0x56, 0x68, 0x54, 0x8e, 0xcb, 0x65, 0x28, 0x98, 0xb1,
	0x2b, 0xff, 0xa6, 0xc1, 0x5a, 0x67, 0xe0, 0x18, 0x7c, 0x36, 0x31, 0x47, 0x53, 0xc3, 0xb4, 0xac,
	0xd8, 0x5f, 0x76, 0x9b, 0xac, 0x54, 0xfd, 0xf2, 0x3e, 0x14, 0x53, 0x6f, 0x42, 0xb1, 0x60, 0xb9,
	0xfe, 0xc4, 0xf5, 0xfd, 0xe1, 0x58, 0x1a, 0xb9, 0xb5, 0x89, 0x19, 0x5c, 0x4b, 0xea, 0x34, 0x78,
	0x17, 0x8a, 0xc5, 0x3b, 0x73, 0x72, 0xf3, 0x7d, 0x65, 0xc3, 0x09, 0x15, 0xbc, 0x1f, 0x47, 0x65,
	0x1a, 0x4c, 0x92, 0xfd, 0x0e, 0xf6, 0xfd, 0xf9, 0xe0, 0x59, 0x2a, 0xb2, 0xa2, 0xf5, 0xef, 0x5e,
	0x4b, 0x25, 0xd0, 0x54, 0x2f, 0xf4, 0x15, 0xcc, 0xaf, 0x62, 0x49, 0x1e, 0x0d, 0xe4, 0xe6, 0xd3,
	0x1b, 0xd7, 0x1a, 0x1b, 0x81, 0x9f, 0xac, 0xf8, 0x37, 0xaf, 0x9d, 0xcf, 0xd3, 0xf3, 0x9f, 0x74,
	0x15, 0xbc, 0x4b, 0x6d, 0xdd, 0x3f, 0x8f, 0x38, 0xc0, 0xec, 0x2d, 0x3c, 0x01, 0x47, 0x18, 0xfd,
	0x2c, 0xe3, 0x86, 0xa1, 0xff, 0xda, 0x45, 0x46, 0x5f, 0xeb, 0x75, 0x91, 0xa2, 0x36, 0x55, 0xd4,
	0xe0, 0x53, 0xf0, 0x10, 0x40, 0x16, 0xec, 0xa9, 0x2d, 0xad, 0xdf, 0xe5, 0x39, 0x28, 0x80, 0x03,
	0x36, 0x8e, 0x51, 0x13, 0x61, 0x2c, 0xb7, 0xf9, 0x34, 0x3c, 0x06, 0x05, 0x16, 0x91, 0x9b, 0x4d,
	0xb5, 0xad, 0xca, 0x3a, 0xe2, 0xb7, 0x60, 0x11, 0x1c, 0xb2, 0x50, 0x1d, 0xe9, 0xb1, 0xb0, 0xaf,
	0x35, 0xf8, 0xcc, 0x73, 0xd9, 0x95, 0xda, 0x46, 0x3d, 0xbd, 0xa3, 0x21, 0x7e, 0x1b, 0x7e, 0x0e,
	0x84, 0xe7, 0xb2, 0x86, 0xda, 0x53, 0x3a, 0x7d, 0x4d, 0xe7, 0xb3, 0xb0, 0x00, 0xf6, 0x59, 0xb4,
	0xa3, 0xff, 0x88, 0x30, 0xbf, 0x73, 0xfe, 0x37, 0x07, 0x78, 0x76, 0x35, 0x48, 0xa9, 0x15, 0x50,
	0x4a, 0xb8, 0xf2, 0x55, 0x2c, 0xdf, 0x54, 0xf1, 0x2a, 0x1b, 0xcb, 0x69, 0xaa, 0xbf, 0xa0, 0x06,
	0xcf, 0x7d, 0x00, 0x6d, 0x77, 0x5a, 0xaa, 0xc2, 0xa7, 0xe1, 0x19, 0x38, 0xdd, 0x80, 0x76, 0x11,
	0x56, 0x90, 0xa6, 0xcb, 0x2d, 0xc4, 0x6f, 0x9d, 0xff, 0xc3, 0xc5, 0x2b, 0x4b, 0xdf, 0x2f, 0x72,
	0xa7, 0x53, 0x70, 0x8c, 0x91, 0x82, 0xd4, 0x9f, 0x10, 0xde, 0x74, 0x1d, 0xd2, 0x68, 0x16, 0xee,
	0xa9, 0x5a, 0xab, 0x8d, 0x78, 0x8e, 0x76, 0x93, 0x45, 0xe8, 0x10, 0x10, 0x5e, 0x0e, 0xe1, 0x05,
	0x86, 0xe2, 0x21, 0x90, 0x71, 0xb3, 0xd0, 0x6a, 0x42, 0x99, 0xba, 0x72, 0x1f, 0x95, 0xb8, 0x87,
	0xa8, 0xc4, 0xbd, 0x8d, 0x4a, 0xdc, 0x5f, 0x8f, 0xa5, 0xd4, 0xc3, 0x63, 0x29, 0xf5, 0xff, 0x63,
	0x29, 0xf5, 0xdb, 0x97, 0xce, 0x28, 0xb8, 0x9e, 0x0f, 0x24, 0xcb, 0x9d, 0xd4, 0x7c, 0xc7, 0xfe,
	0x3a, 0x79, 0xb5, 0x62, 0xbb, 0xf6, 0xe7, 0xf2, 0x57, 0x10, 0xbf, 0x75, 0xfe, 0x20, 0x4b, 0x5e,
	0xf7, 0x6f, 0xdf, 0x07, 0x00, 0x00, 0xff, 0xff, 0x65, 0xfe, 0x43, 0x0d, 0x25, 0x06, 0x00, 0x00,
}

func (m *Reward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Reward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Reward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CreatedAt != 0 {
		i = encodeVarintReward(dAtA, i, uint64(m.CreatedAt))
		i--
		dAtA[i] = 0x48
	}
	if len(m.SourceUID) > 0 {
		i -= len(m.SourceUID)
		copy(dAtA[i:], m.SourceUID)
		i = encodeVarintReward(dAtA, i, uint64(len(m.SourceUID)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.SourceCode) > 0 {
		i -= len(m.SourceCode)
		copy(dAtA[i:], m.SourceCode)
		i = encodeVarintReward(dAtA, i, uint64(len(m.SourceCode)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Source) > 0 {
		i -= len(m.Source)
		copy(dAtA[i:], m.Source)
		i = encodeVarintReward(dAtA, i, uint64(len(m.Source)))
		i--
		dAtA[i] = 0x32
	}
	if m.RewardAmount != nil {
		{
			size, err := m.RewardAmount.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintReward(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.RewardType != 0 {
		i = encodeVarintReward(dAtA, i, uint64(m.RewardType))
		i--
		dAtA[i] = 0x20
	}
	if len(m.CampaignUID) > 0 {
		i -= len(m.CampaignUID)
		copy(dAtA[i:], m.CampaignUID)
		i = encodeVarintReward(dAtA, i, uint64(len(m.CampaignUID)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintReward(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintReward(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RewardAmount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RewardAmount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RewardAmount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.UnlockTs.Size()
		i -= size
		if _, err := m.UnlockTs.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintReward(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.SubaccountAmount.Size()
		i -= size
		if _, err := m.SubaccountAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintReward(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.MainAccountAmount.Size()
		i -= size
		if _, err := m.MainAccountAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintReward(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintReward(dAtA []byte, offset int, v uint64) int {
	offset -= sovReward(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Reward) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovReward(uint64(l))
	}
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovReward(uint64(l))
	}
	l = len(m.CampaignUID)
	if l > 0 {
		n += 1 + l + sovReward(uint64(l))
	}
	if m.RewardType != 0 {
		n += 1 + sovReward(uint64(m.RewardType))
	}
	if m.RewardAmount != nil {
		l = m.RewardAmount.Size()
		n += 1 + l + sovReward(uint64(l))
	}
	l = len(m.Source)
	if l > 0 {
		n += 1 + l + sovReward(uint64(l))
	}
	l = len(m.SourceCode)
	if l > 0 {
		n += 1 + l + sovReward(uint64(l))
	}
	l = len(m.SourceUID)
	if l > 0 {
		n += 1 + l + sovReward(uint64(l))
	}
	if m.CreatedAt != 0 {
		n += 1 + sovReward(uint64(m.CreatedAt))
	}
	return n
}

func (m *RewardAmount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MainAccountAmount.Size()
	n += 1 + l + sovReward(uint64(l))
	l = m.SubaccountAmount.Size()
	n += 1 + l + sovReward(uint64(l))
	l = m.UnlockTs.Size()
	n += 1 + l + sovReward(uint64(l))
	return n
}

func sovReward(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozReward(x uint64) (n int) {
	return sovReward(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Reward) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReward
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Reward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Reward: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CampaignUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CampaignUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardType", wireType)
			}
			m.RewardType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardType |= RewardType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RewardAmount == nil {
				m.RewardAmount = &RewardAmount{}
			}
			if err := m.RewardAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Source = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceCode", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceCode = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipReward(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReward
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RewardAmount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReward
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RewardAmount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RewardAmount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MainAccountAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MainAccountAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubaccountAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SubaccountAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnlockTs", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReward
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReward
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.UnlockTs.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReward(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReward
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipReward(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowReward
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowReward
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowReward
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthReward
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupReward
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthReward
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthReward        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowReward          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupReward = fmt.Errorf("proto: unexpected end of group")
)
