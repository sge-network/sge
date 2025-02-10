// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: sge/legacy/reward/v1beta/tx.proto

package rewardv1beta

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Msg_SetPromoterConf_FullMethodName = "/sge.legacy.reward.v1beta.Msg/SetPromoterConf"
	Msg_CreatePromoter_FullMethodName  = "/sge.legacy.reward.v1beta.Msg/CreatePromoter"
	Msg_CreateCampaign_FullMethodName  = "/sge.legacy.reward.v1beta.Msg/CreateCampaign"
	Msg_UpdateCampaign_FullMethodName  = "/sge.legacy.reward.v1beta.Msg/UpdateCampaign"
	Msg_WithdrawFunds_FullMethodName   = "/sge.legacy.reward.v1beta.Msg/WithdrawFunds"
	Msg_GrantReward_FullMethodName     = "/sge.legacy.reward.v1beta.Msg/GrantReward"
	Msg_UpdateParams_FullMethodName    = "/sge.legacy.reward.v1beta.Msg/UpdateParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// Deprecated: Do not use.
	// SetPromoterConf is a method to set the configurations of a promoter.
	SetPromoterConf(ctx context.Context, in *MsgSetPromoterConf, opts ...grpc.CallOption) (*MsgSetPromoterConfResponse, error)
	// Deprecated: Do not use.
	// CreatePromoter is a method to create a promoter
	CreatePromoter(ctx context.Context, in *MsgCreatePromoter, opts ...grpc.CallOption) (*MsgCreatePromoterResponse, error)
	// Deprecated: Do not use.
	// CreateCampaign is a method to create a campaign
	CreateCampaign(ctx context.Context, in *MsgCreateCampaign, opts ...grpc.CallOption) (*MsgCreateCampaignResponse, error)
	// Deprecated: Do not use.
	// UpdateCampaign is a method to update campaign
	UpdateCampaign(ctx context.Context, in *MsgUpdateCampaign, opts ...grpc.CallOption) (*MsgUpdateCampaignResponse, error)
	// Deprecated: Do not use.
	// WithdrawCampaignFunds is method to withdraw funds from the campaign
	WithdrawFunds(ctx context.Context, in *MsgWithdrawFunds, opts ...grpc.CallOption) (*MsgWithdrawFundsResponse, error)
	// Deprecated: Do not use.
	// GrantReward is method to allocate rewards
	GrantReward(ctx context.Context, in *MsgGrantReward, opts ...grpc.CallOption) (*MsgGrantRewardResponse, error)
	// Deprecated: Do not use.
	// UpdateParams defines a governance operation for updating the x/ovm module
	// parameters. The authority is defined in the keeper.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

// Deprecated: Do not use.
func (c *msgClient) SetPromoterConf(ctx context.Context, in *MsgSetPromoterConf, opts ...grpc.CallOption) (*MsgSetPromoterConfResponse, error) {
	out := new(MsgSetPromoterConfResponse)
	err := c.cc.Invoke(ctx, Msg_SetPromoterConf_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *msgClient) CreatePromoter(ctx context.Context, in *MsgCreatePromoter, opts ...grpc.CallOption) (*MsgCreatePromoterResponse, error) {
	out := new(MsgCreatePromoterResponse)
	err := c.cc.Invoke(ctx, Msg_CreatePromoter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *msgClient) CreateCampaign(ctx context.Context, in *MsgCreateCampaign, opts ...grpc.CallOption) (*MsgCreateCampaignResponse, error) {
	out := new(MsgCreateCampaignResponse)
	err := c.cc.Invoke(ctx, Msg_CreateCampaign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *msgClient) UpdateCampaign(ctx context.Context, in *MsgUpdateCampaign, opts ...grpc.CallOption) (*MsgUpdateCampaignResponse, error) {
	out := new(MsgUpdateCampaignResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateCampaign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *msgClient) WithdrawFunds(ctx context.Context, in *MsgWithdrawFunds, opts ...grpc.CallOption) (*MsgWithdrawFundsResponse, error) {
	out := new(MsgWithdrawFundsResponse)
	err := c.cc.Invoke(ctx, Msg_WithdrawFunds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *msgClient) GrantReward(ctx context.Context, in *MsgGrantReward, opts ...grpc.CallOption) (*MsgGrantRewardResponse, error) {
	out := new(MsgGrantRewardResponse)
	err := c.cc.Invoke(ctx, Msg_GrantReward_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// Deprecated: Do not use.
	// SetPromoterConf is a method to set the configurations of a promoter.
	SetPromoterConf(context.Context, *MsgSetPromoterConf) (*MsgSetPromoterConfResponse, error)
	// Deprecated: Do not use.
	// CreatePromoter is a method to create a promoter
	CreatePromoter(context.Context, *MsgCreatePromoter) (*MsgCreatePromoterResponse, error)
	// Deprecated: Do not use.
	// CreateCampaign is a method to create a campaign
	CreateCampaign(context.Context, *MsgCreateCampaign) (*MsgCreateCampaignResponse, error)
	// Deprecated: Do not use.
	// UpdateCampaign is a method to update campaign
	UpdateCampaign(context.Context, *MsgUpdateCampaign) (*MsgUpdateCampaignResponse, error)
	// Deprecated: Do not use.
	// WithdrawCampaignFunds is method to withdraw funds from the campaign
	WithdrawFunds(context.Context, *MsgWithdrawFunds) (*MsgWithdrawFundsResponse, error)
	// Deprecated: Do not use.
	// GrantReward is method to allocate rewards
	GrantReward(context.Context, *MsgGrantReward) (*MsgGrantRewardResponse, error)
	// Deprecated: Do not use.
	// UpdateParams defines a governance operation for updating the x/ovm module
	// parameters. The authority is defined in the keeper.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) SetPromoterConf(context.Context, *MsgSetPromoterConf) (*MsgSetPromoterConfResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPromoterConf not implemented")
}
func (UnimplementedMsgServer) CreatePromoter(context.Context, *MsgCreatePromoter) (*MsgCreatePromoterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePromoter not implemented")
}
func (UnimplementedMsgServer) CreateCampaign(context.Context, *MsgCreateCampaign) (*MsgCreateCampaignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCampaign not implemented")
}
func (UnimplementedMsgServer) UpdateCampaign(context.Context, *MsgUpdateCampaign) (*MsgUpdateCampaignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCampaign not implemented")
}
func (UnimplementedMsgServer) WithdrawFunds(context.Context, *MsgWithdrawFunds) (*MsgWithdrawFundsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawFunds not implemented")
}
func (UnimplementedMsgServer) GrantReward(context.Context, *MsgGrantReward) (*MsgGrantRewardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantReward not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_SetPromoterConf_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPromoterConf)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPromoterConf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetPromoterConf_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPromoterConf(ctx, req.(*MsgSetPromoterConf))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreatePromoter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreatePromoter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreatePromoter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreatePromoter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreatePromoter(ctx, req.(*MsgCreatePromoter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateCampaign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateCampaign)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateCampaign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateCampaign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateCampaign(ctx, req.(*MsgCreateCampaign))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateCampaign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateCampaign)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateCampaign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateCampaign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateCampaign(ctx, req.(*MsgUpdateCampaign))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawFunds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawFunds)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawFunds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_WithdrawFunds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawFunds(ctx, req.(*MsgWithdrawFunds))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_GrantReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgGrantReward)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).GrantReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_GrantReward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).GrantReward(ctx, req.(*MsgGrantReward))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sge.legacy.reward.v1beta.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetPromoterConf",
			Handler:    _Msg_SetPromoterConf_Handler,
		},
		{
			MethodName: "CreatePromoter",
			Handler:    _Msg_CreatePromoter_Handler,
		},
		{
			MethodName: "CreateCampaign",
			Handler:    _Msg_CreateCampaign_Handler,
		},
		{
			MethodName: "UpdateCampaign",
			Handler:    _Msg_UpdateCampaign_Handler,
		},
		{
			MethodName: "WithdrawFunds",
			Handler:    _Msg_WithdrawFunds_Handler,
		},
		{
			MethodName: "GrantReward",
			Handler:    _Msg_GrantReward_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sge/legacy/reward/v1beta/tx.proto",
}
