package server

import (
	"context"
	"github.com/ervitis/backend-challenge/clientrest/domain"
	"github.com/ervitis/logme"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type (
	Server struct {
		address, port string
		router        domain.IRouter
		log           logme.Loggerme
		srv           *http.Server
	}

	Options func(*Server)
)

func WithPort(port string) Options {
	return func(server *Server) {
		server.port = port
	}
}

func WithAddress(address string) Options {
	return func(server *Server) {
		server.address = address
	}
}

func WithLogger(logger logme.Loggerme) Options {
	return func(server *Server) {
		server.log = logger
	}
}

func WithRouter(r domain.IRouter) Options {
	return func(server *Server) {
		server.router = r
	}
}

func defaultOptions() *Server {
	return &Server{address: "127.0.0.1", port: "8080"}
}

func (s *Server) Listen() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	errChannel := make(chan error)

	go func() {
		s.log.L().Infof("server running on %s", s.getStringConnection())
		if err := s.srv.ListenAndServe(); err != nil {
			errChannel <- err
		}
	}()

	go func() {
		for {
			select {
			case err := <-errChannel:
				s.log.L().Fatalf("error received from handler: %s", err.Error())
				signals <- os.Interrupt
			}
		}
	}()

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.L().Fatalf("error shutdown server: %s", err.Error())
	}
	os.Exit(0)
}

func (s *Server) getStringConnection() string {
	return s.address + ":" + s.port
}

func CreateServer(serverOptions ...Options) *Server {
	opts := defaultOptions()

	for _, opt := range serverOptions {
		opt(opts)
	}

	if opts.router == nil {
		panic("handler router not set with WithRouter")
	}

	timeout := 15 * time.Second
	idleTimeout := 45 * time.Second
	if os.Getenv("DEBUG") != "" {
		timeout = 240 * time.Second
		idleTimeout = 500 * time.Second
	}

	opts.srv = &http.Server{
		Handler: opts.router.GetRouter(),
		Addr: opts.getStringConnection(),
		WriteTimeout: timeout,
		ReadTimeout: timeout,
		IdleTimeout: idleTimeout,
	}

	return opts
}
