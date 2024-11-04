package kmiddleware

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/DoNewsCode/core/srvhttp"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	his         metrics.Histogram
	initMetrics sync.Once
)

func ProvideHistogramMetrics() metrics.Histogram {
	initMetrics.Do(func() {
		vec := stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Total time spent serving requests.",
		}, []string{"module", "service", "route", "status"})
		stdprometheus.DefaultRegisterer.MustRegister(vec)
		his = prometheus.NewHistogram(vec)
	})
	return his
}

func NewLabeledMetricsMiddleware(his metrics.Histogram, module, service string) LabeledMiddleware {
	return func(name string, e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				status := 200
				if sc, ok := err.(srvhttp.StatusCoder); ok {
					status = sc.StatusCode()
				}
				his.With("module", module, "route", name, "service", service, "status", strconv.Itoa(status)).
					Observe(time.Since(begin).Seconds())
			}(time.Now())
			return e(ctx, request)
		}
	}
}

func NewMetricsMiddleware(his metrics.Histogram, module, service, method string) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				status := 200
				if sc, ok := err.(srvhttp.StatusCoder); ok {
					status = sc.StatusCode()
				}
				his.With("module", module, "route", method, "service", service,
					"status", strconv.Itoa(status)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return e(ctx, request)
		}
	}
}
