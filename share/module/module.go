package module

import (
	"context"
	"net/http"

	//mw "gitee.com/tagtic/go-middleware/http/middleware"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/lingwei0604/kitty/pkg/config"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/kkafka"
	pb "github.com/lingwei0604/kitty/proto"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

type Module struct {
	appName     contract.AppName
	cleanup     func()
	handler     http.Handler
	grpcServer  GrpcShareServer
	kafkaServer kkafka.Server
}

func New(appModuleConfig contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) *Module {
	appModule, cleanup, err := injectModule(appModuleConfig, logger, dynConf)
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func (a *Module) ProvideCloser() {
	a.cleanup()
}

func (a *Module) ProvideGrpc(server *grpc.Server) {
	pb.RegisterShareServer(server, a.grpcServer)
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/share/v1/").Handler(http.StripPrefix("/share/v1", a.handler))
}

func (a *Module) ProvideRunGroup(group *run.Group) {
	ctx, cancel := context.WithCancel(context.Background())
	group.Add(func() error {
		return a.kafkaServer.Serve(ctx)
	}, func(err error) {
		cancel()
	})

}
