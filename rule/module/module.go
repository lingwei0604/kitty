package module

import (
	"context"
	"git.yingzhongshare.com/mkt/kitty/pkg/kerr"
	"git.yingzhongshare.com/mkt/kitty/pkg/khttp"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/tracing/opentracing"
	"net/http"
	"strconv"

	"git.yingzhongshare.com/mkt/kitty/pkg/config"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/rule/service"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type Module struct {
	tracer     stdopentracing.Tracer
	repository service.Repository
	endpoints  Endpoints
	logger     log.Logger
	close      func()
}

func New(moduleConfig contract.ConfigReader, logger log.Logger) *Module {
	module, cleanup, err := injectModule(moduleConfig, logger)
	if err != nil {
		panic(err)
	}
	module.close = cleanup
	return module
}

func (m *Module) ProvideHttp(router *mux.Router) {
	//router.PathPrefix("/rule/").Handler(mw.Security(http.StripPrefix("/rule", MakeHTTPHandler(m.endpoints,
	//	httptransport.ServerBefore(
	//		opentracing.HTTPToContext(m.tracer, "app", m.logger),
	//		jwt.HTTPToContext(),
	//		khttp.IpToContext(),
	//		tenantToContext(),
	//	),
	//	httptransport.ServerErrorEncoder(kerr.ErrorEncoder)))))
	router.PathPrefix("/rule/").Handler(http.StripPrefix("/rule", MakeHTTPHandler(m.endpoints,
		httptransport.ServerBefore(
			opentracing.HTTPToContext(m.tracer, "app", m.logger),
			jwt.HTTPToContext(),
			khttp.IpToContext(),
			tenantToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder))))
}

func tenantToContext() httptransport.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		query := request.URL.Query()
		userID, _ := strconv.ParseUint(query.Get("user_id"), 10, 64)
		return context.WithValue(ctx, config.TenantKey, &config.Tenant{
			Channel:     query.Get("channel"),
			VersionCode: query.Get("version_code"),
			UserId:      userID,
			Imei:        query.Get("imei"),
			Oaid:        query.Get("oaid"),
			Suuid:       query.Get("suuid"),
		})
	}
}

func (m *Module) ProvideCloser() {
	m.close()
}

func (m *Module) ProvideRunGroup(group *run.Group) {
	ctx, cancel := context.WithCancel(context.Background())
	group.Add(func() error {
		return m.repository.WatchConfigUpdate(ctx)
	}, func(err error) {
		cancel()
	})
}

func (m *Module) GetRepository() service.Repository {
	return m.repository
}
