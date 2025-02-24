// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sgenetwork/sge/market/stats.proto

package types

import (
	fmt "fmt"
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

// MarketStats holds statistics of the market
type MarketStats struct {
	// resolved_unsettled is the list of universal unique identifiers
	// of resolved markets that have unsettled bets
	ResolvedUnsettled []string `protobuf:"bytes,1,rep,name=resolved_unsettled,json=resolvedUnsettled,proto3" json:"resolved_unsettled,omitempty"`
}

func (m *MarketStats) Reset()         { *m = MarketStats{} }
func (m *MarketStats) String() string { return proto.CompactTextString(m) }
func (*MarketStats) ProtoMessage()    {}
func (*MarketStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_b89d96b591dbbbbd, []int{0}
}
func (m *MarketStats) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketStats.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketStats.Merge(m, src)
}
func (m *MarketStats) XXX_Size() int {
	return m.Size()
}
func (m *MarketStats) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketStats.DiscardUnknown(m)
}

var xxx_messageInfo_MarketStats proto.InternalMessageInfo

func (m *MarketStats) GetResolvedUnsettled() []string {
	if m != nil {
		return m.ResolvedUnsettled
	}
	return nil
}

func init() {
	proto.RegisterType((*MarketStats)(nil), "sgenetwork.sge.market.MarketStats")
}

func init() { proto.RegisterFile("sgenetwork/sge/market/stats.proto", fileDescriptor_b89d96b591dbbbbd) }

var fileDescriptor_b89d96b591dbbbbd = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2c, 0x4e, 0x4f, 0xcd,
	0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x2f, 0x4e, 0x4f, 0xd5, 0xcf, 0x4d, 0x2c, 0xca, 0x4e,
	0x2d, 0xd1, 0x2f, 0x2e, 0x49, 0x2c, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x45,
	0x28, 0xd1, 0x2b, 0x4e, 0x4f, 0xd5, 0x83, 0x28, 0x51, 0xb2, 0xe1, 0xe2, 0xf6, 0x05, 0xb3, 0x82,
	0x41, 0x6a, 0x85, 0x74, 0xb9, 0x84, 0x8a, 0x52, 0x8b, 0xf3, 0x73, 0xca, 0x52, 0x53, 0xe2, 0x4b,
	0xf3, 0x8a, 0x53, 0x4b, 0x4a, 0x72, 0x52, 0x53, 0x24, 0x18, 0x15, 0x98, 0x35, 0x38, 0x83, 0x04,
	0x61, 0x32, 0xa1, 0x30, 0x09, 0x27, 0xaf, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c,
	0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63,
	0x88, 0x32, 0x48, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0x05, 0xb9, 0x48, 0x17,
	0xd9, 0x75, 0x15, 0xfa, 0x39, 0xa9, 0xe9, 0x89, 0xc9, 0x95, 0x30, 0x67, 0x96, 0x54, 0x16, 0xa4,
	0x16, 0x27, 0xb1, 0x81, 0xdd, 0x69, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x89, 0x4e, 0x58,
	0xcc, 0x00, 0x00, 0x00,
}

func (m *MarketStats) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketStats) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketStats) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ResolvedUnsettled) > 0 {
		for iNdEx := len(m.ResolvedUnsettled) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ResolvedUnsettled[iNdEx])
			copy(dAtA[i:], m.ResolvedUnsettled[iNdEx])
			i = encodeVarintStats(dAtA, i, uint64(len(m.ResolvedUnsettled[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintStats(dAtA []byte, offset int, v uint64) int {
	offset -= sovStats(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MarketStats) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ResolvedUnsettled) > 0 {
		for _, s := range m.ResolvedUnsettled {
			l = len(s)
			n += 1 + l + sovStats(uint64(l))
		}
	}
	return n
}

func sovStats(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStats(x uint64) (n int) {
	return sovStats(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MarketStats) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStats
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
			return fmt.Errorf("proto: MarketStats: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketStats: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResolvedUnsettled", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
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
				return ErrInvalidLengthStats
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ResolvedUnsettled = append(m.ResolvedUnsettled, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStats(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStats
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
func skipStats(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStats
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
					return 0, ErrIntOverflowStats
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
					return 0, ErrIntOverflowStats
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
				return 0, ErrInvalidLengthStats
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStats
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStats
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStats        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStats          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStats = fmt.Errorf("proto: unexpected end of group")
)
