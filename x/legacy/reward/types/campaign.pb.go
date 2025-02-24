// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sgenetwork/sge/reward/campaign.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// Campaign is type for defining the campaign properties.
type Campaign struct {
	// creator is the address of campaign creator.
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	// uid is the unique identifier of a campaign.
	UID string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid"`
	// promoter is the address of campaign promoter.
	// Funds for the campaign would be deducted from this account.
	Promoter string `protobuf:"bytes,3,opt,name=promoter,proto3" json:"promoter,omitempty"`
	// start_ts is the start timestamp of a campaign.
	StartTS uint64 `protobuf:"varint,4,opt,name=start_ts,proto3" json:"start_ts"`
	// end_ts is the end timestamp of a campaign.
	EndTS uint64 `protobuf:"varint,5,opt,name=end_ts,proto3" json:"end_ts"`
	// reward_category is the category of reward.
	RewardCategory RewardCategory `protobuf:"varint,6,opt,name=reward_category,json=rewardCategory,proto3,enum=sgenetwork.sge.reward.RewardCategory" json:"reward_category,omitempty"`
	// reward_type is the type of reward.
	RewardType RewardType `protobuf:"varint,7,opt,name=reward_type,json=rewardType,proto3,enum=sgenetwork.sge.reward.RewardType" json:"reward_type,omitempty"`
	// amount_type is the type of reward amount.
	RewardAmountType RewardAmountType `protobuf:"varint,8,opt,name=reward_amount_type,json=rewardAmountType,proto3,enum=sgenetwork.sge.reward.RewardAmountType" json:"reward_amount_type,omitempty"`
	// reward_amount is the amount defined for a reward.
	RewardAmount *RewardAmount `protobuf:"bytes,9,opt,name=reward_amount,json=rewardAmount,proto3" json:"reward_amount,omitempty"`
	// pool is the tracker of campaign funds.
	Pool Pool `protobuf:"bytes,10,opt,name=pool,proto3" json:"pool"`
	// is_active is the flag to check if the campaign is active or not.
	IsActive bool `protobuf:"varint,11,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	// meta is the metadata of the campaign.
	// It is a stringified base64 encoded json.
	Meta string `protobuf:"bytes,13,opt,name=meta,proto3" json:"meta,omitempty"`
	// cap_count is the maximum allowed grant for a certain account.
	CapCount uint64 `protobuf:"varint,14,opt,name=cap_count,json=capCount,proto3" json:"cap_count,omitempty"`
	// constraints is the constrains of a campaign.
	Constraints *CampaignConstraints `protobuf:"bytes,15,opt,name=constraints,proto3" json:"constraints,omitempty"`
}

func (m *Campaign) Reset()         { *m = Campaign{} }
func (m *Campaign) String() string { return proto.CompactTextString(m) }
func (*Campaign) ProtoMessage()    {}
func (*Campaign) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2342810052fed89, []int{0}
}
func (m *Campaign) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Campaign) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Campaign.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Campaign) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Campaign.Merge(m, src)
}
func (m *Campaign) XXX_Size() int {
	return m.Size()
}
func (m *Campaign) XXX_DiscardUnknown() {
	xxx_messageInfo_Campaign.DiscardUnknown(m)
}

var xxx_messageInfo_Campaign proto.InternalMessageInfo

func (m *Campaign) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *Campaign) GetUID() string {
	if m != nil {
		return m.UID
	}
	return ""
}

func (m *Campaign) GetPromoter() string {
	if m != nil {
		return m.Promoter
	}
	return ""
}

func (m *Campaign) GetStartTS() uint64 {
	if m != nil {
		return m.StartTS
	}
	return 0
}

func (m *Campaign) GetEndTS() uint64 {
	if m != nil {
		return m.EndTS
	}
	return 0
}

func (m *Campaign) GetRewardCategory() RewardCategory {
	if m != nil {
		return m.RewardCategory
	}
	return RewardCategory_REWARD_CATEGORY_UNSPECIFIED
}

func (m *Campaign) GetRewardType() RewardType {
	if m != nil {
		return m.RewardType
	}
	return RewardType_REWARD_TYPE_UNSPECIFIED
}

func (m *Campaign) GetRewardAmountType() RewardAmountType {
	if m != nil {
		return m.RewardAmountType
	}
	return RewardAmountType_REWARD_AMOUNT_TYPE_UNSPECIFIED
}

func (m *Campaign) GetRewardAmount() *RewardAmount {
	if m != nil {
		return m.RewardAmount
	}
	return nil
}

func (m *Campaign) GetPool() Pool {
	if m != nil {
		return m.Pool
	}
	return Pool{}
}

func (m *Campaign) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *Campaign) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func (m *Campaign) GetCapCount() uint64 {
	if m != nil {
		return m.CapCount
	}
	return 0
}

func (m *Campaign) GetConstraints() *CampaignConstraints {
	if m != nil {
		return m.Constraints
	}
	return nil
}

// Pool tracks funds assigned and spent to/for a campaign.
type Pool struct {
	Total     cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=total,proto3,customtype=cosmossdk.io/math.Int" json:"total" yaml:"total"`
	Spent     cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=spent,proto3,customtype=cosmossdk.io/math.Int" json:"spent" yaml:"spent"`
	Withdrawn cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=withdrawn,proto3,customtype=cosmossdk.io/math.Int" json:"withdrawn" yaml:"spent"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2342810052fed89, []int{1}
}
func (m *Pool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pool.Merge(m, src)
}
func (m *Pool) XXX_Size() int {
	return m.Size()
}
func (m *Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_Pool proto.InternalMessageInfo

// CampaignConstraints contains campaign constraints and criteria.
type CampaignConstraints struct {
	MaxBetAmount cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=max_bet_amount,json=maxBetAmount,proto3,customtype=cosmossdk.io/math.Int" json:"max_bet_amount" yaml:"max_bet_amount"`
}

func (m *CampaignConstraints) Reset()         { *m = CampaignConstraints{} }
func (m *CampaignConstraints) String() string { return proto.CompactTextString(m) }
func (*CampaignConstraints) ProtoMessage()    {}
func (*CampaignConstraints) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2342810052fed89, []int{2}
}
func (m *CampaignConstraints) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CampaignConstraints) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CampaignConstraints.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CampaignConstraints) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CampaignConstraints.Merge(m, src)
}
func (m *CampaignConstraints) XXX_Size() int {
	return m.Size()
}
func (m *CampaignConstraints) XXX_DiscardUnknown() {
	xxx_messageInfo_CampaignConstraints.DiscardUnknown(m)
}

var xxx_messageInfo_CampaignConstraints proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Campaign)(nil), "sgenetwork.sge.reward.Campaign")
	proto.RegisterType((*Pool)(nil), "sgenetwork.sge.reward.Pool")
	proto.RegisterType((*CampaignConstraints)(nil), "sgenetwork.sge.reward.CampaignConstraints")
}

func init() {
	proto.RegisterFile("sgenetwork/sge/reward/campaign.proto", fileDescriptor_e2342810052fed89)
}

var fileDescriptor_e2342810052fed89 = []byte{
	// 660 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0x4f, 0x6f, 0xd3, 0x3e,
	0x18, 0xc7, 0x9b, 0xdf, 0xba, 0x2e, 0x75, 0xb7, 0xee, 0x27, 0x6f, 0x13, 0x66, 0x93, 0x9a, 0x12,
	0x40, 0x54, 0xa0, 0xa5, 0x68, 0x03, 0x0e, 0xdc, 0xd6, 0x82, 0xc4, 0x00, 0x21, 0xe4, 0x6d, 0x17,
	0x84, 0x54, 0x79, 0x89, 0x95, 0x45, 0x6b, 0xe2, 0xc8, 0xf6, 0xe8, 0x7a, 0xe3, 0x25, 0xf0, 0x62,
	0x78, 0x11, 0x3b, 0x4e, 0x9c, 0x10, 0x12, 0x11, 0xca, 0x6e, 0x3b, 0xf2, 0x0a, 0x90, 0xed, 0xb4,
	0xdb, 0x50, 0xf7, 0x87, 0x4b, 0x6d, 0x3f, 0xf9, 0x7e, 0x3f, 0xcf, 0xd3, 0xe4, 0xf1, 0x03, 0xee,
	0x89, 0x90, 0x26, 0x54, 0x0e, 0x18, 0xdf, 0x6f, 0x8b, 0x90, 0xb6, 0x39, 0x1d, 0x10, 0x1e, 0xb4,
	0x7d, 0x12, 0xa7, 0x24, 0x0a, 0x13, 0x2f, 0xe5, 0x4c, 0x32, 0xb8, 0x74, 0xa6, 0xf2, 0x44, 0x48,
	0x3d, 0xa3, 0x5a, 0xbe, 0xed, 0x33, 0x11, 0x33, 0xd1, 0xd3, 0xa2, 0xb6, 0x39, 0x18, 0xc7, 0xf2,
	0x62, 0xc8, 0x42, 0x66, 0xe2, 0x6a, 0x57, 0x44, 0xdd, 0xc9, 0xd9, 0xcc, 0x62, 0x34, 0xee, 0xe7,
	0x0a, 0xb0, 0xbb, 0x45, 0x7a, 0xb8, 0x06, 0x66, 0x7c, 0x4e, 0x89, 0x64, 0x1c, 0x59, 0x4d, 0xab,
	0x55, 0xed, 0xa0, 0x6f, 0x5f, 0x57, 0x17, 0x8b, 0x4c, 0x1b, 0x41, 0xc0, 0xa9, 0x10, 0x5b, 0x92,
	0x47, 0x49, 0x88, 0x47, 0x42, 0xd8, 0x04, 0x53, 0x07, 0x51, 0x80, 0xfe, 0xd3, 0xfa, 0x7a, 0x9e,
	0x39, 0x53, 0x3b, 0x9b, 0x2f, 0x4e, 0x33, 0x47, 0x45, 0xb1, 0xfa, 0x81, 0x4f, 0x80, 0x9d, 0x72,
	0x16, 0x33, 0x49, 0x39, 0x9a, 0xba, 0x06, 0x3b, 0x56, 0xc2, 0x75, 0x60, 0x0b, 0x49, 0xb8, 0xec,
	0x49, 0x81, 0xca, 0x4d, 0xab, 0x55, 0xee, 0xdc, 0xca, 0x33, 0x67, 0x66, 0x4b, 0xc5, 0xb6, 0xb7,
	0x4e, 0x33, 0x67, 0xfc, 0x18, 0x8f, 0x77, 0xf0, 0x11, 0xa8, 0xd0, 0x24, 0x50, 0x96, 0x69, 0x6d,
	0x59, 0xc8, 0x33, 0x67, 0xfa, 0x65, 0x12, 0x68, 0x43, 0xf1, 0x08, 0x17, 0x2b, 0x7c, 0x07, 0xe6,
	0xcd, 0xab, 0xe8, 0xf9, 0x44, 0xd2, 0x90, 0xf1, 0x21, 0xaa, 0x34, 0xad, 0x56, 0x7d, 0xed, 0xbe,
	0x37, 0xf1, 0x03, 0x78, 0x58, 0x2f, 0xdd, 0x42, 0x8c, 0xeb, 0xfc, 0xc2, 0x19, 0x76, 0x40, 0xad,
	0xe0, 0xc9, 0x61, 0x4a, 0xd1, 0x8c, 0x66, 0xdd, 0xb9, 0x92, 0xb5, 0x3d, 0x4c, 0x29, 0x06, 0x7c,
	0xbc, 0x87, 0x3b, 0x00, 0x16, 0x0c, 0x12, 0xb3, 0x83, 0x44, 0x1a, 0x94, 0xad, 0x51, 0x0f, 0xae,
	0x44, 0x6d, 0x68, 0xbd, 0x06, 0xfe, 0xcf, 0xff, 0x8a, 0xc0, 0x57, 0x60, 0xee, 0x02, 0x16, 0x55,
	0x9b, 0x56, 0xab, 0xb6, 0x76, 0xf7, 0x06, 0x44, 0x3c, 0x7b, 0x9e, 0x06, 0x9f, 0x82, 0x72, 0xca,
	0x58, 0x1f, 0x01, 0x0d, 0x58, 0xb9, 0x04, 0xf0, 0x9e, 0xb1, 0x7e, 0xa7, 0x7c, 0x94, 0x39, 0x25,
	0xac, 0xe5, 0x70, 0x05, 0x54, 0x23, 0xd1, 0x23, 0xbe, 0x8c, 0x3e, 0x51, 0x54, 0x6b, 0x5a, 0x2d,
	0x1b, 0xdb, 0x91, 0xd8, 0xd0, 0x67, 0x08, 0x41, 0x39, 0xa6, 0x92, 0xa0, 0x39, 0xd5, 0x1c, 0x58,
	0xef, 0x95, 0xc1, 0x27, 0x69, 0xcf, 0xd7, 0xd5, 0xd6, 0xd5, 0xc7, 0xc4, 0xb6, 0x4f, 0xd2, 0xae,
	0x2e, 0xe2, 0x2d, 0xa8, 0xf9, 0x2c, 0x11, 0x92, 0x93, 0x28, 0x91, 0x02, 0xcd, 0xeb, 0x5a, 0x1e,
	0x5e, 0x52, 0xcb, 0xa8, 0xbb, 0xbb, 0x67, 0x0e, 0x7c, 0xde, 0xee, 0xfe, 0xb4, 0x40, 0x59, 0x15,
	0x0c, 0xbb, 0x60, 0x5a, 0x32, 0x49, 0xfa, 0x45, 0xf3, 0xaf, 0xaa, 0xfa, 0x7f, 0x64, 0xce, 0x92,
	0xe9, 0x54, 0x11, 0xec, 0x7b, 0x11, 0x6b, 0xc7, 0x44, 0xee, 0x79, 0x9b, 0x89, 0xfc, 0x9d, 0x39,
	0xb3, 0x43, 0x12, 0xf7, 0x9f, 0xbb, 0xda, 0xe3, 0x62, 0xe3, 0x55, 0x10, 0x91, 0xd2, 0x44, 0x16,
	0x37, 0xe2, 0xa6, 0x10, 0xed, 0x71, 0xb1, 0xf1, 0xc2, 0x37, 0xa0, 0x3a, 0x88, 0xe4, 0x5e, 0xc0,
	0xc9, 0x20, 0x29, 0xee, 0xcc, 0x3f, 0x82, 0xce, 0xfc, 0xae, 0x00, 0x0b, 0x13, 0xde, 0x01, 0xfc,
	0x08, 0xea, 0x31, 0x39, 0xec, 0xed, 0x52, 0x39, 0x6a, 0x0a, 0xf3, 0xb7, 0x9f, 0x5d, 0x97, 0x68,
	0xc9, 0x24, 0xba, 0x68, 0x76, 0xf1, 0x6c, 0x4c, 0x0e, 0x3b, 0x54, 0x9a, 0x3e, 0xe9, 0xbc, 0x3e,
	0xca, 0x1b, 0xd6, 0x71, 0xde, 0xb0, 0x7e, 0xe5, 0x0d, 0xeb, 0xcb, 0x49, 0xa3, 0x74, 0x7c, 0xd2,
	0x28, 0x7d, 0x3f, 0x69, 0x94, 0x3e, 0x3c, 0x0e, 0x23, 0xb9, 0x77, 0xb0, 0xeb, 0xf9, 0x2c, 0x56,
	0x53, 0x69, 0xf5, 0xfc, 0x84, 0x3a, 0x6c, 0xf7, 0x69, 0x48, 0xfc, 0xe1, 0x68, 0x54, 0xa9, 0xee,
	0x17, 0xbb, 0x15, 0x3d, 0xaa, 0xd6, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x84, 0x9b, 0x95, 0x09,
	0x3e, 0x05, 0x00, 0x00,
}

func (m *Campaign) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Campaign) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Campaign) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Constraints != nil {
		{
			size, err := m.Constraints.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCampaign(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x7a
	}
	if m.CapCount != 0 {
		i = encodeVarintCampaign(dAtA, i, uint64(m.CapCount))
		i--
		dAtA[i] = 0x70
	}
	if len(m.Meta) > 0 {
		i -= len(m.Meta)
		copy(dAtA[i:], m.Meta)
		i = encodeVarintCampaign(dAtA, i, uint64(len(m.Meta)))
		i--
		dAtA[i] = 0x6a
	}
	if m.IsActive {
		i--
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	{
		size, err := m.Pool.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintCampaign(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if m.RewardAmount != nil {
		{
			size, err := m.RewardAmount.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCampaign(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if m.RewardAmountType != 0 {
		i = encodeVarintCampaign(dAtA, i, uint64(m.RewardAmountType))
		i--
		dAtA[i] = 0x40
	}
	if m.RewardType != 0 {
		i = encodeVarintCampaign(dAtA, i, uint64(m.RewardType))
		i--
		dAtA[i] = 0x38
	}
	if m.RewardCategory != 0 {
		i = encodeVarintCampaign(dAtA, i, uint64(m.RewardCategory))
		i--
		dAtA[i] = 0x30
	}
	if m.EndTS != 0 {
		i = encodeVarintCampaign(dAtA, i, uint64(m.EndTS))
		i--
		dAtA[i] = 0x28
	}
	if m.StartTS != 0 {
		i = encodeVarintCampaign(dAtA, i, uint64(m.StartTS))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Promoter) > 0 {
		i -= len(m.Promoter)
		copy(dAtA[i:], m.Promoter)
		i = encodeVarintCampaign(dAtA, i, uint64(len(m.Promoter)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.UID) > 0 {
		i -= len(m.UID)
		copy(dAtA[i:], m.UID)
		i = encodeVarintCampaign(dAtA, i, uint64(len(m.UID)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintCampaign(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Pool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Withdrawn.Size()
		i -= size
		if _, err := m.Withdrawn.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCampaign(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Spent.Size()
		i -= size
		if _, err := m.Spent.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCampaign(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Total.Size()
		i -= size
		if _, err := m.Total.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCampaign(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *CampaignConstraints) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CampaignConstraints) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CampaignConstraints) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MaxBetAmount.Size()
		i -= size
		if _, err := m.MaxBetAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCampaign(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintCampaign(dAtA []byte, offset int, v uint64) int {
	offset -= sovCampaign(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Campaign) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovCampaign(uint64(l))
	}
	l = len(m.UID)
	if l > 0 {
		n += 1 + l + sovCampaign(uint64(l))
	}
	l = len(m.Promoter)
	if l > 0 {
		n += 1 + l + sovCampaign(uint64(l))
	}
	if m.StartTS != 0 {
		n += 1 + sovCampaign(uint64(m.StartTS))
	}
	if m.EndTS != 0 {
		n += 1 + sovCampaign(uint64(m.EndTS))
	}
	if m.RewardCategory != 0 {
		n += 1 + sovCampaign(uint64(m.RewardCategory))
	}
	if m.RewardType != 0 {
		n += 1 + sovCampaign(uint64(m.RewardType))
	}
	if m.RewardAmountType != 0 {
		n += 1 + sovCampaign(uint64(m.RewardAmountType))
	}
	if m.RewardAmount != nil {
		l = m.RewardAmount.Size()
		n += 1 + l + sovCampaign(uint64(l))
	}
	l = m.Pool.Size()
	n += 1 + l + sovCampaign(uint64(l))
	if m.IsActive {
		n += 2
	}
	l = len(m.Meta)
	if l > 0 {
		n += 1 + l + sovCampaign(uint64(l))
	}
	if m.CapCount != 0 {
		n += 1 + sovCampaign(uint64(m.CapCount))
	}
	if m.Constraints != nil {
		l = m.Constraints.Size()
		n += 1 + l + sovCampaign(uint64(l))
	}
	return n
}

func (m *Pool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Total.Size()
	n += 1 + l + sovCampaign(uint64(l))
	l = m.Spent.Size()
	n += 1 + l + sovCampaign(uint64(l))
	l = m.Withdrawn.Size()
	n += 1 + l + sovCampaign(uint64(l))
	return n
}

func (m *CampaignConstraints) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MaxBetAmount.Size()
	n += 1 + l + sovCampaign(uint64(l))
	return n
}

func sovCampaign(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCampaign(x uint64) (n int) {
	return sovCampaign(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Campaign) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCampaign
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
			return fmt.Errorf("proto: Campaign: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Campaign: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Promoter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Promoter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTS", wireType)
			}
			m.StartTS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartTS |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTS", wireType)
			}
			m.EndTS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndTS |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardCategory", wireType)
			}
			m.RewardCategory = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardCategory |= RewardCategory(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardType", wireType)
			}
			m.RewardType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAmountType", wireType)
			}
			m.RewardAmountType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardAmountType |= RewardAmountType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
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
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Pool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsActive = bool(v != 0)
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Meta = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 14:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CapCount", wireType)
			}
			m.CapCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CapCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Constraints", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Constraints == nil {
				m.Constraints = &CampaignConstraints{}
			}
			if err := m.Constraints.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCampaign(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCampaign
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
func (m *Pool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCampaign
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
			return fmt.Errorf("proto: Pool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Total", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Total.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Spent", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Spent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Withdrawn", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Withdrawn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCampaign(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCampaign
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
func (m *CampaignConstraints) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCampaign
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
			return fmt.Errorf("proto: CampaignConstraints: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CampaignConstraints: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBetAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCampaign
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
				return ErrInvalidLengthCampaign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCampaign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxBetAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCampaign(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCampaign
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
func skipCampaign(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCampaign
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
					return 0, ErrIntOverflowCampaign
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
					return 0, ErrIntOverflowCampaign
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
				return 0, ErrInvalidLengthCampaign
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCampaign
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCampaign
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCampaign        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCampaign          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCampaign = fmt.Errorf("proto: unexpected end of group")
)
