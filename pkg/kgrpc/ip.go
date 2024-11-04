package kgrpc

import (
	"context"

	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func IpToContext() grpctransport.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		remote, _ := peer.FromContext(ctx)
		return context.WithValue(ctx, contract.IpKey, remote.Addr.String())
	}
}
