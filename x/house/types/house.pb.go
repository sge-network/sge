// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/house/house.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

// Params defines the parameters for the house module.
type Params struct {
	// minimum_deposit is the minum amount of deposit acceptable.
	MinimumDeposit uint64 `protobuf:"varint,1,opt,name=minimum_deposit,json=minimumDeposit,proto3" json:"minimum_deposit,omitempty" yaml:"minimum_deposit"`
	// house_participation_fee is the % of deposit to be paid for house participation by the user
	HouseParticipationFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=house_participation_fee,json=houseParticipationFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"house_participation_fee"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ac648d293874eae, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMinimumDeposit() uint64 {
	if m != nil {
		return m.MinimumDeposit
	}
	return 0
}

// Deposit represents the deposit against a sport event held by an account.
type Deposit struct {
	// depositor_address is the bech32-encoded address of the depositor.
	DepositorAddress string `protobuf:"bytes,1,opt,name=depositor_address,json=depositorAddress,proto3" json:"depositor_address,omitempty" yaml:"depositor_address"`
	// sport_event_uid is the uid of sport event against which deposit is being made.
	SportEventUid string `protobuf:"bytes,2,opt,name=sport_event_uid,json=sportEventUid,proto3" json:"sport_event_uid,omitempty" yaml:"sport_event_uid"`
	// amount is the amount being deposited.
	Amount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount" yaml:"amount"`
	// fee is deducted from amount at the point of deposit.
	Fee github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=fee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"fee" yaml:"fee"`
	// liquidity is the liquidity being provided to the house after fee deduction.
	Liquidity github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=liquidity,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"liquidity" yaml:"liquidity"`
	// id corresponding to the book participant
	ParticipantId uint64 `protobuf:"varint,6,opt,name=participant_id,json=participantId,proto3" json:"participant_id,omitempty" yaml:"participant_id"`
}

func (m *Deposit) Reset()      { *m = Deposit{} }
func (*Deposit) ProtoMessage() {}
func (*Deposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ac648d293874eae, []int{1}
}
func (m *Deposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Deposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Deposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Deposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deposit.Merge(m, src)
}
func (m *Deposit) XXX_Size() int {
	return m.Size()
}
func (m *Deposit) XXX_DiscardUnknown() {
	xxx_messageInfo_Deposit.DiscardUnknown(m)
}

var xxx_messageInfo_Deposit proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "sgenetwork.sge.house.Params")
	proto.RegisterType((*Deposit)(nil), "sgenetwork.sge.house.Deposit")
}

func init() { proto.RegisterFile("sge/house/house.proto", fileDescriptor_7ac648d293874eae) }

var fileDescriptor_7ac648d293874eae = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xc7, 0x6d, 0x52, 0x02, 0x39, 0x91, 0xb6, 0x58, 0x0d, 0x98, 0x08, 0xd9, 0xd5, 0x0d, 0x28,
	0x4b, 0x6d, 0x55, 0x6c, 0x11, 0x12, 0x60, 0x0a, 0x52, 0x16, 0x54, 0x59, 0x42, 0x48, 0x2c, 0xc6,
	0xf1, 0xbd, 0xb8, 0xa7, 0xd6, 0x3e, 0xe3, 0x3b, 0x17, 0xf2, 0x0d, 0x18, 0x19, 0x19, 0xf3, 0x4d,
	0x58, 0x3b, 0x76, 0x44, 0x0c, 0x16, 0x4a, 0x16, 0xe6, 0xf0, 0x05, 0x90, 0xef, 0x8e, 0x10, 0xda,
	0x29, 0x4b, 0x72, 0xfa, 0xbf, 0xff, 0xfd, 0xee, 0xff, 0xce, 0xf7, 0x50, 0x8f, 0xa7, 0xe0, 0x9f,
	0xb0, 0x8a, 0xeb, 0x5f, 0xaf, 0x28, 0x99, 0x60, 0xd6, 0x1e, 0x4f, 0x21, 0x07, 0xf1, 0x91, 0x95,
	0xa7, 0x1e, 0x4f, 0xc1, 0x93, 0xb5, 0xfe, 0x5e, 0xca, 0x52, 0x26, 0x0d, 0x7e, 0xb3, 0x52, 0xde,
	0xbe, 0x93, 0x30, 0x9e, 0x31, 0xee, 0x8f, 0x63, 0x0e, 0xfe, 0xf9, 0xe1, 0x18, 0x44, 0x7c, 0xe8,
	0x27, 0x8c, 0xe6, 0xaa, 0x8e, 0xbf, 0x99, 0xa8, 0x7d, 0x1c, 0x97, 0x71, 0xc6, 0xad, 0x17, 0x68,
	0x27, 0xa3, 0x39, 0xcd, 0xaa, 0x2c, 0x22, 0x50, 0x30, 0x4e, 0x85, 0x6d, 0xee, 0x9b, 0x83, 0xad,
	0xa0, 0xbf, 0xac, 0xdd, 0x7b, 0xd3, 0x38, 0x3b, 0x1b, 0xe2, 0x2b, 0x06, 0x1c, 0x6e, 0x6b, 0xe5,
	0x48, 0x09, 0xd6, 0x04, 0xdd, 0x97, 0x71, 0xa2, 0x22, 0x2e, 0x05, 0x4d, 0x68, 0x11, 0x0b, 0xca,
	0xf2, 0x68, 0x02, 0x60, 0xdf, 0xd8, 0x37, 0x07, 0x9d, 0xc0, 0xbb, 0xa8, 0x5d, 0xe3, 0x47, 0xed,
	0x3e, 0x4a, 0xa9, 0x38, 0xa9, 0xc6, 0x5e, 0xc2, 0x32, 0x5f, 0x67, 0x54, 0x7f, 0x07, 0x9c, 0x9c,
	0xfa, 0x62, 0x5a, 0x00, 0xf7, 0x8e, 0x20, 0x09, 0x7b, 0x12, 0x77, 0xbc, 0x4e, 0x7b, 0x05, 0x30,
	0xbc, 0xfd, 0x75, 0xe6, 0x1a, 0xbf, 0x66, 0xae, 0x89, 0x7f, 0xb7, 0xd0, 0xad, 0xbf, 0xa7, 0x8f,
	0xd0, 0x5d, 0x9d, 0x8c, 0x95, 0x51, 0x4c, 0x48, 0x09, 0x9c, 0xcb, 0x26, 0x3a, 0xc1, 0xc3, 0x65,
	0xed, 0xda, 0xaa, 0x89, 0x6b, 0x16, 0x1c, 0xee, 0xae, 0xb4, 0xe7, 0x4a, 0xb2, 0x02, 0xb4, 0xc3,
	0x0b, 0x56, 0x8a, 0x08, 0xce, 0x21, 0x17, 0x51, 0x45, 0x89, 0x6e, 0x60, 0xed, 0x36, 0xae, 0x18,
	0x70, 0xd8, 0x95, 0xca, 0xcb, 0x46, 0x78, 0x43, 0x89, 0xf5, 0x16, 0xb5, 0xe3, 0x8c, 0x55, 0xb9,
	0xb0, 0x5b, 0x72, 0xeb, 0xd3, 0x0d, 0x7a, 0x1f, 0xe5, 0x62, 0x59, 0xbb, 0x5d, 0x75, 0x90, 0xa2,
	0xe0, 0x50, 0xe3, 0xac, 0xd7, 0xa8, 0xd5, 0xdc, 0xe8, 0x96, 0xa4, 0x3e, 0xd9, 0x98, 0x8a, 0x14,
	0x75, 0x02, 0x80, 0xc3, 0x06, 0x64, 0xbd, 0x47, 0x9d, 0x33, 0xfa, 0xa1, 0xa2, 0x84, 0x8a, 0xa9,
	0x7d, 0x53, 0x52, 0x83, 0x8d, 0xa9, 0xbb, 0x8a, 0xba, 0x02, 0xe1, 0xf0, 0x1f, 0xd4, 0x7a, 0x86,
	0xb6, 0x57, 0x2f, 0x22, 0x17, 0x11, 0x25, 0x76, 0x5b, 0xbe, 0xad, 0x07, 0xcb, 0xda, 0xed, 0xa9,
	0x8d, 0xff, 0xd7, 0x71, 0xd8, 0x5d, 0x13, 0x46, 0x64, 0x78, 0xe7, 0xf3, 0xcc, 0x35, 0xf4, 0x57,
	0x37, 0x82, 0xe0, 0x62, 0xee, 0x98, 0x97, 0x73, 0xc7, 0xfc, 0x39, 0x77, 0xcc, 0x2f, 0x0b, 0xc7,
	0xb8, 0x5c, 0x38, 0xc6, 0xf7, 0x85, 0x63, 0xbc, 0x1b, 0xac, 0x05, 0xe6, 0x29, 0x1c, 0xe8, 0x49,
	0x69, 0xd6, 0xfe, 0x27, 0x3d, 0x4d, 0x32, 0xf6, 0xb8, 0x2d, 0x47, 0xe0, 0xf1, 0x9f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xa9, 0x13, 0xa9, 0xf5, 0x67, 0x03, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.MinimumDeposit != that1.MinimumDeposit {
		return false
	}
	if !this.HouseParticipationFee.Equal(that1.HouseParticipationFee) {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.HouseParticipationFee.Size()
		i -= size
		if _, err := m.HouseParticipationFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintHouse(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.MinimumDeposit != 0 {
		i = encodeVarintHouse(dAtA, i, uint64(m.MinimumDeposit))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Deposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Deposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Deposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ParticipantId != 0 {
		i = encodeVarintHouse(dAtA, i, uint64(m.ParticipantId))
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.Liquidity.Size()
		i -= size
		if _, err := m.Liquidity.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintHouse(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.Fee.Size()
		i -= size
		if _, err := m.Fee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintHouse(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintHouse(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.SportEventUid) > 0 {
		i -= len(m.SportEventUid)
		copy(dAtA[i:], m.SportEventUid)
		i = encodeVarintHouse(dAtA, i, uint64(len(m.SportEventUid)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DepositorAddress) > 0 {
		i -= len(m.DepositorAddress)
		copy(dAtA[i:], m.DepositorAddress)
		i = encodeVarintHouse(dAtA, i, uint64(len(m.DepositorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintHouse(dAtA []byte, offset int, v uint64) int {
	offset -= sovHouse(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MinimumDeposit != 0 {
		n += 1 + sovHouse(uint64(m.MinimumDeposit))
	}
	l = m.HouseParticipationFee.Size()
	n += 1 + l + sovHouse(uint64(l))
	return n
}

func (m *Deposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DepositorAddress)
	if l > 0 {
		n += 1 + l + sovHouse(uint64(l))
	}
	l = len(m.SportEventUid)
	if l > 0 {
		n += 1 + l + sovHouse(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovHouse(uint64(l))
	l = m.Fee.Size()
	n += 1 + l + sovHouse(uint64(l))
	l = m.Liquidity.Size()
	n += 1 + l + sovHouse(uint64(l))
	if m.ParticipantId != 0 {
		n += 1 + sovHouse(uint64(m.ParticipantId))
	}
	return n
}

func sovHouse(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozHouse(x uint64) (n int) {
	return sovHouse(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHouse
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumDeposit", wireType)
			}
			m.MinimumDeposit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinimumDeposit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HouseParticipationFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
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
				return ErrInvalidLengthHouse
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHouse
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.HouseParticipationFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHouse(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHouse
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
func (m *Deposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHouse
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
			return fmt.Errorf("proto: Deposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Deposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
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
				return ErrInvalidLengthHouse
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHouse
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SportEventUid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
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
				return ErrInvalidLengthHouse
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHouse
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SportEventUid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
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
				return ErrInvalidLengthHouse
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHouse
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
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
				return ErrInvalidLengthHouse
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHouse
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Fee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
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
				return ErrInvalidLengthHouse
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHouse
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liquidity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParticipantId", wireType)
			}
			m.ParticipantId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHouse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ParticipantId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHouse(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHouse
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
func skipHouse(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHouse
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
					return 0, ErrIntOverflowHouse
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
					return 0, ErrIntOverflowHouse
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
				return 0, ErrInvalidLengthHouse
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupHouse
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthHouse
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthHouse        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHouse          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupHouse = fmt.Errorf("proto: unexpected end of group")
)