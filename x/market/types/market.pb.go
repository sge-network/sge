// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/market/market.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

// MarketStatus is the market status enumeration
type MarketStatus int32

const (
	// unspecified market
	MarketStatus_MARKET_STATUS_UNSPECIFIED MarketStatus = 0
	// market is active
	MarketStatus_MARKET_STATUS_ACTIVE MarketStatus = 1
	// market is inactive
	MarketStatus_MARKET_STATUS_INACTIVE MarketStatus = 2
	// market is canceled
	MarketStatus_MARKET_STATUS_CANCELED MarketStatus = 3
	// market is aborted
	MarketStatus_MARKET_STATUS_ABORTED MarketStatus = 4
	// result of the market is declared
	MarketStatus_MARKET_STATUS_RESULT_DECLARED MarketStatus = 5
)

var MarketStatus_name = map[int32]string{
	0: "MARKET_STATUS_UNSPECIFIED",
	1: "MARKET_STATUS_ACTIVE",
	2: "MARKET_STATUS_INACTIVE",
	3: "MARKET_STATUS_CANCELED",
	4: "MARKET_STATUS_ABORTED",
	5: "MARKET_STATUS_RESULT_DECLARED",
}

var MarketStatus_value = map[string]int32{
	"MARKET_STATUS_UNSPECIFIED":     0,
	"MARKET_STATUS_ACTIVE":          1,
	"MARKET_STATUS_INACTIVE":        2,
	"MARKET_STATUS_CANCELED":        3,
	"MARKET_STATUS_ABORTED":         4,
	"MARKET_STATUS_RESULT_DECLARED": 5,
}

func (x MarketStatus) String() string {
	return proto.EnumName(MarketStatus_name, int32(x))
}

func (MarketStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_935a8ad1d6bee065, []int{0}
}

// Market is the representation of the market to be stored in
// the market state.
type Market struct {
	// uid is the universal unique identifier of the market.
	UID string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid"`
	// start_ts is the start timestamp of the market.
	StartTS uint64 `protobuf:"varint,2,opt,name=start_ts,proto3" json:"start_ts"`
	// end_ts is the end timestamp of the market.
	EndTS uint64 `protobuf:"varint,3,opt,name=end_ts,proto3" json:"end_ts"`
	// odds is the list of odds of the market.
	Odds []*Odds `protobuf:"bytes,4,rep,name=odds,proto3" json:"odds,omitempty"`
	// winner_odds_uids is the list of winner odds universal unique identifiers.
	WinnerOddsUIDs []string `protobuf:"bytes,5,rep,name=winner_odds_uids,proto3" json:"winner_odds_uids"`
	// status is the current status of the market.
	Status MarketStatus `protobuf:"varint,6,opt,name=status,proto3,enum=sgenetwork.sge.market.MarketStatus" json:"status,omitempty"`
	// resolution_ts is the timestamp of the resolution of market.
	ResolutionTS uint64 `protobuf:"varint,7,opt,name=resolution_ts,proto3" json:"resolution_ts"`
	// creator is the address of the creator of market.
	Creator string `protobuf:"bytes,8,opt,name=creator,proto3" json:"creator,omitempty"`
	// meta contains human-readable metadata of the market.
	Meta string `protobuf:"bytes,9,opt,name=meta,proto3" json:"meta,omitempty"`
	// book_uid is the unique identifier corresponding to the book
	BookUID string `protobuf:"bytes,10,opt,name=book_uid,proto3" json:"book_uid"`
	// max_total_payout is the sum of total payout to be paid after settlement if all of bets win.
	MaxTotalPayout cosmossdk_io_math.Int `protobuf:"bytes,11,opt,name=max_total_payout,json=maxTotalPayout,proto3,customtype=cosmossdk.io/math.Int" json:"max_total_payout"`
	// price_stats is statistics related to the price fluctuation.
	PriceStats *PriceStats `protobuf:"bytes,12,opt,name=price_stats,json=priceStats,proto3" json:"price_stats,omitempty"`
	// price_pool is price fund pool information.
	PricePool *PricePool `protobuf:"bytes,13,opt,name=price_pool,json=pricePool,proto3" json:"price_pool,omitempty"`
}

func (m *Market) Reset()         { *m = Market{} }
func (m *Market) String() string { return proto.CompactTextString(m) }
func (*Market) ProtoMessage()    {}
func (*Market) Descriptor() ([]byte, []int) {
	return fileDescriptor_935a8ad1d6bee065, []int{0}
}
func (m *Market) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Market) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Market.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Market) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Market.Merge(m, src)
}
func (m *Market) XXX_Size() int {
	return m.Size()
}
func (m *Market) XXX_DiscardUnknown() {
	xxx_messageInfo_Market.DiscardUnknown(m)
}

var xxx_messageInfo_Market proto.InternalMessageInfo

func (m *Market) GetUID() string {
	if m != nil {
		return m.UID
	}
	return ""
}

func (m *Market) GetStartTS() uint64 {
	if m != nil {
		return m.StartTS
	}
	return 0
}

func (m *Market) GetEndTS() uint64 {
	if m != nil {
		return m.EndTS
	}
	return 0
}

func (m *Market) GetOdds() []*Odds {
	if m != nil {
		return m.Odds
	}
	return nil
}

func (m *Market) GetWinnerOddsUIDs() []string {
	if m != nil {
		return m.WinnerOddsUIDs
	}
	return nil
}

func (m *Market) GetStatus() MarketStatus {
	if m != nil {
		return m.Status
	}
	return MarketStatus_MARKET_STATUS_UNSPECIFIED
}

func (m *Market) GetResolutionTS() uint64 {
	if m != nil {
		return m.ResolutionTS
	}
	return 0
}

func (m *Market) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *Market) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func (m *Market) GetBookUID() string {
	if m != nil {
		return m.BookUID
	}
	return ""
}

func (m *Market) GetPriceStats() *PriceStats {
	if m != nil {
		return m.PriceStats
	}
	return nil
}

func (m *Market) GetPricePool() *PricePool {
	if m != nil {
		return m.PricePool
	}
	return nil
}

// PriceStats is a type for the sge price fluctuation definitions and statistics.
type PriceStats struct {
	// max_wager_sge_price is the maximum price of sge token in dollars during wager among all of the bets.
	MaxWagerSgePrice github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=max_wager_sge_price,json=maxWagerSgePrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"max_wager_sge_price"`
	// sge_price is the price of sge token in dollars during resolution.
	ResolutionSgePrice github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=resolution_sge_price,json=resolutionSgePrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"resolution_sge_price"`
}

func (m *PriceStats) Reset()         { *m = PriceStats{} }
func (m *PriceStats) String() string { return proto.CompactTextString(m) }
func (*PriceStats) ProtoMessage()    {}
func (*PriceStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_935a8ad1d6bee065, []int{1}
}
func (m *PriceStats) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PriceStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PriceStats.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PriceStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PriceStats.Merge(m, src)
}
func (m *PriceStats) XXX_Size() int {
	return m.Size()
}
func (m *PriceStats) XXX_DiscardUnknown() {
	xxx_messageInfo_PriceStats.DiscardUnknown(m)
}

var xxx_messageInfo_PriceStats proto.InternalMessageInfo

// PricePool is a type for the sge price fluctuation pool stats.
type PricePool struct {
	// resolution_funds is the funds allocated for the price lock reimbursement on resolution.
	ResolutionFunds cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=resolution_funds,json=resolutionFunds,proto3,customtype=cosmossdk.io/math.Int" json:"resolution_funds"`
	// sge_price is the funds spent on automated orderbook settlement.
	SpentFunds cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=spent_funds,json=spentFunds,proto3,customtype=cosmossdk.io/math.Int" json:"spent_funds"`
	// returned_funds is the remaining funds to be returned to the market creator after orderbook settlement.
	ReturnedFunds cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=returned_funds,json=returnedFunds,proto3,customtype=cosmossdk.io/math.Int" json:"returned_funds"`
}

func (m *PricePool) Reset()         { *m = PricePool{} }
func (m *PricePool) String() string { return proto.CompactTextString(m) }
func (*PricePool) ProtoMessage()    {}
func (*PricePool) Descriptor() ([]byte, []int) {
	return fileDescriptor_935a8ad1d6bee065, []int{2}
}
func (m *PricePool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PricePool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PricePool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PricePool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PricePool.Merge(m, src)
}
func (m *PricePool) XXX_Size() int {
	return m.Size()
}
func (m *PricePool) XXX_DiscardUnknown() {
	xxx_messageInfo_PricePool.DiscardUnknown(m)
}

var xxx_messageInfo_PricePool proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("sgenetwork.sge.market.MarketStatus", MarketStatus_name, MarketStatus_value)
	proto.RegisterType((*Market)(nil), "sgenetwork.sge.market.Market")
	proto.RegisterType((*PriceStats)(nil), "sgenetwork.sge.market.PriceStats")
	proto.RegisterType((*PricePool)(nil), "sgenetwork.sge.market.PricePool")
}

func init() { proto.RegisterFile("sge/market/market.proto", fileDescriptor_935a8ad1d6bee065) }

var fileDescriptor_935a8ad1d6bee065 = []byte{
	// 770 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xdd, 0x6e, 0xe3, 0x44,
	0x18, 0x8d, 0xf3, 0xd7, 0xcd, 0xa4, 0x0d, 0xd6, 0x6c, 0xc3, 0x7a, 0x8b, 0x1a, 0xbb, 0x45, 0x5a,
	0x05, 0xd0, 0x3a, 0x52, 0xf7, 0x12, 0x09, 0x94, 0xc4, 0xde, 0x25, 0xa2, 0x3f, 0xd1, 0xd8, 0xa1,
	0x12, 0x12, 0x32, 0x6e, 0x3c, 0xb8, 0x51, 0x12, 0x8f, 0xe5, 0x19, 0xab, 0xed, 0x5b, 0xf0, 0x42,
	0xdc, 0xf7, 0x06, 0xa9, 0x97, 0x08, 0x24, 0x0b, 0xb9, 0x77, 0x7d, 0x0a, 0x34, 0x63, 0xe7, 0x8f,
	0x52, 0xa8, 0xb8, 0xc9, 0x7c, 0xf3, 0x9d, 0xef, 0x9c, 0xf1, 0x99, 0xe3, 0x18, 0xbc, 0xa2, 0x3e,
	0xee, 0xcc, 0xdd, 0x68, 0x8a, 0x59, 0xbe, 0xe8, 0x61, 0x44, 0x18, 0x81, 0x4d, 0xea, 0xe3, 0x00,
	0xb3, 0x2b, 0x12, 0x4d, 0x75, 0xea, 0x63, 0x3d, 0x03, 0xf7, 0x76, 0x7d, 0xe2, 0x13, 0x31, 0xd1,
	0xe1, 0x55, 0x36, 0xbc, 0xd7, 0x5c, 0x53, 0x21, 0x9e, 0x47, 0xb3, 0xf6, 0xe1, 0x6d, 0x05, 0x54,
	0x4f, 0x44, 0x17, 0x6a, 0xa0, 0x14, 0x4f, 0x3c, 0x45, 0xd2, 0xa4, 0x76, 0xad, 0xd7, 0x48, 0x13,
	0xb5, 0x34, 0x1a, 0x18, 0x0f, 0x89, 0xca, 0xbb, 0x88, 0xff, 0xc0, 0x77, 0xe0, 0x05, 0x65, 0x6e,
	0xc4, 0x1c, 0x46, 0x95, 0xa2, 0x26, 0xb5, 0xcb, 0xbd, 0x57, 0x69, 0xa2, 0x6e, 0x59, 0xbc, 0x67,
	0x5b, 0x0f, 0x89, 0xba, 0x84, 0xd1, 0xb2, 0x82, 0x5f, 0x80, 0x2a, 0x0e, 0x3c, 0x4e, 0x29, 0x09,
	0xca, 0xcb, 0x34, 0x51, 0x2b, 0x66, 0xe0, 0x09, 0x42, 0x0e, 0xa1, 0x7c, 0x85, 0x1d, 0x50, 0xe6,
	0x0f, 0xa7, 0x94, 0xb5, 0x52, 0xbb, 0x7e, 0xf4, 0x89, 0xfe, 0x8f, 0x0e, 0xf5, 0x33, 0xcf, 0xa3,
	0x48, 0x0c, 0x42, 0x04, 0xe4, 0xab, 0x49, 0x10, 0xe0, 0xc8, 0xe1, 0x5b, 0x27, 0x9e, 0x78, 0x54,
	0xa9, 0x68, 0xa5, 0x76, 0xad, 0xf7, 0x26, 0x4d, 0xd4, 0xc6, 0xb9, 0xc0, 0xf8, 0xfc, 0x68, 0x60,
	0xd0, 0x87, 0x44, 0x7d, 0x34, 0x8d, 0x1e, 0x75, 0xe0, 0x97, 0xa0, 0x4a, 0x99, 0xcb, 0x62, 0xaa,
	0x54, 0x35, 0xa9, 0xdd, 0x38, 0xfa, 0xf4, 0x89, 0xc7, 0xc8, 0xee, 0xcd, 0x12, 0xa3, 0x28, 0xa7,
	0xc0, 0x0f, 0x60, 0x27, 0xc2, 0x94, 0xcc, 0x62, 0x36, 0x21, 0x01, 0x77, 0xbd, 0x25, 0x5c, 0x1f,
	0xa4, 0x89, 0xba, 0x8d, 0x96, 0x80, 0x30, 0xbf, 0x39, 0x88, 0x36, 0xb7, 0x50, 0x01, 0x5b, 0xe3,
	0x08, 0xbb, 0x8c, 0x44, 0xca, 0x0b, 0x1e, 0x09, 0x5a, 0x6c, 0x21, 0x04, 0xe5, 0x39, 0x66, 0xae,
	0x52, 0x13, 0x6d, 0x51, 0xf3, 0x68, 0x2e, 0x08, 0x99, 0x72, 0x03, 0x0a, 0x10, 0x09, 0x8a, 0x68,
	0x7a, 0x84, 0x4c, 0xb3, 0x14, 0x97, 0x30, 0x5a, 0x56, 0xf0, 0x03, 0x90, 0xe7, 0xee, 0xb5, 0xc3,
	0x08, 0x73, 0x67, 0x4e, 0xe8, 0xde, 0x90, 0x98, 0x29, 0x75, 0x41, 0xde, 0xbf, 0x4d, 0xd4, 0xc2,
	0xef, 0x89, 0xda, 0x1c, 0x13, 0x3a, 0x27, 0x94, 0x7a, 0x53, 0x7d, 0x42, 0x3a, 0x73, 0x97, 0x5d,
	0xea, 0x83, 0x80, 0xa1, 0xc6, 0xdc, 0xbd, 0xb6, 0x39, 0x6b, 0x28, 0x48, 0xb0, 0x07, 0xea, 0x61,
	0x34, 0x19, 0x63, 0x87, 0x5f, 0x02, 0x55, 0xb6, 0x35, 0xa9, 0x5d, 0x3f, 0x3a, 0x78, 0xe2, 0xda,
	0x86, 0x7c, 0x92, 0xdf, 0x1a, 0x45, 0x20, 0x5c, 0xd6, 0xf0, 0x6b, 0x90, 0xed, 0x9c, 0x90, 0x90,
	0x99, 0xb2, 0x23, 0x24, 0xb4, 0x7f, 0x93, 0x18, 0x12, 0x32, 0x43, 0xb5, 0x70, 0x51, 0x1e, 0xfe,
	0x2a, 0x01, 0xb0, 0xd2, 0x86, 0x3f, 0x80, 0x97, 0xdc, 0xdc, 0x95, 0xeb, 0xe3, 0xc8, 0xa1, 0x3e,
	0x76, 0xc4, 0x64, 0xfe, 0x7a, 0xeb, 0xb9, 0xbf, 0x37, 0xfe, 0x84, 0x5d, 0xc6, 0x17, 0xfa, 0x98,
	0xcc, 0x3b, 0x99, 0xd5, 0x7c, 0x79, 0x4b, 0xbd, 0x69, 0x87, 0xdd, 0x84, 0x98, 0xea, 0x06, 0x1e,
	0x23, 0x7e, 0x4f, 0xe7, 0x5c, 0xc9, 0xf2, 0xb1, 0x38, 0x03, 0xfe, 0x08, 0x76, 0xd7, 0xf2, 0x5a,
	0xe9, 0x17, 0xff, 0x97, 0x3e, 0x5c, 0x69, 0x2d, 0x4e, 0x38, 0xfc, 0x43, 0x02, 0xb5, 0xa5, 0x51,
	0xf8, 0x0d, 0x90, 0xd7, 0xce, 0xfb, 0x29, 0x0e, 0x3c, 0x9a, 0x7b, 0xf9, 0x8f, 0xac, 0x3e, 0x5a,
	0xd1, 0xde, 0x73, 0x16, 0xfc, 0x0a, 0xd4, 0x69, 0x88, 0x03, 0x96, 0x8b, 0x14, 0x9f, 0x23, 0x02,
	0x04, 0x23, 0xe3, 0x1b, 0xa0, 0x11, 0x61, 0x16, 0x47, 0x01, 0xf6, 0x72, 0x89, 0xd2, 0x73, 0x24,
	0x76, 0x16, 0x24, 0xa1, 0xf2, 0xf9, 0x2f, 0x12, 0xd8, 0x5e, 0xff, 0x03, 0xc1, 0x7d, 0xf0, 0xfa,
	0xa4, 0x8b, 0xbe, 0x35, 0x6d, 0xc7, 0xb2, 0xbb, 0xf6, 0xc8, 0x72, 0x46, 0xa7, 0xd6, 0xd0, 0xec,
	0x0f, 0xde, 0x0f, 0x4c, 0x43, 0x2e, 0x40, 0x05, 0xec, 0x6e, 0xc2, 0xdd, 0xbe, 0x3d, 0xf8, 0xce,
	0x94, 0x25, 0xb8, 0x07, 0x3e, 0xde, 0x44, 0x06, 0xa7, 0x39, 0x56, 0x7c, 0x8c, 0xf5, 0xbb, 0xa7,
	0x7d, 0xf3, 0xd8, 0x34, 0xe4, 0x12, 0x7c, 0x0d, 0x9a, 0x7f, 0x53, 0xec, 0x9d, 0x21, 0xdb, 0x34,
	0xe4, 0x32, 0x3c, 0x00, 0xfb, 0x9b, 0x10, 0x32, 0xad, 0xd1, 0xb1, 0xed, 0x18, 0x66, 0xff, 0xb8,
	0x8b, 0x4c, 0x43, 0xae, 0xf4, 0xfa, 0xb7, 0x69, 0x4b, 0xba, 0x4b, 0x5b, 0xd2, 0x9f, 0x69, 0x4b,
	0xfa, 0xf9, 0xbe, 0x55, 0xb8, 0xbb, 0x6f, 0x15, 0x7e, 0xbb, 0x6f, 0x15, 0xbe, 0xff, 0x6c, 0x2d,
	0x73, 0xea, 0xe3, 0xb7, 0xf9, 0xfb, 0xcb, 0xeb, 0xce, 0xf5, 0xe2, 0x13, 0x2c, 0xa2, 0xbf, 0xa8,
	0x8a, 0x8f, 0xf0, 0xbb, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xff, 0x14, 0xe9, 0x3d, 0xe3, 0x05,
	0x00, 0x00,
}

func (m *Market) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Market) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Market) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PricePool != nil {
		{
			size, err := m.PricePool.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarket(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6a
	}
	if m.PriceStats != nil {
		{
			size, err := m.PriceStats.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarket(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x62
	}
	{
		size := m.MaxTotalPayout.Size()
		i -= size
		if _, err := m.MaxTotalPayout.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	if len(m.BookUID) > 0 {
		i -= len(m.BookUID)
		copy(dAtA[i:], m.BookUID)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.BookUID)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Meta) > 0 {
		i -= len(m.Meta)
		copy(dAtA[i:], m.Meta)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.Meta)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x42
	}
	if m.ResolutionTS != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.ResolutionTS))
		i--
		dAtA[i] = 0x38
	}
	if m.Status != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x30
	}
	if len(m.WinnerOddsUIDs) > 0 {
		for iNdEx := len(m.WinnerOddsUIDs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.WinnerOddsUIDs[iNdEx])
			copy(dAtA[i:], m.WinnerOddsUIDs[iNdEx])
			i = encodeVarintMarket(dAtA, i, uint64(len(m.WinnerOddsUIDs[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Odds) > 0 {
		for iNdEx := len(m.Odds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Odds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarket(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.EndTS != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.EndTS))
		i--
		dAtA[i] = 0x18
	}
	if m.StartTS != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.StartTS))
		i--
		dAtA[i] = 0x10
	}
	if len(m.UID) > 0 {
		i -= len(m.UID)
		copy(dAtA[i:], m.UID)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.UID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PriceStats) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PriceStats) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PriceStats) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ResolutionSgePrice.Size()
		i -= size
		if _, err := m.ResolutionSgePrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.MaxWagerSgePrice.Size()
		i -= size
		if _, err := m.MaxWagerSgePrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *PricePool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PricePool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PricePool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ReturnedFunds.Size()
		i -= size
		if _, err := m.ReturnedFunds.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.SpentFunds.Size()
		i -= size
		if _, err := m.SpentFunds.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.ResolutionFunds.Size()
		i -= size
		if _, err := m.ResolutionFunds.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintMarket(dAtA []byte, offset int, v uint64) int {
	offset -= sovMarket(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Market) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.UID)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	if m.StartTS != 0 {
		n += 1 + sovMarket(uint64(m.StartTS))
	}
	if m.EndTS != 0 {
		n += 1 + sovMarket(uint64(m.EndTS))
	}
	if len(m.Odds) > 0 {
		for _, e := range m.Odds {
			l = e.Size()
			n += 1 + l + sovMarket(uint64(l))
		}
	}
	if len(m.WinnerOddsUIDs) > 0 {
		for _, s := range m.WinnerOddsUIDs {
			l = len(s)
			n += 1 + l + sovMarket(uint64(l))
		}
	}
	if m.Status != 0 {
		n += 1 + sovMarket(uint64(m.Status))
	}
	if m.ResolutionTS != 0 {
		n += 1 + sovMarket(uint64(m.ResolutionTS))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	l = len(m.Meta)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	l = len(m.BookUID)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	l = m.MaxTotalPayout.Size()
	n += 1 + l + sovMarket(uint64(l))
	if m.PriceStats != nil {
		l = m.PriceStats.Size()
		n += 1 + l + sovMarket(uint64(l))
	}
	if m.PricePool != nil {
		l = m.PricePool.Size()
		n += 1 + l + sovMarket(uint64(l))
	}
	return n
}

func (m *PriceStats) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MaxWagerSgePrice.Size()
	n += 1 + l + sovMarket(uint64(l))
	l = m.ResolutionSgePrice.Size()
	n += 1 + l + sovMarket(uint64(l))
	return n
}

func (m *PricePool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ResolutionFunds.Size()
	n += 1 + l + sovMarket(uint64(l))
	l = m.SpentFunds.Size()
	n += 1 + l + sovMarket(uint64(l))
	l = m.ReturnedFunds.Size()
	n += 1 + l + sovMarket(uint64(l))
	return n
}

func sovMarket(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMarket(x uint64) (n int) {
	return sovMarket(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Market) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarket
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
			return fmt.Errorf("proto: Market: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Market: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTS", wireType)
			}
			m.StartTS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTS", wireType)
			}
			m.EndTS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Odds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Odds = append(m.Odds, &Odds{})
			if err := m.Odds[len(m.Odds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WinnerOddsUIDs", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WinnerOddsUIDs = append(m.WinnerOddsUIDs, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= MarketStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResolutionTS", wireType)
			}
			m.ResolutionTS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ResolutionTS |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Meta = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BookUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BookUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxTotalPayout", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxTotalPayout.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PriceStats", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PriceStats == nil {
				m.PriceStats = &PriceStats{}
			}
			if err := m.PriceStats.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PricePool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PricePool == nil {
				m.PricePool = &PricePool{}
			}
			if err := m.PricePool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarket
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
func (m *PriceStats) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarket
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
			return fmt.Errorf("proto: PriceStats: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PriceStats: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxWagerSgePrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxWagerSgePrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResolutionSgePrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ResolutionSgePrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarket
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
func (m *PricePool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarket
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
			return fmt.Errorf("proto: PricePool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PricePool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResolutionFunds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ResolutionFunds.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpentFunds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SpentFunds.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReturnedFunds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ReturnedFunds.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarket
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
func skipMarket(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMarket
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
					return 0, ErrIntOverflowMarket
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
					return 0, ErrIntOverflowMarket
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
				return 0, ErrInvalidLengthMarket
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMarket
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMarket
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMarket        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMarket          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMarket = fmt.Errorf("proto: unexpected end of group")
)
