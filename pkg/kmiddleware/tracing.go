package kmiddleware

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/lingwei0604/kitty/pkg/config"
	"github.com/lingwei0604/kitty/pkg/contract"
	stdtracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var openAllTrace = false

func init() {
	if os.Getenv("OpenAllTrace") != "" {
		openAllTrace = true
	}
}

type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

func NewTraceServerMiddleware(tracer stdtracing.Tracer, env string) LabeledMiddleware {
	return func(s string, endpoint endpoint.Endpoint) endpoint.Endpoint {
		name := fmt.Sprintf("%s(%s)", s, env)
		return TraceConsumer(tracer, name)(endpoint)
	}
}

func NewClientServerMiddleware(tracer stdtracing.Tracer, env string) LabeledMiddleware {
	return func(s string, endpoint endpoint.Endpoint) endpoint.Endpoint {
		name := fmt.Sprintf("%s(%s)", s, env)
		return TraceProducer(tracer, name)(endpoint)
	}
}

// TraceConsumer returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
//
// If `ctx` already has a Span, it is re-used and the operation name is
// overwritten. If `ctx` does not yet have a Span, one is created here.
func TraceConsumer(tracer stdtracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			serverSpan := stdtracing.SpanFromContext(ctx)
			if serverSpan == nil {
				// All we can do is create a new root span.
				serverSpan = tracer.StartSpan(operationName)
			} else {
				serverSpan.SetOperationName(operationName)
			}
			defer serverSpan.Finish()

			if openAllTrace {
				ext.SamplingPriority.Set(serverSpan, 1)
			} else {
				if ip, ok := ctx.Value(contract.IpKey).(string); ok && ip == "43.228.38.242" {
					ext.SamplingPriority.Set(serverSpan, 1)
				}
			}

			ext.SpanKindConsumer.Set(serverSpan)
			tenant := config.GetTenant(ctx)
			if tenant.UserId != 0 {
				serverSpan.SetTag("package.name", tenant.PackageName)
				serverSpan.SetTag("suuid", tenant.Suuid)
				serverSpan.SetTag("user.id", tenant.UserId)
			}
			ctx = stdtracing.ContextWithSpan(ctx, serverSpan)
			resp, err := next(ctx, request)
			if err != nil {
				ext.Error.Set(serverSpan, true)
				serverSpan.LogKV("error", err.Error())
			}
			serverSpan.LogKV("request", request)
			serverSpan.LogKV("response", resp)
			return resp, err
		}
	}
}

// TraceProducer returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
func TraceProducer(tracer stdtracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var clientSpan stdtracing.Span
			if parentSpan := stdtracing.SpanFromContext(ctx); parentSpan != nil {
				clientSpan = tracer.StartSpan(
					operationName,
					stdtracing.ChildOf(parentSpan.Context()),
				)
			} else {
				clientSpan = tracer.StartSpan(operationName)
			}
			defer clientSpan.Finish()
			ext.SpanKindConsumer.Set(clientSpan)
			ctx = stdtracing.ContextWithSpan(ctx, clientSpan)
			resp, err := next(ctx, request)
			if err != nil {
				ext.Error.Set(clientSpan, true)
				clientSpan.LogKV("error", err.Error())
			}
			clientSpan.LogKV("request", request)
			clientSpan.LogKV("response", resp)
			return resp, err
		}
	}
}
