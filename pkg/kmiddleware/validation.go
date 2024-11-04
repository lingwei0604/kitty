package kmiddleware

import (
	"context"

	"git.yingzhongshare.com/mkt/kitty/app/msg"
	"git.yingzhongshare.com/mkt/kitty/pkg/kerr"
	"github.com/go-kit/kit/endpoint"
)

type validator interface {
	Validate() error
}

func NewValidationMiddleware() endpoint.Middleware {
	return func(in endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			if t, ok := req.(validator); ok {
				err = t.Validate()
				if err != nil {
					return nil, kerr.InvalidArgumentErr(err, msg.InvalidParams)
				}
			}
			resp, err = in(ctx, req)
			return
		}
	}
}
