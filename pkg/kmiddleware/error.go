package kmiddleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/lingwei0604/kitty/app/msg"
	"github.com/lingwei0604/kitty/pkg/kerr"
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
