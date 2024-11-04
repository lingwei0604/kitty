package handlers

import (
	"github.com/go-kit/kit/log"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/wechat"
	pb "github.com/lingwei0604/kitty/proto"
)

func NewAppService(
	conf contract.ConfigReader,
	log log.Logger,
	ur UserRepository,
	cr CodeRepository,
	fr FileRepository,
	sender contract.SmsSender,
	wechat wechat.Wechater,
	dispatcher contract.Dispatcher,
	allocator PreAllocator,
	userCache UserCache,
) appService {
	return appService{
		conf:         conf,
		logger:       log,
		ur:           ur,
		cr:           cr,
		sender:       sender,
		wechat:       wechat,
		fr:           fr,
		dispatcher:   dispatcher,
		preAllocator: allocator,
		userCache:    userCache,
	}
}

type ServerMiddleware func(server pb.AppServer) pb.AppServer

func NewInputEnrichMiddleware() ServerMiddleware {
	return func(server pb.AppServer) pb.AppServer {
		return &InputEnrichedAppService{AppServer: server}
	}
}

func Chain(outer ServerMiddleware, others ...ServerMiddleware) ServerMiddleware {
	return func(next pb.AppServer) pb.AppServer {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func ProvideAppServer(service appService) pb.AppServer {
	return Chain(
		NewInputEnrichMiddleware(),
	)(service)
}
