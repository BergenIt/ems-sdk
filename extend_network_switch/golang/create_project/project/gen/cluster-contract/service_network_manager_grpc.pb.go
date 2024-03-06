// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: service_network_manager.proto

package cluster_contract

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

// NetworkManagerClient is the client API for NetworkManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkManagerClient interface {
	// процедура для сохранения настроек коммутатора
	CreateConfig(ctx context.Context, in *CreateNetworkConfigRequest, opts ...grpc.CallOption) (*CreateNetworkConfigResponse, error)
}

type networkManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkManagerClient(cc grpc.ClientConnInterface) NetworkManagerClient {
	return &networkManagerClient{cc}
}

func (c *networkManagerClient) CreateConfig(ctx context.Context, in *CreateNetworkConfigRequest, opts ...grpc.CallOption) (*CreateNetworkConfigResponse, error) {
	out := new(CreateNetworkConfigResponse)
	err := c.cc.Invoke(ctx, "/tool_cluster.v4.NetworkManager/CreateConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkManagerServer is the server API for NetworkManager service.
// All implementations must embed UnimplementedNetworkManagerServer
// for forward compatibility
type NetworkManagerServer interface {
	// процедура для сохранения настроек коммутатора
	CreateConfig(context.Context, *CreateNetworkConfigRequest) (*CreateNetworkConfigResponse, error)
	mustEmbedUnimplementedNetworkManagerServer()
}

// UnimplementedNetworkManagerServer must be embedded to have forward compatible implementations.
type UnimplementedNetworkManagerServer struct {
}

func (UnimplementedNetworkManagerServer) CreateConfig(context.Context, *CreateNetworkConfigRequest) (*CreateNetworkConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateConfig not implemented")
}
func (UnimplementedNetworkManagerServer) mustEmbedUnimplementedNetworkManagerServer() {}

// UnsafeNetworkManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkManagerServer will
// result in compilation errors.
type UnsafeNetworkManagerServer interface {
	mustEmbedUnimplementedNetworkManagerServer()
}

func RegisterNetworkManagerServer(s grpc.ServiceRegistrar, srv NetworkManagerServer) {
	s.RegisterService(&NetworkManager_ServiceDesc, srv)
}

func _NetworkManager_CreateConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNetworkConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkManagerServer).CreateConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tool_cluster.v4.NetworkManager/CreateConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkManagerServer).CreateConfig(ctx, req.(*CreateNetworkConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NetworkManager_ServiceDesc is the grpc.ServiceDesc for NetworkManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tool_cluster.v4.NetworkManager",
	HandlerType: (*NetworkManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateConfig",
			Handler:    _NetworkManager_CreateConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_network_manager.proto",
}