package module

import (
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/kerr"
	"github.com/lingwei0604/kitty/pkg/khttp"
	"github.com/lingwei0604/kitty/preload/svc"
	kitty "github.com/lingwei0604/kitty/proto"
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
