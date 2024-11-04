package kmiddleware

import (
	"context"
	"errors"
	"fmt"

	"git.yingzhongshare.com/mkt/kitty/app/msg"
	"git.yingzhongshare.com/mkt/kitty/pkg/kerr"
	"github.com/go-kit/kit/endpoint"
)

func NewErrorMarshallerMiddleware(handlePanic bool) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func() {
				if !handlePanic {
					return
				}
				if er := recover(); er != nil {
					err = kerr.InternalErr(fmt.Errorf("panic: %s", er), msg.ServerBug)
				}
			}()
			response, err = e(ctx, request)
			if err != nil {
				var serverError kerr.ServerError
				if !errors.As(err, &serverError) {
					serverError = kerr.UnknownErr(err)
				}
				// Brings kerr.SeverError to the uppermost level
				return response, serverError
			}

			return response, nil
		}
	}
}
