// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sgenetwork/sge/bet/params.proto

package types

import (
	fmt "fmt"
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

// Params defines the parameters for the module.
type Params struct {
	// batch_settlement_count is the batch settlement bet count.
	BatchSettlementCount uint32 `protobuf:"varint,1,opt,name=batch_settlement_count,json=batchSettlementCount,proto3" json:"batch_settlement_count,omitempty"`
	// max_bet_by_uid_query_count is the maximum bet by uid query items count.
	MaxBetByUidQueryCount uint32 `protobuf:"varint,2,opt,name=max_bet_by_uid_query_count,json=maxBetByUidQueryCount,proto3" json:"max_bet_by_uid_query_count,omitempty"`
	// constraints is the bet constraints.
	Constraints Constraints `protobuf:"bytes,3,opt,name=constraints,proto3" json:"constraints" yaml:"constraints"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_f6bea002c4bbd0b9, []int{0}
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

func (m *Params) GetBatchSettlementCount() uint32 {
	if m != nil {
		return m.BatchSettlementCount
	}
	return 0
}

func (m *Params) GetMaxBetByUidQueryCount() uint32 {
	if m != nil {
		return m.MaxBetByUidQueryCount
	}
	return 0
}

func (m *Params) GetConstraints() Constraints {
	if m != nil {
		return m.Constraints
	}
	return Constraints{}
}

func init() {
	proto.RegisterType((*Params)(nil), "sgenetwork.sge.bet.Params")
}

func init() { proto.RegisterFile("sgenetwork/sge/bet/params.proto", fileDescriptor_f6bea002c4bbd0b9) }

var fileDescriptor_f6bea002c4bbd0b9 = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xb1, 0x4e, 0xeb, 0x30,
	0x18, 0x85, 0xe3, 0x7b, 0x51, 0x87, 0x54, 0x2c, 0x51, 0x41, 0x55, 0x86, 0xa4, 0xaa, 0x18, 0xba,
	0x34, 0x91, 0x80, 0x85, 0x8e, 0xe9, 0xc0, 0x0a, 0x45, 0x2c, 0x48, 0xc8, 0xb2, 0xd3, 0x5f, 0x6e,
	0x44, 0x1d, 0x97, 0xf8, 0x8f, 0xa8, 0xdf, 0x82, 0x91, 0x91, 0xc7, 0xe9, 0xd8, 0x91, 0xa9, 0x42,
	0xcd, 0x1b, 0xf0, 0x04, 0x28, 0x6e, 0x45, 0x23, 0xc1, 0x66, 0xe9, 0x3b, 0xe7, 0xd8, 0xfe, 0xdc,
	0x50, 0x0b, 0xc8, 0x01, 0x5f, 0x54, 0xf1, 0x14, 0x6b, 0x01, 0x31, 0x07, 0x8c, 0x17, 0xac, 0x60,
	0x52, 0x47, 0x8b, 0x42, 0xa1, 0xf2, 0xbc, 0x43, 0x20, 0xd2, 0x02, 0x22, 0x0e, 0xe8, 0x77, 0x84,
	0x12, 0xca, 0xe2, 0xb8, 0x3e, 0xed, 0x92, 0xfe, 0xd9, 0x1f, 0x53, 0xa9, 0xca, 0x35, 0x16, 0x2c,
	0xcb, 0x71, 0xbf, 0xd7, 0xdf, 0x10, 0xb7, 0x75, 0x63, 0x2f, 0xf0, 0x2e, 0xdd, 0x53, 0xce, 0x30,
	0x9d, 0x51, 0x0d, 0x88, 0x73, 0x90, 0x90, 0x23, 0x4d, 0x55, 0x99, 0x63, 0x97, 0xf4, 0xc8, 0xe0,
	0x78, 0xd2, 0xb1, 0xf4, 0xee, 0x07, 0x8e, 0x6b, 0xe6, 0x5d, 0xb9, 0xbe, 0x64, 0x4b, 0xca, 0x01,
	0x29, 0x37, 0xb4, 0xcc, 0xa6, 0xf4, 0xb9, 0x84, 0xc2, 0xec, 0x9b, 0xff, 0x6c, 0xf3, 0x44, 0xb2,
	0x65, 0x02, 0x98, 0x98, 0xfb, 0x6c, 0x7a, 0x5b, 0xd3, 0x5d, 0xf5, 0xd1, 0x6d, 0x37, 0x1e, 0xd4,
	0xfd, 0xdf, 0x23, 0x83, 0xf6, 0x79, 0x18, 0xfd, 0xfe, 0x61, 0x34, 0x3e, 0xc4, 0x12, 0x7f, 0xb5,
	0x09, 0x9d, 0xaf, 0x4d, 0xe8, 0x19, 0x26, 0xe7, 0xa3, 0x7e, 0x63, 0xa1, 0x3f, 0x69, 0xee, 0x8d,
	0x8e, 0xde, 0xde, 0x43, 0x27, 0xb9, 0x5e, 0x6d, 0x03, 0xb2, 0xde, 0x06, 0xe4, 0x73, 0x1b, 0x90,
	0xd7, 0x2a, 0x70, 0xd6, 0x55, 0xe0, 0x7c, 0x54, 0x81, 0xf3, 0x30, 0x14, 0x19, 0xce, 0x4a, 0x1e,
	0xa5, 0x4a, 0xd6, 0x82, 0x86, 0x4d, 0x59, 0xcb, 0x78, 0x0e, 0x82, 0xa5, 0xc6, 0x5a, 0x43, 0xb3,
	0x00, 0xcd, 0x5b, 0x56, 0xd8, 0xc5, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc9, 0x83, 0xc5, 0x76,
	0xa3, 0x01, 0x00, 0x00,
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
		size, err := m.Constraints.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.MaxBetByUidQueryCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxBetByUidQueryCount))
		i--
		dAtA[i] = 0x10
	}
	if m.BatchSettlementCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.BatchSettlementCount))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
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
	if m.BatchSettlementCount != 0 {
		n += 1 + sovParams(uint64(m.BatchSettlementCount))
	}
	if m.MaxBetByUidQueryCount != 0 {
		n += 1 + sovParams(uint64(m.MaxBetByUidQueryCount))
	}
	l = m.Constraints.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
				return fmt.Errorf("proto: wrong wireType = %d for field BatchSettlementCount", wireType)
			}
			m.BatchSettlementCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BatchSettlementCount |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBetByUidQueryCount", wireType)
			}
			m.MaxBetByUidQueryCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxBetByUidQueryCount |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Constraints", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Constraints.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
