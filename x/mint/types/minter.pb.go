// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sgenetwork/sge/mint/minter.proto

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

// Minter represents the minting state.
type Minter struct {
	// inflation is the current annual inflation rate.
	Inflation cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=inflation,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"inflation"`
	// phase_step is the index of phases slice + 1.
	PhaseStep int32 `protobuf:"varint,2,opt,name=phase_step,json=phaseStep,proto3" json:"phase_step,omitempty"`
	// phase_provisions is the current phase expected provisions.
	PhaseProvisions cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=phase_provisions,json=phaseProvisions,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"phase_provisions" yaml:"phase_provisions"`
	// truncated_tokens holds current truncated tokens because of Dec to Int
	// conversion in the minting.
	TruncatedTokens cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=truncated_tokens,json=truncatedTokens,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"truncated_tokens"`
}

func (m *Minter) Reset()         { *m = Minter{} }
func (m *Minter) String() string { return proto.CompactTextString(m) }
func (*Minter) ProtoMessage()    {}
func (*Minter) Descriptor() ([]byte, []int) {
	return fileDescriptor_a969bd1a18aad64f, []int{0}
}
func (m *Minter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Minter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Minter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Minter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Minter.Merge(m, src)
}
func (m *Minter) XXX_Size() int {
	return m.Size()
}
func (m *Minter) XXX_DiscardUnknown() {
	xxx_messageInfo_Minter.DiscardUnknown(m)
}

var xxx_messageInfo_Minter proto.InternalMessageInfo

func (m *Minter) GetPhaseStep() int32 {
	if m != nil {
		return m.PhaseStep
	}
	return 0
}

func init() {
	proto.RegisterType((*Minter)(nil), "sgenetwork.sge.mint.Minter")
}

func init() { proto.RegisterFile("sgenetwork/sge/mint/minter.proto", fileDescriptor_a969bd1a18aad64f) }

var fileDescriptor_a969bd1a18aad64f = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0xc1, 0x4e, 0xf2, 0x40,
	0x14, 0x85, 0x5b, 0xfe, 0x5f, 0x12, 0x66, 0x03, 0xa9, 0x26, 0x56, 0x8c, 0x85, 0xb0, 0x91, 0x0d,
	0x6d, 0x8c, 0x3b, 0x77, 0x12, 0x96, 0x18, 0x0d, 0xba, 0x32, 0x26, 0xa4, 0x94, 0xeb, 0x30, 0x81,
	0xce, 0x6d, 0x3a, 0x17, 0x91, 0xb7, 0xf0, 0x61, 0x7c, 0x08, 0x76, 0x12, 0x57, 0xc6, 0x05, 0x31,
	0xf0, 0x06, 0x3e, 0x81, 0xe9, 0x8c, 0x01, 0xc3, 0xca, 0xb8, 0x69, 0x3a, 0xe7, 0x9c, 0x7c, 0x67,
	0x26, 0x87, 0x55, 0x15, 0x07, 0x09, 0x34, 0xc1, 0x74, 0x18, 0x28, 0x0e, 0x41, 0x2c, 0x24, 0xe9,
	0x0f, 0xa4, 0x7e, 0x92, 0x22, 0xa1, 0xb3, 0xbb, 0x49, 0xf8, 0x8a, 0x83, 0x9f, 0x99, 0xe5, 0x83,
	0x08, 0x55, 0x8c, 0xaa, 0xab, 0x23, 0x81, 0x39, 0x98, 0x7c, 0x79, 0x8f, 0x23, 0x47, 0xa3, 0x67,
	0x7f, 0x46, 0xad, 0xbd, 0xe4, 0x58, 0xfe, 0x42, 0x63, 0x9d, 0x4b, 0x56, 0x10, 0xf2, 0x7e, 0x14,
	0x92, 0x40, 0xe9, 0xda, 0x55, 0xbb, 0x5e, 0x68, 0x9e, 0xcc, 0x16, 0x15, 0xeb, 0x7d, 0x51, 0x39,
	0x34, 0x24, 0xd5, 0x1f, 0xfa, 0x02, 0x83, 0x38, 0xa4, 0x81, 0xdf, 0x06, 0x1e, 0x46, 0xd3, 0x16,
	0x44, 0xaf, 0xcf, 0x0d, 0xf6, 0x5d, 0xd4, 0x82, 0xa8, 0xb3, 0x61, 0x38, 0x47, 0x8c, 0x25, 0x83,
	0x50, 0x41, 0x57, 0x11, 0x24, 0x6e, 0xae, 0x6a, 0xd7, 0x77, 0x3a, 0x05, 0xad, 0x5c, 0x13, 0x24,
	0xce, 0x84, 0x95, 0x8c, 0x9d, 0xa4, 0xf8, 0x20, 0x94, 0x40, 0xa9, 0xdc, 0x7f, 0xba, 0xb6, 0xfd,
	0x8b, 0xda, 0xcf, 0x45, 0x65, 0x7f, 0x1a, 0xc6, 0xa3, 0xb3, 0xda, 0x36, 0xa4, 0xb6, 0x75, 0xa3,
	0xa2, 0x0e, 0x5c, 0xad, 0x7d, 0xe7, 0x8e, 0x95, 0x28, 0x1d, 0xcb, 0x28, 0x24, 0xe8, 0x77, 0x09,
	0x87, 0x20, 0x95, 0xfb, 0xff, 0xaf, 0xef, 0x2d, 0xae, 0x51, 0x37, 0x9a, 0xd4, 0x3c, 0x9f, 0x2d,
	0x3d, 0x7b, 0xbe, 0xf4, 0xec, 0x8f, 0xa5, 0x67, 0x3f, 0xad, 0x3c, 0x6b, 0xbe, 0xf2, 0xac, 0xb7,
	0x95, 0x67, 0xdd, 0x1e, 0x73, 0x41, 0x83, 0x71, 0xcf, 0x8f, 0x30, 0xce, 0x36, 0x6d, 0xfc, 0xdc,
	0xf7, 0xd1, 0x2c, 0x4c, 0xd3, 0x04, 0x54, 0x2f, 0xaf, 0xb7, 0x39, 0xfd, 0x0a, 0x00, 0x00, 0xff,
	0xff, 0xea, 0x98, 0x4c, 0x79, 0x05, 0x02, 0x00, 0x00,
}

func (m *Minter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Minter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Minter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TruncatedTokens.Size()
		i -= size
		if _, err := m.TruncatedTokens.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMinter(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.PhaseProvisions.Size()
		i -= size
		if _, err := m.PhaseProvisions.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMinter(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.PhaseStep != 0 {
		i = encodeVarintMinter(dAtA, i, uint64(m.PhaseStep))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.Inflation.Size()
		i -= size
		if _, err := m.Inflation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMinter(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintMinter(dAtA []byte, offset int, v uint64) int {
	offset -= sovMinter(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Minter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Inflation.Size()
	n += 1 + l + sovMinter(uint64(l))
	if m.PhaseStep != 0 {
		n += 1 + sovMinter(uint64(m.PhaseStep))
	}
	l = m.PhaseProvisions.Size()
	n += 1 + l + sovMinter(uint64(l))
	l = m.TruncatedTokens.Size()
	n += 1 + l + sovMinter(uint64(l))
	return n
}

func sovMinter(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMinter(x uint64) (n int) {
	return sovMinter(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Minter) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMinter
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
			return fmt.Errorf("proto: Minter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Minter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Inflation", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMinter
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
				return ErrInvalidLengthMinter
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMinter
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Inflation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PhaseStep", wireType)
			}
			m.PhaseStep = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMinter
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PhaseStep |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PhaseProvisions", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMinter
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
				return ErrInvalidLengthMinter
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMinter
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PhaseProvisions.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TruncatedTokens", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMinter
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
				return ErrInvalidLengthMinter
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMinter
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TruncatedTokens.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMinter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMinter
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
func skipMinter(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMinter
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
					return 0, ErrIntOverflowMinter
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
					return 0, ErrIntOverflowMinter
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
				return 0, ErrInvalidLengthMinter
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMinter
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMinter
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMinter        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMinter          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMinter = fmt.Errorf("proto: unexpected end of group")
)
