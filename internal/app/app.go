package app

import (
	"car-pooling-challenge/internal/domain"
	"car-pooling-challenge/pkg/config"
	"car-pooling-challenge/pkg/logger"
	"car-pooling-challenge/pkg/router"
	"car-pooling-challenge/pkg/server"
	"log"
	"net/http"
	"net/http/pprof"
)

type App struct {
	logger  *log.Logger
	config  *config.Config
	server  *server.Server
	pooling *domain.Pooling
}

func New(cfg *config.Config) *App {

	var router = router.New()
	var pooling = domain.NewPooling(cfg)

	router.Register("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, pooling)
	})
	router.Register("/cars", func(w http.ResponseWriter, r *http.Request) {
		carsHandler(w, r, pooling)
	})
	router.Register("/journey", func(w http.ResponseWriter, r *http.Request) {
		journeyHandler(w, r, pooling)
	})
	router.Register("/dropoff", func(w http.ResponseWriter, r *http.Request) {
		dropoffHandler(w, r, pooling)
	})
	router.Register("/locate", func(w http.ResponseWriter, r *http.Request) {
		locateHandler(w, r, pooling)
	})

	if cfg.ServerDebug {
		router.Register("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Register("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Register("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Register("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		router.Register("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	}
	// visualize the heap profiler
	// > go tool pprof http://localhost:9091/debug/pprof/profile
	// > (pprof) web

	return &App{
		logger:  logger.NewLogger("App"),
		config:  cfg,
		server:  server.New(cfg, router),
		pooling: pooling,
	}
}

func (a *App) Run() {
	a.pooling.Run()
	a.server.ListenAndServeWithGracefulShutdown()
	a.pooling.CancelContext()
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) GetLogger() *log.Logger {
	return a.logger
}

func (a *App) GetServer() *server.Server {
	return a.server
}
