package module

import (
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/kerr"
	"git.yingzhongshare.com/mkt/kitty/pkg/khttp"
	"git.yingzhongshare.com/mkt/kitty/preload/svc"
	kitty "git.yingzhongshare.com/mkt/kitty/proto"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type Module struct {
	cleanup   func()
	endpoints svc.Endpoints
	logger    log.Logger
	conf      contract.ConfigReader
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/preload/v1").Handler(http.StripPrefix("/preload/v1", svc.MakeHTTPHandler(a.endpoints,
		httptransport.ServerBefore(
			khttp.IpToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	)))
}

func New(conf contract.ConfigReader, logger log.Logger) *Module {
	module, cleanup, err := injectModule(conf, logger)
	if err != nil {
		panic(err)
	}
	module.cleanup = cleanup
	return module
}

func provideModule(conf contract.ConfigReader, logger log.Logger, server kitty.PreloadServer) *Module {
	return &Module{
		conf:      conf,
		endpoints: svc.NewEndpoints(server),
		logger:    logger,
	}
}
