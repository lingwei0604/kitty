package kgrpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/lingwei0604/kitty/pkg/contract"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func IpToContext() grpctransport.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		remote, _ := peer.FromContext(ctx)
		return context.WithValue(ctx, contract.IpKey, remote.Addr.String())
	}
}
