package ots3

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/lingwei0604/kitty/pkg/contract"
	"net/http"
)

type Module struct {
	handler http.Handler
}

func (m Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/upload").Handler(m.handler)
}

func New(conf contract.ConfigReader, logger log.Logger) *Module {
	return injectModule(conf, logger)
}
