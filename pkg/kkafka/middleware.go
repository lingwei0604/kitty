package kkafka

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Middleware func(h Handler) Handler

func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func ErrorLogMiddleware(logger log.Logger) Middleware {
	return func(h Handler) Handler {
		return HandleFunc(func(ctx context.Context, msg kafka.Message) error {
			err := h.Handle(ctx, msg)
			if err != nil {
				level.Warn(logger).Log("err", err.Error(), "topic", msg.Topic)
				return err
			}
			return nil
		})
	}
}

func TracingConsumerMiddleware(tracer opentracing.Tracer, opName string) Middleware {
	return func(h Handler) Handler {
		return HandleFunc(func(ctx context.Context, msg kafka.Message) error {
			carrier := getCarrier(&msg)
			spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
			span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, opName, opentracing.FollowsFrom(spanContext))
			defer span.Finish()

			ext.SpanKind.Set(span, ext.SpanKindConsumerEnum)
			ext.PeerService.Set(span, "kafka")
			span.SetTag("topic", msg.Topic)
			span.SetTag("partition", msg.Partition)
			span.SetTag("offset", msg.Offset)

			err = h.Handle(ctx, msg)
			if err != nil {
				span.LogKV("error", err.Error())
				ext.Error.Set(span, true)
				// This is user's fault, so we are not returning any error here
				return err
			}
			return nil
		})
	}
}

func TracingProducerMiddleware(tracer opentracing.Tracer, opName string) Middleware {
	return func(h Handler) Handler {
		return HandleFunc(func(ctx context.Context, msg kafka.Message) error {
			span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, opName)
			defer span.Finish()
			ext.SpanKind.Set(span, ext.SpanKindProducerEnum)

			carrier := make(opentracing.TextMapCarrier)
			err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
			if err != nil {
				return errors.Wrap(err, "unable to inject tracing context")
			}

			var header kafka.Header
			for k, v := range carrier {
				header.Key = k
				header.Value = []byte(v)
			}
			msg.Headers = append(msg.Headers, header)
			span.LogKV("message", string(msg.Value))
			err = h.Handle(ctx, msg)
			if err != nil {
				span.LogKV("error", err.Error())
				ext.Error.Set(span, true)
				// This is user's fault, so we are not returning any error here
				return err
			}

			return nil
		})
	}
}

func getCarrier(msg *kafka.Message) opentracing.TextMapCarrier {

	var mapCarrier = make(opentracing.TextMapCarrier)
	if msg.Headers != nil {
		for _, v := range msg.Headers {
			mapCarrier[v.Key] = string(v.Value)
		}
	}
	return mapCarrier
}

// ContextToKafka returns an kafka RequestResponseFunc that injects an OpenTracing Span
// found in `ctx` into the http headers. If no such Span can be found, the
// RequestFunc is a noop.
func ContextToKafka(tracer opentracing.Tracer, logger log.Logger) RequestResponseFunc {
	return func(ctx context.Context, msg *kafka.Message) context.Context {
		// Try to find a Span in the Context.
		if span := opentracing.SpanFromContext(ctx); span != nil {
			// Add standard OpenTracing tags.
			ext.SpanKind.Set(span, ext.SpanKindProducerEnum)

			carrier := make(opentracing.TextMapCarrier)
			err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
			if err != nil {
				level.Warn(logger).Log("err", fmt.Sprintf("unable to inject tracing context: %s", err.Error()))
			}

			var header kafka.Header
			for k, v := range carrier {
				header.Key = k
				header.Value = []byte(v)
			}
			msg.Headers = append(msg.Headers, header)
		}
		return ctx
	}
}

// KafkaToContext returns an http RequestFunc that tries to join with an
// OpenTracing trace found in `req` and starts a new Span called
// `operationName` accordingly. If no trace could be found in `req`, the Span
// will be a trace root. The Span is incorporated in the returned Context and
// can be retrieved with opentracing.SpanFromContext(ctx).
func KafkaToContext(tracer opentracing.Tracer, operationName string, logger log.Logger) RequestResponseFunc {
	return func(ctx context.Context, msg *kafka.Message) context.Context {

		carrier := getCarrier(msg)
		spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			logger.Log("err", err)
		}
		span := tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			level.Warn(logger).Log("err", fmt.Sprintf("unable to extract tracing context: %s", err.Error()))
		}
		ext.SpanKind.Set(span, ext.SpanKindConsumerEnum)
		ext.PeerService.Set(span, "kafka")
		span.SetTag("topic", msg.Topic)
		span.SetTag("partition", msg.Partition)
		span.SetTag("offset", msg.Offset)

		return opentracing.ContextWithSpan(ctx, span)
	}
}

func ErrHandler(logger log.Logger) transport.ErrorHandler {
	return transport.ErrorHandlerFunc(func(ctx context.Context, err error) {
		level.Warn(logger).Log("err", err.Error())
	})
}
