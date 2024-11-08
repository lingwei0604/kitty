// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 0035b7fc88
// Version Date: 2022-11-02T08:53:09Z

// Package grpc provides a gRPC client for the App service.
package grpc

import (
	"context"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	"github.com/lingwei0604/kitty/app/svc"
	pb "github.com/lingwei0604/kitty/proto"
)

// app.v2
// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...grpctransport.ClientOption) (svc.Endpoints, error) {
	var preregisterEndpoint endpoint.Endpoint
	{
		preregisterEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"PreRegister",
			EncodeGRPCPreRegisterRequest,
			DecodeGRPCPreRegisterResponse,
			pb.PreRegisterReply{},
			options...,
		).Endpoint()
	}

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"Login",
			EncodeGRPCLoginRequest,
			DecodeGRPCLoginResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var bindwechatEndpoint endpoint.Endpoint
	{
		bindwechatEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"BindWechat",
			EncodeGRPCBindWechatRequest,
			DecodeGRPCBindWechatResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var getcodeEndpoint endpoint.Endpoint
	{
		getcodeEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"GetCode",
			EncodeGRPCGetCodeRequest,
			DecodeGRPCGetCodeResponse,
			pb.GenericReply{},
			options...,
		).Endpoint()
	}

	var getinfoEndpoint endpoint.Endpoint
	{
		getinfoEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"GetInfo",
			EncodeGRPCGetInfoRequest,
			DecodeGRPCGetInfoResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var getinfobatchEndpoint endpoint.Endpoint
	{
		getinfobatchEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"GetInfoBatch",
			EncodeGRPCGetInfoBatchRequest,
			DecodeGRPCGetInfoBatchResponse,
			pb.UserInfoBatchReply{},
			options...,
		).Endpoint()
	}

	var updateinfoEndpoint endpoint.Endpoint
	{
		updateinfoEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"UpdateInfo",
			EncodeGRPCUpdateInfoRequest,
			DecodeGRPCUpdateInfoResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var bindEndpoint endpoint.Endpoint
	{
		bindEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"Bind",
			EncodeGRPCBindRequest,
			DecodeGRPCBindResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var bindadEndpoint endpoint.Endpoint
	{
		bindadEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"BindAd",
			EncodeGRPCBindAdRequest,
			DecodeGRPCBindAdResponse,
			pb.GenericReply{},
			options...,
		).Endpoint()
	}

	var unbindEndpoint endpoint.Endpoint
	{
		unbindEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"Unbind",
			EncodeGRPCUnbindRequest,
			DecodeGRPCUnbindResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var refreshEndpoint endpoint.Endpoint
	{
		refreshEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"Refresh",
			EncodeGRPCRefreshRequest,
			DecodeGRPCRefreshResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var softdeleteEndpoint endpoint.Endpoint
	{
		softdeleteEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"SoftDelete",
			EncodeGRPCSoftDeleteRequest,
			DecodeGRPCSoftDeleteResponse,
			pb.UserInfoReply{},
			options...,
		).Endpoint()
	}

	var devicelookupEndpoint endpoint.Endpoint
	{
		devicelookupEndpoint = grpctransport.NewClient(
			conn,
			"app.v2.App",
			"DeviceLookup",
			EncodeGRPCDeviceLookupRequest,
			DecodeGRPCDeviceLookupResponse,
			pb.DeviceLookupReply{},
			options...,
		).Endpoint()
	}

	return svc.Endpoints{
		PreRegisterEndpoint:  preregisterEndpoint,
		LoginEndpoint:        loginEndpoint,
		BindWechatEndpoint:   bindwechatEndpoint,
		GetCodeEndpoint:      getcodeEndpoint,
		GetInfoEndpoint:      getinfoEndpoint,
		GetInfoBatchEndpoint: getinfobatchEndpoint,
		UpdateInfoEndpoint:   updateinfoEndpoint,
		BindEndpoint:         bindEndpoint,
		BindAdEndpoint:       bindadEndpoint,
		UnbindEndpoint:       unbindEndpoint,
		RefreshEndpoint:      refreshEndpoint,
		SoftDeleteEndpoint:   softdeleteEndpoint,
		DeviceLookupEndpoint: devicelookupEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCPreRegisterResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC preregister reply to a user-domain preregister response. Primarily useful in a client.
func DecodeGRPCPreRegisterResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.PreRegisterReply)
	return reply, nil
}

// DecodeGRPCLoginResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC login reply to a user-domain login response. Primarily useful in a client.
func DecodeGRPCLoginResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCBindWechatResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC bindwechat reply to a user-domain bindwechat response. Primarily useful in a client.
func DecodeGRPCBindWechatResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCGetCodeResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getcode reply to a user-domain getcode response. Primarily useful in a client.
func DecodeGRPCGetCodeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GenericReply)
	return reply, nil
}

// DecodeGRPCGetInfoResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getinfo reply to a user-domain getinfo response. Primarily useful in a client.
func DecodeGRPCGetInfoResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCGetInfoBatchResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getinfobatch reply to a user-domain getinfobatch response. Primarily useful in a client.
func DecodeGRPCGetInfoBatchResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoBatchReply)
	return reply, nil
}

// DecodeGRPCUpdateInfoResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC updateinfo reply to a user-domain updateinfo response. Primarily useful in a client.
func DecodeGRPCUpdateInfoResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCBindResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC bind reply to a user-domain bind response. Primarily useful in a client.
func DecodeGRPCBindResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCBindAdResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC bindad reply to a user-domain bindad response. Primarily useful in a client.
func DecodeGRPCBindAdResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GenericReply)
	return reply, nil
}

// DecodeGRPCUnbindResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC unbind reply to a user-domain unbind response. Primarily useful in a client.
func DecodeGRPCUnbindResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCRefreshResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC refresh reply to a user-domain refresh response. Primarily useful in a client.
func DecodeGRPCRefreshResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCSoftDeleteResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC softdelete reply to a user-domain softdelete response. Primarily useful in a client.
func DecodeGRPCSoftDeleteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCDeviceLookupResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC devicelookup reply to a user-domain devicelookup response. Primarily useful in a client.
func DecodeGRPCDeviceLookupResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.DeviceLookupReply)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCPreRegisterRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain preregister request to a gRPC preregister request. Primarily useful in a client.
func EncodeGRPCPreRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.PreRegisterRequest)
	return req, nil
}

// EncodeGRPCLoginRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain login request to a gRPC login request. Primarily useful in a client.
func EncodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserLoginRequest)
	return req, nil
}

// EncodeGRPCBindWechatRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain bindwechat request to a gRPC bindwechat request. Primarily useful in a client.
func EncodeGRPCBindWechatRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.BindWechatRequest)
	return req, nil
}

// EncodeGRPCGetCodeRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getcode request to a gRPC getcode request. Primarily useful in a client.
func EncodeGRPCGetCodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetCodeRequest)
	return req, nil
}

// EncodeGRPCGetInfoRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getinfo request to a gRPC getinfo request. Primarily useful in a client.
func EncodeGRPCGetInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserInfoRequest)
	return req, nil
}

// EncodeGRPCGetInfoBatchRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getinfobatch request to a gRPC getinfobatch request. Primarily useful in a client.
func EncodeGRPCGetInfoBatchRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserInfoBatchRequest)
	return req, nil
}

// EncodeGRPCUpdateInfoRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain updateinfo request to a gRPC updateinfo request. Primarily useful in a client.
func EncodeGRPCUpdateInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserInfoUpdateRequest)
	return req, nil
}

// EncodeGRPCBindRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain bind request to a gRPC bind request. Primarily useful in a client.
func EncodeGRPCBindRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserBindRequest)
	return req, nil
}

// EncodeGRPCBindAdRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain bindad request to a gRPC bindad request. Primarily useful in a client.
func EncodeGRPCBindAdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserBindAdRequest)
	return req, nil
}

// EncodeGRPCUnbindRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain unbind request to a gRPC unbind request. Primarily useful in a client.
func EncodeGRPCUnbindRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserUnbindRequest)
	return req, nil
}

// EncodeGRPCRefreshRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain refresh request to a gRPC refresh request. Primarily useful in a client.
func EncodeGRPCRefreshRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserRefreshRequest)
	return req, nil
}

// EncodeGRPCSoftDeleteRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain softdelete request to a gRPC softdelete request. Primarily useful in a client.
func EncodeGRPCSoftDeleteRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserSoftDeleteRequest)
	return req, nil
}

// EncodeGRPCDeviceLookupRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain devicelookup request to a gRPC devicelookup request. Primarily useful in a client.
func EncodeGRPCDeviceLookupRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.DeviceLookupRequest)
	return req, nil
}
