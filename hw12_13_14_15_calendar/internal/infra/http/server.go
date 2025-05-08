package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/app/provider"
)

type Server struct {
	httpServer *http.Server
}

var ErrIncorrectPort = errors.New("incorrect http port")

func NewServer() *Server {
	conf := provider.CalendarContainer.GetConfig()

	port, err := strconv.Atoi(conf.HTTP.Port)
	if err != nil || port <= 0 || port > 65535 {
		provider.CalendarContainer.GetLogger().Error(context.TODO(), "Fatal error: "+ErrIncorrectPort.Error())
		os.Exit(1)
	}

	addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	router := newRouter()

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
