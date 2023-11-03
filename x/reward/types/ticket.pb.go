// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/reward/ticket.proto

package types

import (
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

// CreateCampaignPayload is the type for campaign creation payload.
type CreateCampaignPayload struct {
	// promoter is the address of campaign promoter.
	// Funds for the campaign would be deducted from this account.
	Promoter string `protobuf:"bytes,1,opt,name=promoter,proto3" json:"promoter,omitempty"`
	// start_ts is the start timestamp of the campaign.
	StartTs uint64 `protobuf:"varint,2,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	// end_ts is the end timestamp of the campaign.
	EndTs uint64 `protobuf:"varint,3,opt,name=end_ts,json=endTs,proto3" json:"end_ts,omitempty"`
	// reward_type is the type of reward.
	RewardType RewardType `protobuf:"varint,4,opt,name=reward_type,json=rewardType,proto3,enum=sgenetwork.sge.reward.RewardType" json:"reward_type,omitempty"`
	// Reward amount
	RewardAmountType RewardAmountType `protobuf:"varint,7,opt,name=reward_amount_type,json=rewardAmountType,proto3,enum=sgenetwork.sge.reward.RewardAmountType" json:"reward_amount_type,omitempty"`
	// reward_amount is the amount of reward.
	RewardAmount *RewardAmount `protobuf:"bytes,8,opt,name=reward_amount,json=rewardAmount,proto3" json:"reward_amount,omitempty"`
	// pool is the tracker of pool funds of the campaign.
	Pool Pool `protobuf:"bytes,9,opt,name=pool,proto3" json:"pool"`
	// validations required on the campaign reward distribution.
	// It is a stringified base64 encoded json.
	Validations string `protobuf:"bytes,10,opt,name=validations,proto3" json:"validations,omitempty"`
	// meta is the metadata of the campaign.
	// It is a stringified base64 encoded json.
	Meta string `protobuf:"bytes,11,opt,name=meta,proto3" json:"meta,omitempty"`
}

func (m *CreateCampaignPayload) Reset()         { *m = CreateCampaignPayload{} }
func (m *CreateCampaignPayload) String() string { return proto.CompactTextString(m) }
func (*CreateCampaignPayload) ProtoMessage()    {}
func (*CreateCampaignPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d710bc1249ca8ae, []int{0}
}
func (m *CreateCampaignPayload) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CreateCampaignPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CreateCampaignPayload.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CreateCampaignPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateCampaignPayload.Merge(m, src)
}
func (m *CreateCampaignPayload) XXX_Size() int {
	return m.Size()
}
func (m *CreateCampaignPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateCampaignPayload.DiscardUnknown(m)
}

var xxx_messageInfo_CreateCampaignPayload proto.InternalMessageInfo

func (m *CreateCampaignPayload) GetPromoter() string {
	if m != nil {
		return m.Promoter
	}
	return ""
}

func (m *CreateCampaignPayload) GetStartTs() uint64 {
	if m != nil {
		return m.StartTs
	}
	return 0
}

func (m *CreateCampaignPayload) GetEndTs() uint64 {
	if m != nil {
		return m.EndTs
	}
	return 0
}

func (m *CreateCampaignPayload) GetRewardType() RewardType {
	if m != nil {
		return m.RewardType
	}
	return RewardType_REWARD_TYPE_UNSPECIFIED
}

func (m *CreateCampaignPayload) GetRewardAmountType() RewardAmountType {
	if m != nil {
		return m.RewardAmountType
	}
	return RewardAmountType_REWARD_AMOUNT_TYPE_UNSPECIFIED
}

func (m *CreateCampaignPayload) GetRewardAmount() *RewardAmount {
	if m != nil {
		return m.RewardAmount
	}
	return nil
}

func (m *CreateCampaignPayload) GetPool() Pool {
	if m != nil {
		return m.Pool
	}
	return Pool{}
}

func (m *CreateCampaignPayload) GetValidations() string {
	if m != nil {
		return m.Validations
	}
	return ""
}

func (m *CreateCampaignPayload) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

// UpdateCampaignPayload is the type for campaign update payload.
type UpdateCampaignPayload struct {
	// end_ts is the end timestamp of the campaign.
	EndTs uint64 `protobuf:"varint,3,opt,name=end_ts,json=endTs,proto3" json:"end_ts,omitempty"`
	// funds is the additional amount added to the campaign.
	Funds string `protobuf:"bytes,4,opt,name=funds,proto3" json:"funds,omitempty"`
}

func (m *UpdateCampaignPayload) Reset()         { *m = UpdateCampaignPayload{} }
func (m *UpdateCampaignPayload) String() string { return proto.CompactTextString(m) }
func (*UpdateCampaignPayload) ProtoMessage()    {}
func (*UpdateCampaignPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d710bc1249ca8ae, []int{1}
}
func (m *UpdateCampaignPayload) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateCampaignPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateCampaignPayload.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpdateCampaignPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateCampaignPayload.Merge(m, src)
}
func (m *UpdateCampaignPayload) XXX_Size() int {
	return m.Size()
}
func (m *UpdateCampaignPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateCampaignPayload.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateCampaignPayload proto.InternalMessageInfo

func (m *UpdateCampaignPayload) GetEndTs() uint64 {
	if m != nil {
		return m.EndTs
	}
	return 0
}

func (m *UpdateCampaignPayload) GetFunds() string {
	if m != nil {
		return m.Funds
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateCampaignPayload)(nil), "sgenetwork.sge.reward.CreateCampaignPayload")
	proto.RegisterType((*UpdateCampaignPayload)(nil), "sgenetwork.sge.reward.UpdateCampaignPayload")
}

func init() { proto.RegisterFile("sge/reward/ticket.proto", fileDescriptor_5d710bc1249ca8ae) }

var fileDescriptor_5d710bc1249ca8ae = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x8e, 0xd3, 0x30,
	0x14, 0x86, 0x63, 0x26, 0x33, 0xd3, 0x3a, 0x80, 0x90, 0x35, 0x11, 0x69, 0x91, 0x42, 0x28, 0x0b,
	0xc2, 0x82, 0x44, 0x2a, 0xe2, 0x00, 0xb4, 0x2c, 0x58, 0x56, 0x51, 0xbb, 0x61, 0x53, 0xb9, 0xcd,
	0xc3, 0x44, 0x4d, 0x62, 0xcb, 0x76, 0x29, 0xb9, 0x05, 0xb7, 0xe1, 0x0a, 0x5d, 0x76, 0xc9, 0x0a,
	0xa1, 0xf6, 0x22, 0xa8, 0x4e, 0x54, 0x05, 0xd4, 0xa2, 0x59, 0xf9, 0x7f, 0xfe, 0x7f, 0x7d, 0xb2,
	0xdf, 0x7b, 0xf8, 0xa9, 0x62, 0x10, 0x4b, 0xd8, 0x50, 0x99, 0xc6, 0x3a, 0x5b, 0xae, 0x40, 0x47,
	0x42, 0x72, 0xcd, 0x89, 0xab, 0x18, 0x94, 0xa0, 0x37, 0x5c, 0xae, 0x22, 0xc5, 0x20, 0xaa, 0x33,
	0xfd, 0x3b, 0xc6, 0x19, 0x37, 0x89, 0xf8, 0xa8, 0xea, 0x70, 0xbf, 0x4d, 0xa9, 0x8f, 0xc6, 0xe8,
	0xb5, 0x8c, 0x25, 0x2d, 0x04, 0xcd, 0x58, 0x59, 0x5b, 0x83, 0x1f, 0x57, 0xd8, 0x1d, 0x4b, 0xa0,
	0x1a, 0xc6, 0x8d, 0x31, 0xa1, 0x55, 0xce, 0x69, 0x4a, 0xfa, 0xb8, 0x23, 0x24, 0x2f, 0xb8, 0x06,
	0xe9, 0xa1, 0x00, 0x85, 0xdd, 0xe4, 0x54, 0x93, 0x1e, 0xee, 0x28, 0x4d, 0xa5, 0x9e, 0x6b, 0xe5,
	0x3d, 0x08, 0x50, 0x68, 0x27, 0xb7, 0xa6, 0x9e, 0x2a, 0xe2, 0xe2, 0x1b, 0x28, 0xd3, 0xa3, 0x71,
	0x65, 0x8c, 0x6b, 0x28, 0xd3, 0xa9, 0x22, 0x23, 0xec, 0xd4, 0x0f, 0x98, 0xeb, 0x4a, 0x80, 0x67,
	0x07, 0x28, 0x7c, 0x3c, 0x7c, 0x11, 0x9d, 0xfd, 0x5e, 0x94, 0x98, 0x63, 0x5a, 0x09, 0x48, 0xb0,
	0x3c, 0x69, 0x32, 0xc3, 0xa4, 0x61, 0xd0, 0x82, 0xaf, 0x4b, 0x5d, 0xa3, 0x6e, 0x0d, 0xea, 0xd5,
	0x7f, 0x51, 0xef, 0x4d, 0xde, 0x00, 0x9f, 0xc8, 0x7f, 0x6e, 0xc8, 0x47, 0xfc, 0xe8, 0x2f, 0xac,
	0xd7, 0x09, 0x50, 0xe8, 0x0c, 0x5f, 0xde, 0x83, 0x98, 0x3c, 0x6c, 0xd3, 0xc8, 0x3b, 0x6c, 0x0b,
	0xce, 0x73, 0xaf, 0x6b, 0x00, 0xcf, 0x2e, 0x00, 0x26, 0x9c, 0xe7, 0x23, 0x7b, 0xfb, 0xeb, 0xb9,
	0x95, 0x98, 0x38, 0x09, 0xb0, 0xf3, 0x95, 0xe6, 0x59, 0x4a, 0x75, 0xc6, 0x4b, 0xe5, 0x61, 0xd3,
	0xec, 0xf6, 0x15, 0x21, 0xd8, 0x2e, 0x40, 0x53, 0xcf, 0x31, 0x96, 0xd1, 0x83, 0x0f, 0xd8, 0x9d,
	0x89, 0xf4, 0xcc, 0xe0, 0x2e, 0x4c, 0xe0, 0x0e, 0x5f, 0x7f, 0x5e, 0x97, 0xa9, 0x32, 0xbd, 0xef,
	0x26, 0x75, 0x31, 0x1a, 0x6f, 0xf7, 0x3e, 0xda, 0xed, 0x7d, 0xf4, 0x7b, 0xef, 0xa3, 0xef, 0x07,
	0xdf, 0xda, 0x1d, 0x7c, 0xeb, 0xe7, 0xc1, 0xb7, 0x3e, 0xbd, 0x66, 0x99, 0xfe, 0xb2, 0x5e, 0x44,
	0x4b, 0x5e, 0xc4, 0x8a, 0xc1, 0x9b, 0xe6, 0x27, 0x47, 0x1d, 0x7f, 0x3b, 0x2d, 0x6b, 0x25, 0x40,
	0x2d, 0x6e, 0xcc, 0x2e, 0xbd, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0xed, 0xca, 0xef, 0xb5, 0xc7,
	0x02, 0x00, 0x00,
}

func (m *CreateCampaignPayload) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CreateCampaignPayload) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CreateCampaignPayload) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Meta) > 0 {
		i -= len(m.Meta)
		copy(dAtA[i:], m.Meta)
		i = encodeVarintTicket(dAtA, i, uint64(len(m.Meta)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.Validations) > 0 {
		i -= len(m.Validations)
		copy(dAtA[i:], m.Validations)
		i = encodeVarintTicket(dAtA, i, uint64(len(m.Validations)))
		i--
		dAtA[i] = 0x52
	}
	{
		size, err := m.Pool.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTicket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if m.RewardAmount != nil {
		{
			size, err := m.RewardAmount.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTicket(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.RewardAmountType != 0 {
		i = encodeVarintTicket(dAtA, i, uint64(m.RewardAmountType))
		i--
		dAtA[i] = 0x38
	}
	if m.RewardType != 0 {
		i = encodeVarintTicket(dAtA, i, uint64(m.RewardType))
		i--
		dAtA[i] = 0x20
	}
	if m.EndTs != 0 {
		i = encodeVarintTicket(dAtA, i, uint64(m.EndTs))
		i--
		dAtA[i] = 0x18
	}
	if m.StartTs != 0 {
		i = encodeVarintTicket(dAtA, i, uint64(m.StartTs))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Promoter) > 0 {
		i -= len(m.Promoter)
		copy(dAtA[i:], m.Promoter)
		i = encodeVarintTicket(dAtA, i, uint64(len(m.Promoter)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UpdateCampaignPayload) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateCampaignPayload) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpdateCampaignPayload) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Funds) > 0 {
		i -= len(m.Funds)
		copy(dAtA[i:], m.Funds)
		i = encodeVarintTicket(dAtA, i, uint64(len(m.Funds)))
		i--
		dAtA[i] = 0x22
	}
	if m.EndTs != 0 {
		i = encodeVarintTicket(dAtA, i, uint64(m.EndTs))
		i--
		dAtA[i] = 0x18
	}
	return len(dAtA) - i, nil
}

func encodeVarintTicket(dAtA []byte, offset int, v uint64) int {
	offset -= sovTicket(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CreateCampaignPayload) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Promoter)
	if l > 0 {
		n += 1 + l + sovTicket(uint64(l))
	}
	if m.StartTs != 0 {
		n += 1 + sovTicket(uint64(m.StartTs))
	}
	if m.EndTs != 0 {
		n += 1 + sovTicket(uint64(m.EndTs))
	}
	if m.RewardType != 0 {
		n += 1 + sovTicket(uint64(m.RewardType))
	}
	if m.RewardAmountType != 0 {
		n += 1 + sovTicket(uint64(m.RewardAmountType))
	}
	if m.RewardAmount != nil {
		l = m.RewardAmount.Size()
		n += 1 + l + sovTicket(uint64(l))
	}
	l = m.Pool.Size()
	n += 1 + l + sovTicket(uint64(l))
	l = len(m.Validations)
	if l > 0 {
		n += 1 + l + sovTicket(uint64(l))
	}
	l = len(m.Meta)
	if l > 0 {
		n += 1 + l + sovTicket(uint64(l))
	}
	return n
}

func (m *UpdateCampaignPayload) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EndTs != 0 {
		n += 1 + sovTicket(uint64(m.EndTs))
	}
	l = len(m.Funds)
	if l > 0 {
		n += 1 + l + sovTicket(uint64(l))
	}
	return n
}

func sovTicket(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTicket(x uint64) (n int) {
	return sovTicket(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CreateCampaignPayload) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTicket
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
			return fmt.Errorf("proto: CreateCampaignPayload: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CreateCampaignPayload: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Promoter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
				return ErrInvalidLengthTicket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTicket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Promoter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTs", wireType)
			}
			m.StartTs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartTs |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTs", wireType)
			}
			m.EndTs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndTs |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardType", wireType)
			}
			m.RewardType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAmountType", wireType)
			}
			m.RewardAmountType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
				return ErrInvalidLengthTicket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTicket
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
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
				return ErrInvalidLengthTicket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTicket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Pool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validations", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
				return ErrInvalidLengthTicket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTicket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validations = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
				return ErrInvalidLengthTicket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTicket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Meta = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTicket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTicket
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
func (m *UpdateCampaignPayload) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTicket
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
			return fmt.Errorf("proto: UpdateCampaignPayload: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateCampaignPayload: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTs", wireType)
			}
			m.EndTs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndTs |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Funds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTicket
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
				return ErrInvalidLengthTicket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTicket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Funds = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTicket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTicket
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
func skipTicket(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTicket
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
					return 0, ErrIntOverflowTicket
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
					return 0, ErrIntOverflowTicket
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
				return 0, ErrInvalidLengthTicket
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTicket
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTicket
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTicket        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTicket          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTicket = fmt.Errorf("proto: unexpected end of group")
)
