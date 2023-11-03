// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/orderbook/participation.proto

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

// OrderBookParticipation represents the participants of an order book.
type OrderBookParticipation struct {
	// index is the index of the participation in the participation queue.
	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty" yaml:"index"`
	// order_book_uid is the unique identifier corresponding to the order book.
	OrderBookUID string `protobuf:"bytes,2,opt,name=order_book_uid,proto3" json:"order_book_uid"`
	// participant_address is the bech32-encoded address of the participant.
	ParticipantAddress string `protobuf:"bytes,3,opt,name=participant_address,json=participantAddress,proto3" json:"participant_address,omitempty" yaml:"participant_address"`
	// liquidity is the total initial liquidity provided.
	Liquidity cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=liquidity,proto3,customtype=cosmossdk.io/math.Int" json:"liquidity" yaml:"liquidity"`
	// fee is the amount of fee to be paid if participation happens.
	Fee cosmossdk_io_math.Int `protobuf:"bytes,5,opt,name=fee,proto3,customtype=cosmossdk.io/math.Int" json:"fee" yaml:"fee"`
	// current_round_liquidity is the liquidity provided for the current round.
	CurrentRoundLiquidity cosmossdk_io_math.Int `protobuf:"bytes,6,opt,name=current_round_liquidity,json=currentRoundLiquidity,proto3,customtype=cosmossdk.io/math.Int" json:"current_round_liquidity" yaml:"current_round_liquidity"`
	// exposures_not_filled represents if all of the exposures of the
	// participation are filled or not.
	ExposuresNotFilled uint64 `protobuf:"varint,7,opt,name=exposures_not_filled,json=exposuresNotFilled,proto3" json:"exposures_not_filled,omitempty" yaml:"exposures_not_filled"`
	// total_bet_amount is the total bet amount corresponding to all exposures.
	TotalBetAmount cosmossdk_io_math.Int `protobuf:"bytes,8,opt,name=total_bet_amount,json=totalBetAmount,proto3,customtype=cosmossdk.io/math.Int" json:"total_bet_amount" yaml:"total_bet_amount"`
	// current_round_total_bet_amount is the total bet amount corresponding to all
	// exposures in the current round.
	CurrentRoundTotalBetAmount cosmossdk_io_math.Int `protobuf:"bytes,9,opt,name=current_round_total_bet_amount,json=currentRoundTotalBetAmount,proto3,customtype=cosmossdk.io/math.Int" json:"current_round_total_bet_amount" yaml:"current_round_total_bet_amount"`
	// max_loss is the total bet amount corresponding to all exposure.
	MaxLoss cosmossdk_io_math.Int `protobuf:"bytes,10,opt,name=max_loss,json=maxLoss,proto3,customtype=cosmossdk.io/math.Int" json:"max_loss" yaml:"max_loss"`
	// current_round_max_loss is the current round max loss.
	CurrentRoundMaxLoss cosmossdk_io_math.Int `protobuf:"bytes,11,opt,name=current_round_max_loss,json=currentRoundMaxLoss,proto3,customtype=cosmossdk.io/math.Int" json:"current_round_max_loss" yaml:"current_round_max_loss"`
	// current_round_max_loss_odds_uid is the total max loss corresponding to
	// all exposures.
	CurrentRoundMaxLossOddsUID string `protobuf:"bytes,12,opt,name=current_round_max_loss_odds_uid,proto3" json:"current_round_max_loss_odds_uid" yaml:"current_round_max_loss_odds_uid"`
	// actual_profit is the actual profit of the participation fulfillment.
	ActualProfit cosmossdk_io_math.Int `protobuf:"bytes,13,opt,name=actual_profit,json=actualProfit,proto3,customtype=cosmossdk.io/math.Int" json:"actual_profit" yaml:"actual_profit"`
	// is_settled represents if the participation is settled or not.
	IsSettled bool `protobuf:"varint,14,opt,name=is_settled,json=isSettled,proto3" json:"is_settled,omitempty" yaml:"is_settled"`
}

func (m *OrderBookParticipation) Reset()      { *m = OrderBookParticipation{} }
func (*OrderBookParticipation) ProtoMessage() {}
func (*OrderBookParticipation) Descriptor() ([]byte, []int) {
	return fileDescriptor_2962bcb47b63c36a, []int{0}
}
func (m *OrderBookParticipation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderBookParticipation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderBookParticipation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderBookParticipation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderBookParticipation.Merge(m, src)
}
func (m *OrderBookParticipation) XXX_Size() int {
	return m.Size()
}
func (m *OrderBookParticipation) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderBookParticipation.DiscardUnknown(m)
}

var xxx_messageInfo_OrderBookParticipation proto.InternalMessageInfo

// ParticipationBetPair represents the book participation and bet bond.
type ParticipationBetPair struct {
	// order_book_uid is the unique identifier corresponding to the order book
	OrderBookUID string `protobuf:"bytes,1,opt,name=order_book_uid,proto3" json:"order_book_uid"`
	// participation_index is the index of participation corresponding to the bet
	// fulfillment.
	ParticipationIndex uint64 `protobuf:"varint,2,opt,name=participation_index,json=participationIndex,proto3" json:"participation_index,omitempty" yaml:"participation_index"`
	// bet_uid is the bet universal unique identifier of the bet that is
	// fulfilled.
	BetUID string `protobuf:"bytes,3,opt,name=bet_uid,proto3" json:"bet_uid"`
}

func (m *ParticipationBetPair) Reset()         { *m = ParticipationBetPair{} }
func (m *ParticipationBetPair) String() string { return proto.CompactTextString(m) }
func (*ParticipationBetPair) ProtoMessage()    {}
func (*ParticipationBetPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_2962bcb47b63c36a, []int{1}
}
func (m *ParticipationBetPair) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ParticipationBetPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ParticipationBetPair.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ParticipationBetPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParticipationBetPair.Merge(m, src)
}
func (m *ParticipationBetPair) XXX_Size() int {
	return m.Size()
}
func (m *ParticipationBetPair) XXX_DiscardUnknown() {
	xxx_messageInfo_ParticipationBetPair.DiscardUnknown(m)
}

var xxx_messageInfo_ParticipationBetPair proto.InternalMessageInfo

func (m *ParticipationBetPair) GetOrderBookUID() string {
	if m != nil {
		return m.OrderBookUID
	}
	return ""
}

func (m *ParticipationBetPair) GetParticipationIndex() uint64 {
	if m != nil {
		return m.ParticipationIndex
	}
	return 0
}

func (m *ParticipationBetPair) GetBetUID() string {
	if m != nil {
		return m.BetUID
	}
	return ""
}

func init() {
	proto.RegisterType((*OrderBookParticipation)(nil), "sgenetwork.sge.orderbook.OrderBookParticipation")
	proto.RegisterType((*ParticipationBetPair)(nil), "sgenetwork.sge.orderbook.ParticipationBetPair")
}

func init() { proto.RegisterFile("sge/orderbook/participation.proto", fileDescriptor_2962bcb47b63c36a) }

var fileDescriptor_2962bcb47b63c36a = []byte{
	// 740 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x9b, 0xfd, 0xaf, 0xe9, 0xc6, 0xc8, 0xda, 0x2d, 0x2a, 0x22, 0x1e, 0x96, 0x98, 0x26,
	0x01, 0x2d, 0x08, 0x90, 0xd0, 0x24, 0x40, 0x0b, 0x68, 0xd2, 0x60, 0xd0, 0x11, 0xe0, 0xb2, 0x8b,
	0x95, 0x36, 0x6e, 0x67, 0xb5, 0x8d, 0x4b, 0xec, 0x68, 0xdd, 0x37, 0x40, 0x9c, 0x38, 0x72, 0xdc,
	0x97, 0xe0, 0x1b, 0x70, 0xd8, 0x71, 0x47, 0xc4, 0x21, 0x42, 0x1d, 0x07, 0xc4, 0xb1, 0x9f, 0x00,
	0xc5, 0xee, 0xd2, 0x76, 0x0b, 0xe4, 0xc2, 0xcd, 0x79, 0xfd, 0xbe, 0xbf, 0xe7, 0x79, 0xdd, 0xfa,
	0x35, 0xb8, 0xce, 0x1b, 0xa4, 0xcc, 0x7c, 0x97, 0xf8, 0x55, 0xc6, 0x9a, 0xe5, 0x8e, 0xe3, 0x0b,
	0x5a, 0xa3, 0x1d, 0x47, 0x50, 0xe6, 0x95, 0x3a, 0x3e, 0x13, 0x4c, 0x37, 0x78, 0x83, 0x78, 0x44,
	0x1c, 0x30, 0xbf, 0x59, 0xe2, 0x0d, 0x52, 0x8a, 0xb3, 0x8b, 0xf9, 0x06, 0x6b, 0x30, 0x99, 0x54,
	0x8e, 0x56, 0x2a, 0x1f, 0x7d, 0x05, 0x60, 0xb9, 0x12, 0xe5, 0x58, 0x8c, 0x35, 0x77, 0x47, 0x81,
	0xfa, 0x1a, 0x98, 0xa6, 0x9e, 0x4b, 0xba, 0x86, 0xb6, 0xaa, 0xad, 0x4f, 0x59, 0x8b, 0xfd, 0x10,
	0xe6, 0x0e, 0x9d, 0x76, 0x6b, 0x03, 0xc9, 0x30, 0xb2, 0xd5, 0xb6, 0xfe, 0x1c, 0x2c, 0x48, 0x15,
	0x1c, 0xc9, 0xe0, 0x80, 0xba, 0xc6, 0xc4, 0xaa, 0xb6, 0x9e, 0xb5, 0x50, 0x2f, 0x84, 0xb9, 0x98,
	0xfd, 0x6e, 0xfb, 0xd9, 0xef, 0x10, 0x9e, 0xcb, 0xb4, 0xcf, 0x7d, 0xeb, 0x15, 0xb0, 0x14, 0x77,
	0xe5, 0x09, 0xec, 0xb8, 0xae, 0x4f, 0x38, 0x37, 0x26, 0x25, 0xd0, 0xec, 0x87, 0xb0, 0xa8, 0x1c,
	0x24, 0x24, 0x21, 0x5b, 0x1f, 0x89, 0x6e, 0xaa, 0xa0, 0x5e, 0x01, 0xd9, 0x16, 0x7d, 0x1f, 0x50,
	0x97, 0x8a, 0x43, 0x63, 0x4a, 0x62, 0xee, 0x1e, 0x87, 0x30, 0xf3, 0x3d, 0x84, 0x85, 0x1a, 0xe3,
	0x6d, 0xc6, 0xb9, 0xdb, 0x2c, 0x51, 0x56, 0x6e, 0x3b, 0x62, 0xbf, 0xb4, 0xed, 0x89, 0x7e, 0x08,
	0x17, 0x95, 0x46, 0x5c, 0x87, 0xec, 0x21, 0x43, 0x7f, 0x04, 0x26, 0xeb, 0x84, 0x18, 0xd3, 0x12,
	0x75, 0x33, 0x0d, 0x05, 0x14, 0xaa, 0x4e, 0x08, 0xb2, 0xa3, 0x3a, 0xfd, 0x00, 0xac, 0xd4, 0x02,
	0xdf, 0x27, 0x9e, 0xc0, 0x3e, 0x0b, 0x3c, 0x17, 0x0f, 0xdd, 0xcd, 0x48, 0xe4, 0x93, 0x34, 0xa4,
	0xa9, 0x90, 0x7f, 0xa1, 0x20, 0xbb, 0x30, 0xd8, 0xb1, 0xa3, 0x8d, 0x9d, 0xd8, 0xf7, 0x6b, 0x90,
	0x27, 0xdd, 0x0e, 0xe3, 0x81, 0x4f, 0x38, 0xf6, 0x98, 0xc0, 0x75, 0xda, 0x6a, 0x11, 0xd7, 0x98,
	0x95, 0x3f, 0x2e, 0xec, 0x87, 0xf0, 0xaa, 0x02, 0x27, 0x65, 0x21, 0x5b, 0x8f, 0xc3, 0xaf, 0x98,
	0xd8, 0x92, 0x41, 0xbd, 0x0a, 0x16, 0x05, 0x13, 0x4e, 0x0b, 0x57, 0x89, 0xc0, 0x4e, 0x9b, 0x05,
	0x9e, 0x30, 0xe6, 0x64, 0x13, 0x0f, 0xd3, 0x9a, 0x58, 0x51, 0x5a, 0xe7, 0xcb, 0x91, 0xbd, 0x20,
	0x43, 0x16, 0x11, 0x9b, 0x32, 0xa0, 0x7f, 0xd4, 0x80, 0x39, 0xde, 0xea, 0x05, 0xc9, 0xac, 0x94,
	0xdc, 0x4a, 0x93, 0xbc, 0x91, 0x74, 0x6e, 0x17, 0x0d, 0x14, 0x47, 0x8f, 0xef, 0xed, 0xb8, 0x99,
	0x17, 0x60, 0xae, 0xed, 0x74, 0x71, 0x8b, 0x71, 0x6e, 0x00, 0xa9, 0x7a, 0x27, 0x4d, 0xf5, 0xb2,
	0x52, 0x3d, 0x2b, 0x43, 0xf6, 0x6c, 0xdb, 0xe9, 0xee, 0x30, 0xce, 0x75, 0x0e, 0x96, 0xc7, 0xbd,
	0xc4, 0xe8, 0x4b, 0x12, 0xfd, 0x38, 0x0d, 0x7d, 0x2d, 0xa9, 0xa1, 0xa1, 0xd0, 0xd2, 0x68, 0x23,
	0x2f, 0x07, 0xa2, 0x5f, 0x34, 0x00, 0x93, 0x0b, 0x30, 0x73, 0x5d, 0x2e, 0x6f, 0x6f, 0x4e, 0xca,
	0x37, 0x7b, 0x21, 0x2c, 0x3e, 0xbd, 0x88, 0xa8, 0xb8, 0x2e, 0x57, 0x77, 0x39, 0x0d, 0xd4, 0x0f,
	0xe1, 0xda, 0xbf, 0x2c, 0xc6, 0x89, 0xc8, 0x4e, 0x43, 0xe9, 0x7b, 0x60, 0xde, 0xa9, 0x89, 0xc0,
	0x69, 0xe1, 0x8e, 0xcf, 0xea, 0x54, 0x18, 0xf3, 0xd2, 0xe4, 0x83, 0xb4, 0x33, 0xca, 0x2b, 0x03,
	0x63, 0xb5, 0xc8, 0xce, 0xa9, 0xef, 0x5d, 0xf9, 0xa9, 0xdf, 0x07, 0x80, 0x72, 0xcc, 0x89, 0x10,
	0xd1, 0x7d, 0x58, 0x58, 0xd5, 0xd6, 0xe7, 0xac, 0x42, 0x3f, 0x84, 0x57, 0x06, 0xc3, 0x2e, 0xde,
	0x43, 0x76, 0x96, 0xf2, 0x37, 0x6a, 0xbd, 0x91, 0xfb, 0x70, 0x04, 0x33, 0x9f, 0x8f, 0x60, 0xe6,
	0xd7, 0x11, 0xcc, 0xa0, 0x9f, 0x1a, 0xc8, 0x8f, 0x4d, 0x4f, 0x8b, 0x88, 0x5d, 0x87, 0xfa, 0x09,
	0xc3, 0x51, 0xfb, 0x2f, 0xc3, 0x31, 0xd2, 0xc0, 0x6a, 0x3c, 0x4f, 0xc8, 0x1b, 0x9c, 0x34, 0x1c,
	0x87, 0x49, 0xa3, 0xc3, 0x31, 0x8a, 0x6e, 0xcb, 0xc9, 0x5d, 0x06, 0xb3, 0xd1, 0x5f, 0x3f, 0x72,
	0xa5, 0x26, 0x6c, 0xa1, 0x17, 0xc2, 0x19, 0x8b, 0x08, 0xe5, 0xe7, 0x6c, 0xd3, 0x3e, 0x5b, 0x58,
	0x5b, 0xc7, 0x3d, 0x53, 0x3b, 0xe9, 0x99, 0xda, 0x8f, 0x9e, 0xa9, 0x7d, 0x3a, 0x35, 0x33, 0x27,
	0xa7, 0x66, 0xe6, 0xdb, 0xa9, 0x99, 0xd9, 0xbb, 0xd5, 0xa0, 0x62, 0x3f, 0xa8, 0x96, 0x6a, 0xac,
	0x5d, 0xe6, 0x0d, 0x72, 0x7b, 0xf0, 0x06, 0x45, 0xeb, 0x72, 0x77, 0xe4, 0xcd, 0x12, 0x87, 0x1d,
	0xc2, 0xab, 0x33, 0xf2, 0xf1, 0xb9, 0xf7, 0x27, 0x00, 0x00, 0xff, 0xff, 0xa2, 0xb1, 0xc3, 0x6d,
	0xd1, 0x06, 0x00, 0x00,
}

func (m *OrderBookParticipation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderBookParticipation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderBookParticipation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsSettled {
		i--
		if m.IsSettled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x70
	}
	{
		size := m.ActualProfit.Size()
		i -= size
		if _, err := m.ActualProfit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	if len(m.CurrentRoundMaxLossOddsUID) > 0 {
		i -= len(m.CurrentRoundMaxLossOddsUID)
		copy(dAtA[i:], m.CurrentRoundMaxLossOddsUID)
		i = encodeVarintParticipation(dAtA, i, uint64(len(m.CurrentRoundMaxLossOddsUID)))
		i--
		dAtA[i] = 0x62
	}
	{
		size := m.CurrentRoundMaxLoss.Size()
		i -= size
		if _, err := m.CurrentRoundMaxLoss.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	{
		size := m.MaxLoss.Size()
		i -= size
		if _, err := m.MaxLoss.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.CurrentRoundTotalBetAmount.Size()
		i -= size
		if _, err := m.CurrentRoundTotalBetAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.TotalBetAmount.Size()
		i -= size
		if _, err := m.TotalBetAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	if m.ExposuresNotFilled != 0 {
		i = encodeVarintParticipation(dAtA, i, uint64(m.ExposuresNotFilled))
		i--
		dAtA[i] = 0x38
	}
	{
		size := m.CurrentRoundLiquidity.Size()
		i -= size
		if _, err := m.CurrentRoundLiquidity.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.Fee.Size()
		i -= size
		if _, err := m.Fee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.Liquidity.Size()
		i -= size
		if _, err := m.Liquidity.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParticipation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.ParticipantAddress) > 0 {
		i -= len(m.ParticipantAddress)
		copy(dAtA[i:], m.ParticipantAddress)
		i = encodeVarintParticipation(dAtA, i, uint64(len(m.ParticipantAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.OrderBookUID) > 0 {
		i -= len(m.OrderBookUID)
		copy(dAtA[i:], m.OrderBookUID)
		i = encodeVarintParticipation(dAtA, i, uint64(len(m.OrderBookUID)))
		i--
		dAtA[i] = 0x12
	}
	if m.Index != 0 {
		i = encodeVarintParticipation(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ParticipationBetPair) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ParticipationBetPair) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ParticipationBetPair) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BetUID) > 0 {
		i -= len(m.BetUID)
		copy(dAtA[i:], m.BetUID)
		i = encodeVarintParticipation(dAtA, i, uint64(len(m.BetUID)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ParticipationIndex != 0 {
		i = encodeVarintParticipation(dAtA, i, uint64(m.ParticipationIndex))
		i--
		dAtA[i] = 0x10
	}
	if len(m.OrderBookUID) > 0 {
		i -= len(m.OrderBookUID)
		copy(dAtA[i:], m.OrderBookUID)
		i = encodeVarintParticipation(dAtA, i, uint64(len(m.OrderBookUID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParticipation(dAtA []byte, offset int, v uint64) int {
	offset -= sovParticipation(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *OrderBookParticipation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Index != 0 {
		n += 1 + sovParticipation(uint64(m.Index))
	}
	l = len(m.OrderBookUID)
	if l > 0 {
		n += 1 + l + sovParticipation(uint64(l))
	}
	l = len(m.ParticipantAddress)
	if l > 0 {
		n += 1 + l + sovParticipation(uint64(l))
	}
	l = m.Liquidity.Size()
	n += 1 + l + sovParticipation(uint64(l))
	l = m.Fee.Size()
	n += 1 + l + sovParticipation(uint64(l))
	l = m.CurrentRoundLiquidity.Size()
	n += 1 + l + sovParticipation(uint64(l))
	if m.ExposuresNotFilled != 0 {
		n += 1 + sovParticipation(uint64(m.ExposuresNotFilled))
	}
	l = m.TotalBetAmount.Size()
	n += 1 + l + sovParticipation(uint64(l))
	l = m.CurrentRoundTotalBetAmount.Size()
	n += 1 + l + sovParticipation(uint64(l))
	l = m.MaxLoss.Size()
	n += 1 + l + sovParticipation(uint64(l))
	l = m.CurrentRoundMaxLoss.Size()
	n += 1 + l + sovParticipation(uint64(l))
	l = len(m.CurrentRoundMaxLossOddsUID)
	if l > 0 {
		n += 1 + l + sovParticipation(uint64(l))
	}
	l = m.ActualProfit.Size()
	n += 1 + l + sovParticipation(uint64(l))
	if m.IsSettled {
		n += 2
	}
	return n
}

func (m *ParticipationBetPair) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OrderBookUID)
	if l > 0 {
		n += 1 + l + sovParticipation(uint64(l))
	}
	if m.ParticipationIndex != 0 {
		n += 1 + sovParticipation(uint64(m.ParticipationIndex))
	}
	l = len(m.BetUID)
	if l > 0 {
		n += 1 + l + sovParticipation(uint64(l))
	}
	return n
}

func sovParticipation(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParticipation(x uint64) (n int) {
	return sovParticipation(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *OrderBookParticipation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParticipation
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
			return fmt.Errorf("proto: OrderBookParticipation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderBookParticipation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderBookUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrderBookUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParticipantAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ParticipantAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liquidity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Fee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentRoundLiquidity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CurrentRoundLiquidity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExposuresNotFilled", wireType)
			}
			m.ExposuresNotFilled = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExposuresNotFilled |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBetAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalBetAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentRoundTotalBetAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CurrentRoundTotalBetAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxLoss", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxLoss.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentRoundMaxLoss", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CurrentRoundMaxLoss.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentRoundMaxLossOddsUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurrentRoundMaxLossOddsUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActualProfit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ActualProfit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsSettled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
			m.IsSettled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipParticipation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParticipation
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
func (m *ParticipationBetPair) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParticipation
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
			return fmt.Errorf("proto: ParticipationBetPair: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ParticipationBetPair: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderBookUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrderBookUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParticipationIndex", wireType)
			}
			m.ParticipationIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BetUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParticipation
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
				return ErrInvalidLengthParticipation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParticipation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BetUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParticipation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParticipation
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
func skipParticipation(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParticipation
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
					return 0, ErrIntOverflowParticipation
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
					return 0, ErrIntOverflowParticipation
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
				return 0, ErrInvalidLengthParticipation
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParticipation
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParticipation
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParticipation        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParticipation          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParticipation = fmt.Errorf("proto: unexpected end of group")
)
