package internalhttp

import (
	"net/http"

	mainhandler "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/infra/http/main"
)

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(mainhandler.MainHandler()))
	return mux
}
