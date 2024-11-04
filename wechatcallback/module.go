package wechatcallback

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/lingwei0604/kitty/pkg/contract"
)

type Module struct {
	Handler *Handler
	cleanup func()
}

func New(reader contract.ConfigReader, logger log.Logger) *Module {
	appModule, cleanup, err := injectModule(reader, logger)
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.HandleFunc("/wechat/echo", a.Handler.Echo).Methods("GET")
	router.HandleFunc("/wechat/unbind/{packageName}", a.Handler.UnbindWechatUser).Methods("POST")
}

func (a *Module) ProvideCloser() {
	a.cleanup()
}
