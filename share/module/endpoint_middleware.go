package module

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/kmiddleware"
	"github.com/lingwei0604/kitty/share/svc"
	"github.com/opentracing/opentracing-go"
)

func provideEndpointsMiddleware(l log.Logger, securityConfig *kmiddleware.SecurityConfig, hist metrics.Histogram, tracer opentracing.Tracer, env contract.Env, appName contract.AppName) overallMiddleware {
	return func(in svc.Endpoints) svc.Endpoints {
		in.WrapAllExcept(kmiddleware.NewValidationMiddleware())
		in.WrapAllExcept(kmiddleware.NewLoggingMiddleware(l, env.IsLocal()))
		in.WrapAllLabeledExcept(kmiddleware.NewLabeledMetricsMiddleware(hist, appName.String(), "share"))
		in.WrapAllLabeledExcept(kmiddleware.NewTraceServerMiddleware(tracer, env.String()))
		in.WrapAllExcept(kmiddleware.NewConfigMiddleware())
		in.WrapAllExcept(kmiddleware.NewAuthenticationMiddleware(securityConfig))
		in.WrapAllExcept(kmiddleware.NewErrorMarshallerMiddleware(env.IsProd()))
		return in
	}
}
