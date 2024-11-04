package module

import (
	"context"

	"git.yingzhongshare.com/mkt/kitty/app/svc"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/kmiddleware"
	pb "git.yingzhongshare.com/mkt/kitty/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/opentracing/opentracing-go"
)

// newLoginToBindMiddleware deprecated
func newLoginToBindMiddleware(bind endpoint.Endpoint) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if ctx.Value(jwt.JWTTokenContextKey) == nil {
				return e(ctx, request)
			}
			loginRequest, ok := request.(pb.UserLoginRequest)
			if !ok {
				return e(ctx, request)
			}
			if len(loginRequest.Mobile) <= 0 && len(loginRequest.Wechat) <= 0 {
				return e(ctx, request)
			}
			bindReq := pb.UserBindRequest{
				Mobile: loginRequest.Mobile,
				Code:   loginRequest.Code,
				Wechat: loginRequest.Wechat,
			}
			return bind(ctx, bindReq)
		}
	}
}

func provideEndpointsMiddleware(l log.Logger, securityConfig *kmiddleware.SecurityConfig, hist metrics.Histogram, tracer opentracing.Tracer, env contract.Env, appName contract.AppName) overallMiddleware {
	return func(in svc.Endpoints) svc.Endpoints {
		in.WrapAllExcept(kmiddleware.NewValidationMiddleware())
		in.WrapAllExcept(kmiddleware.NewLoggingMiddleware(l, env.IsLocal()))
		in.WrapAllLabeledExcept(kmiddleware.NewLabeledMetricsMiddleware(hist, appName.String(), "user"))
		in.WrapAllLabeledExcept(kmiddleware.NewTraceServerMiddleware(tracer, env.String()))
		in.WrapAllExcept(kmiddleware.NewConfigMiddleware())
		in.GetInfoEndpoint = kmiddleware.NewOptionalAuthenticationMiddleware(securityConfig)(in.GetInfoEndpoint)
		in.WrapAllExcept(kmiddleware.NewAuthenticationMiddleware(securityConfig), "Login", "GetCode", "GetInfo", "PreRegister")
		in.WrapAllExcept(kmiddleware.NewErrorMarshallerMiddleware(env.IsProd()))
		return in
	}
}
