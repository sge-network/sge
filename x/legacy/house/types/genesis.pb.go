// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sgenetwork/sge/house/genesis.proto

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

// GenesisState defines the house module's genesis state.
type GenesisState struct {
	// params defines the parameters of the house module at genesis
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// deposit_list defines the deposits active at genesis.
	DepositList []Deposit `protobuf:"bytes,2,rep,name=deposit_list,json=depositList,proto3" json:"deposit_list"`
	// withdrawal_list defines the withdrawals active at genesis.
	WithdrawalList []Withdrawal `protobuf:"bytes,3,rep,name=withdrawal_list,json=withdrawalList,proto3" json:"withdrawal_list"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_02191d749ca2ced3, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetDepositList() []Deposit {
	if m != nil {
		return m.DepositList
	}
	return nil
}

func (m *GenesisState) GetWithdrawalList() []Withdrawal {
	if m != nil {
		return m.WithdrawalList
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "sgenetwork.sge.house.GenesisState")
}

func init() {
	proto.RegisterFile("sgenetwork/sge/house/genesis.proto", fileDescriptor_02191d749ca2ced3)
}

var fileDescriptor_02191d749ca2ced3 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0x1b, 0x27, 0x3b, 0xb4, 0x43, 0xa1, 0xec, 0x30, 0x8a, 0xc6, 0x3a, 0x2f, 0xbb, 0x98,
	0xc0, 0xbc, 0x79, 0x1c, 0xa2, 0x08, 0x82, 0xa2, 0x07, 0xc1, 0x8b, 0x64, 0x5b, 0x48, 0x83, 0x9d,
	0x29, 0xcd, 0x37, 0xea, 0xfe, 0x85, 0x3f, 0x6b, 0xc7, 0x1d, 0x3d, 0x89, 0xb4, 0x7f, 0x44, 0x9a,
	0x44, 0xb7, 0x43, 0x6e, 0x21, 0x3c, 0xef, 0xf3, 0xbd, 0xbc, 0xe1, 0x50, 0x0b, 0xfe, 0xce, 0xa1,
	0x52, 0xe5, 0x1b, 0xd5, 0x82, 0xd3, 0x4c, 0x2d, 0x35, 0xa7, 0xed, 0x9f, 0x96, 0x9a, 0x14, 0xa5,
	0x02, 0x15, 0xf7, 0xb7, 0x0c, 0xd1, 0x82, 0x13, 0xc3, 0x24, 0x7d, 0xa1, 0x84, 0x32, 0x00, 0x6d,
	0x5f, 0x96, 0x4d, 0xfc, 0xbe, 0x39, 0x2f, 0x94, 0x96, 0xe0, 0x98, 0x53, 0x2f, 0x53, 0xb0, 0x92,
	0x2d, 0xdc, 0xc9, 0xe4, 0xcc, 0x8b, 0x54, 0x12, 0xb2, 0x79, 0xc9, 0x2a, 0x0b, 0x0d, 0x1b, 0x14,
	0xf6, 0x6e, 0x6c, 0xd3, 0x27, 0x60, 0xc0, 0xe3, 0xcb, 0xb0, 0x6b, 0x2d, 0x03, 0x94, 0xa2, 0x51,
	0x34, 0x3e, 0x22, 0xbe, 0xe6, 0xe4, 0xc1, 0x30, 0x93, 0xfd, 0xf5, 0xf7, 0x49, 0xf0, 0xe8, 0x12,
	0xf1, 0x75, 0xd8, 0x73, 0x2d, 0x5f, 0x73, 0xa9, 0x61, 0xb0, 0x97, 0x76, 0x46, 0xd1, 0xf8, 0xd8,
	0x6f, 0xb8, 0xb2, 0xa4, 0x53, 0x44, 0x2e, 0x78, 0x27, 0x35, 0xc4, 0xf7, 0xe1, 0xe1, 0x5f, 0x4d,
	0x96, 0x5b, 0x55, 0xc7, 0xa8, 0x52, 0xbf, 0xea, 0xf9, 0x1f, 0x76, 0xb6, 0x83, 0x6d, 0xbc, 0x15,
	0x4e, 0x6e, 0xd7, 0x35, 0x46, 0x9b, 0x1a, 0xa3, 0x9f, 0x1a, 0xa3, 0xcf, 0x06, 0x07, 0x9b, 0x06,
	0x07, 0x5f, 0x0d, 0x0e, 0x5e, 0xa8, 0x90, 0x90, 0x2d, 0xa7, 0x64, 0xa6, 0x16, 0xed, 0x48, 0xe7,
	0xbb, 0x83, 0x7d, 0xd0, 0x9c, 0x0b, 0x36, 0x5b, 0xb9, 0xe5, 0x60, 0x55, 0x70, 0x3d, 0xed, 0x9a,
	0xdd, 0x2e, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0x45, 0x4f, 0x82, 0x6f, 0xf5, 0x01, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.WithdrawalList) > 0 {
		for iNdEx := len(m.WithdrawalList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.WithdrawalList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.DepositList) > 0 {
		for iNdEx := len(m.DepositList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DepositList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.DepositList) > 0 {
		for _, e := range m.DepositList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.WithdrawalList) > 0 {
		for _, e := range m.WithdrawalList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositList = append(m.DepositList, Deposit{})
			if err := m.DepositList[len(m.DepositList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawalList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WithdrawalList = append(m.WithdrawalList, Withdrawal{})
			if err := m.WithdrawalList[len(m.WithdrawalList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
