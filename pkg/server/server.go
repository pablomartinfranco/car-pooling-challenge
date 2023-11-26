package server

import (
	"car-pooling-challenge/pkg/config"
	"car-pooling-challenge/pkg/logger"
	"car-pooling-challenge/pkg/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	logger   *log.Logger
	config   *config.Config
	router   *router.Router
	instance *http.Server
}

func New(config *config.Config, router *router.Router) *Server {

	var mux = http.NewServeMux()

	for _, route := range router.GetRoutes() {
		mux.HandleFunc(route.Pattern, route.HandlerFunc)
	}

	var address = config.ServerHost + ":" + config.ServerPort

	var server = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &Server{
		logger:   logger.NewLogger("Server"),
		config:   config,
		router:   router,
		instance: server,
	}
}

func (s *Server) GetInstance() *http.Server {
	return s.instance
}

func (s *Server) ListenAndServe() {
	if err := s.instance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Printf("Oops... Server is not starting! Reason: %v", err)
		panic(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("Shutting down...")
	return s.instance.Shutdown(ctx)
}

func (s *Server) ListenAndServeWithGracefulShutdown() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		var sigint = make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint // Wait os.Interrupt signal.

		// Create a context with a timeout to give existing connections time to finish.
		var duration = time.Duration(s.config.ServerReadTimeout) * time.Second
		_ctx, _cancel := context.WithTimeout(context.Background(), duration)
		defer _cancel()

		// Shutdown the server and block until all connections are closed.
		if err := s.Shutdown(_ctx); err != nil {
			s.logger.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		// Cancel the main context to stop the server gracefully.
		cancel()
	}()

	go func() {
		s.ListenAndServe()
	}()
	s.logger.Printf("Server started on %s:%s", s.config.ServerHost, s.config.ServerPort)

	// Block until a signal is received to exit.
	<-ctx.Done()
	s.logger.Println("Server stopped gracefully.")
}
