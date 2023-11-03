// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/house/deposit.proto

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

// Deposit represents the deposit against a market held by an account.
type Deposit struct {
	// creator is the bech32-encoded address of the depositor.
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty" yaml:"creator"`
	// creator is the bech32-encoded address of the depositor.
	DepositorAddress string `protobuf:"bytes,2,opt,name=depositor_address,json=depositorAddress,proto3" json:"depositor_address,omitempty" yaml:"depositor_address"`
	// market_uid is the uid of market/order book against which deposit is being
	// made.
	MarketUID string `protobuf:"bytes,3,opt,name=market_uid,proto3" json:"market_uid"`
	// participation_index is the index corresponding to the order book
	// participation
	ParticipationIndex uint64 `protobuf:"varint,4,opt,name=participation_index,json=participationIndex,proto3" json:"participation_index,omitempty" yaml:"participation_index"`
	// amount is the amount being deposited on an order book to be a house
	Amount cosmossdk_io_math.Int `protobuf:"bytes,5,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount" yaml:"amount"`
	// withdrawal_count is the total count of the withdrawals from an order book
	WithdrawalCount uint64 `protobuf:"varint,6,opt,name=withdrawal_count,json=withdrawalCount,proto3" json:"withdrawal_count,omitempty" yaml:"withdrawals"`
	// total_withdrawal_amount is the total amount withdrawn from the liquidity
	// provided
	TotalWithdrawalAmount cosmossdk_io_math.Int `protobuf:"bytes,7,opt,name=total_withdrawal_amount,json=totalWithdrawalAmount,proto3,customtype=cosmossdk.io/math.Int" json:"total_withdrawal_amount" yaml:"total_withdrawal_amount"`
}

func (m *Deposit) Reset()      { *m = Deposit{} }
func (*Deposit) ProtoMessage() {}
func (*Deposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6f2840908fc45a1, []int{0}
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
	proto.RegisterType((*Deposit)(nil), "sgenetwork.sge.house.Deposit")
}

func init() { proto.RegisterFile("sge/house/deposit.proto", fileDescriptor_c6f2840908fc45a1) }

var fileDescriptor_c6f2840908fc45a1 = []byte{
	// 441 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x3f, 0x6f, 0xd3, 0x40,
	0x18, 0xc6, 0x6d, 0x08, 0x89, 0x7a, 0xe2, 0x4f, 0x39, 0x5a, 0x6a, 0x15, 0xe4, 0xab, 0x6e, 0xca,
	0x00, 0xf6, 0xc0, 0x56, 0x06, 0x94, 0x50, 0x21, 0x65, 0x40, 0x48, 0x96, 0x10, 0x12, 0x8b, 0x75,
	0xb5, 0x4f, 0xce, 0x29, 0xb1, 0x5f, 0xeb, 0xee, 0xa2, 0xb4, 0xdf, 0x80, 0x91, 0x91, 0x31, 0xe2,
	0xd3, 0x74, 0xec, 0x88, 0x18, 0x4e, 0xc8, 0x59, 0x10, 0xa3, 0x3f, 0x01, 0xca, 0x9d, 0x69, 0x83,
	0x00, 0xb1, 0xdd, 0x3d, 0xef, 0xef, 0x79, 0xde, 0x67, 0x78, 0xd1, 0x81, 0x2a, 0x78, 0x3c, 0x85,
	0x85, 0xe2, 0x71, 0xce, 0x6b, 0x50, 0x42, 0x47, 0xb5, 0x04, 0x0d, 0x78, 0x4f, 0x15, 0xbc, 0xe2,
	0x7a, 0x09, 0x72, 0x16, 0xa9, 0x82, 0x47, 0x96, 0x39, 0xdc, 0x2b, 0xa0, 0x00, 0x0b, 0xc4, 0x9b,
	0x97, 0x63, 0xe9, 0xe7, 0x1e, 0x1a, 0x9c, 0x38, 0x37, 0x7e, 0x82, 0x06, 0x99, 0xe4, 0x4c, 0x83,
	0x0c, 0xfc, 0x23, 0x7f, 0xb8, 0x33, 0xc6, 0xad, 0x21, 0x77, 0xcf, 0x59, 0x39, 0x3f, 0xa6, 0xdd,
	0x80, 0x26, 0xbf, 0x10, 0x3c, 0x41, 0xf7, 0xbb, 0xb5, 0x20, 0x53, 0x96, 0xe7, 0x92, 0x2b, 0x15,
	0xdc, 0xb0, 0xbe, 0xc7, 0xad, 0x21, 0x81, 0xf3, 0xfd, 0x81, 0xd0, 0x64, 0xf7, 0x4a, 0x1b, 0x39,
	0x09, 0x3f, 0x47, 0xa8, 0x64, 0x72, 0xc6, 0x75, 0xba, 0x10, 0x79, 0x70, 0xd3, 0x66, 0x3c, 0x6a,
	0x0c, 0xd9, 0x79, 0x6d, 0xd5, 0xb7, 0x93, 0x93, 0x1f, 0x86, 0x6c, 0x21, 0xc9, 0xd6, 0x1b, 0xbf,
	0x41, 0x0f, 0x6a, 0x26, 0xb5, 0xc8, 0x44, 0xcd, 0xb4, 0x80, 0x2a, 0x15, 0x55, 0xce, 0xcf, 0x82,
	0xde, 0x91, 0x3f, 0xec, 0x8d, 0xc3, 0xd6, 0x90, 0x43, 0xd7, 0xe4, 0x2f, 0x10, 0x4d, 0xf0, 0x6f,
	0xea, 0x64, 0x23, 0xe2, 0x57, 0xa8, 0xcf, 0x4a, 0x58, 0x54, 0x3a, 0xb8, 0x65, 0x9b, 0x44, 0x17,
	0x86, 0x78, 0x5f, 0x0d, 0xd9, 0xcf, 0x40, 0x95, 0xa0, 0x54, 0x3e, 0x8b, 0x04, 0xc4, 0x25, 0xd3,
	0xd3, 0x68, 0x52, 0xe9, 0xd6, 0x90, 0x3b, 0x6e, 0x81, 0x33, 0xd1, 0xa4, 0x73, 0xe3, 0x11, 0xda,
	0x5d, 0x0a, 0x3d, 0xcd, 0x25, 0x5b, 0xb2, 0x79, 0x9a, 0xd9, 0xc4, 0xbe, 0x6d, 0xf5, 0xb0, 0x35,
	0x04, 0x3b, 0xd3, 0x35, 0xa1, 0x68, 0x72, 0xef, 0xfa, 0xf7, 0xd2, 0x46, 0x2c, 0xd1, 0x81, 0x06,
	0xcd, 0xe6, 0xe9, 0x56, 0x50, 0xd7, 0x6d, 0x60, 0xbb, 0xbd, 0xf8, 0x5f, 0xb7, 0xd0, 0xad, 0xf9,
	0x47, 0x0a, 0x4d, 0xf6, 0xed, 0xe4, 0xdd, 0xd5, 0x60, 0x64, 0xf5, 0xe3, 0xdb, 0x1f, 0x56, 0xc4,
	0xfb, 0xb4, 0x22, 0xde, 0xf7, 0x15, 0xf1, 0xc6, 0xe3, 0x8b, 0x26, 0xf4, 0x2f, 0x9b, 0xd0, 0xff,
	0xd6, 0x84, 0xfe, 0xc7, 0x75, 0xe8, 0x5d, 0xae, 0x43, 0xef, 0xcb, 0x3a, 0xf4, 0xde, 0x0f, 0x0b,
	0xa1, 0xa7, 0x8b, 0xd3, 0x28, 0x83, 0x32, 0x56, 0x05, 0x7f, 0xda, 0x9d, 0xdd, 0xe6, 0x1d, 0x9f,
	0x75, 0xc7, 0xa9, 0xcf, 0x6b, 0xae, 0x4e, 0xfb, 0xf6, 0xde, 0x9e, 0xfd, 0x0c, 0x00, 0x00, 0xff,
	0xff, 0xce, 0xd4, 0x22, 0x0c, 0xb6, 0x02, 0x00, 0x00,
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
	{
		size := m.TotalWithdrawalAmount.Size()
		i -= size
		if _, err := m.TotalWithdrawalAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDeposit(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.WithdrawalCount != 0 {
		i = encodeVarintDeposit(dAtA, i, uint64(m.WithdrawalCount))
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDeposit(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.ParticipationIndex != 0 {
		i = encodeVarintDeposit(dAtA, i, uint64(m.ParticipationIndex))
		i--
		dAtA[i] = 0x20
	}
	if len(m.MarketUID) > 0 {
		i -= len(m.MarketUID)
		copy(dAtA[i:], m.MarketUID)
		i = encodeVarintDeposit(dAtA, i, uint64(len(m.MarketUID)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DepositorAddress) > 0 {
		i -= len(m.DepositorAddress)
		copy(dAtA[i:], m.DepositorAddress)
		i = encodeVarintDeposit(dAtA, i, uint64(len(m.DepositorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintDeposit(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDeposit(dAtA []byte, offset int, v uint64) int {
	offset -= sovDeposit(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Deposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovDeposit(uint64(l))
	}
	l = len(m.DepositorAddress)
	if l > 0 {
		n += 1 + l + sovDeposit(uint64(l))
	}
	l = len(m.MarketUID)
	if l > 0 {
		n += 1 + l + sovDeposit(uint64(l))
	}
	if m.ParticipationIndex != 0 {
		n += 1 + sovDeposit(uint64(m.ParticipationIndex))
	}
	l = m.Amount.Size()
	n += 1 + l + sovDeposit(uint64(l))
	if m.WithdrawalCount != 0 {
		n += 1 + sovDeposit(uint64(m.WithdrawalCount))
	}
	l = m.TotalWithdrawalAmount.Size()
	n += 1 + l + sovDeposit(uint64(l))
	return n
}

func sovDeposit(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDeposit(x uint64) (n int) {
	return sovDeposit(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Deposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDeposit
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
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
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
				return ErrInvalidLengthDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
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
				return ErrInvalidLengthDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarketUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
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
				return ErrInvalidLengthDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MarketUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParticipationIndex", wireType)
			}
			m.ParticipationIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ParticipationIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
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
				return ErrInvalidLengthDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawalCount", wireType)
			}
			m.WithdrawalCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.WithdrawalCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalWithdrawalAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeposit
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
				return ErrInvalidLengthDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalWithdrawalAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDeposit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDeposit
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
func skipDeposit(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDeposit
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
					return 0, ErrIntOverflowDeposit
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
					return 0, ErrIntOverflowDeposit
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
				return 0, ErrInvalidLengthDeposit
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDeposit
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDeposit
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDeposit        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDeposit          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDeposit = fmt.Errorf("proto: unexpected end of group")
)
