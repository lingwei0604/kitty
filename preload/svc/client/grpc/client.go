// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: e12fd89529
// Version Date: 2021-03-04T06:59:01Z

// Package grpc provides a gRPC client for the Preload service.
package grpc

import (
	"context"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	"github.com/lingwei0604/kitty/preload/svc"
	pb "github.com/lingwei0604/kitty/proto"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...grpctransport.ClientOption) (svc.Endpoints, error) {
	var listinfoEndpoint endpoint.Endpoint
	{
		listinfoEndpoint = grpctransport.NewClient(
			conn,
			"preload.v1.Preload",
			"ListInfo",
			EncodeGRPCListInfoRequest,
			DecodeGRPCListInfoResponse,
			pb.PreloadResp{},
			options...,
		).Endpoint()
	}

	return svc.Endpoints{
		ListInfoEndpoint: listinfoEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCListInfoResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC listinfo reply to a user-domain listinfo response. Primarily useful in a client.
func DecodeGRPCListInfoResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.PreloadResp)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCListInfoRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain listinfo request to a gRPC listinfo request. Primarily useful in a client.
func EncodeGRPCListInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.PreloadReq)
	return req, nil
}
