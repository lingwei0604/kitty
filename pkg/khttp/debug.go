package khttp

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
)

func Debug(router *mux.Router) {
	m := mux.NewRouter()
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)
	m.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
	router.PathPrefix("/debug/").Handler(m)
}
