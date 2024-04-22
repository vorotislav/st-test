// Package http содержит описание http-сервиса, который скрывает подготовку http-сервера и создание обработчиков.
package http

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"st-test/internal/http/handler/api"
	"st-test/internal/http/handler/healthz"
	"st-test/internal/http/handler/middlewares/apptype"
	"st-test/internal/settings"
	"st-test/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Service описание сервиса. Хранит логгер, http-сервер и настройки.
type Service struct {
	logger   *zap.Logger
	server   *http.Server
	settings *settings.APISettings
}

// NewService получает логгер, настройки и хранилище и создаёт объект Сервис.
func NewService(log *zap.Logger, set *settings.APISettings, store *storage.Store) *Service {
	serLog := log.Named("http-service")

	mux := chi.NewRouter()

	// A good base middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(time.Second))

	mux.Use(apptype.ApplicationType(log))

	// api handlers
	apiHandler := api.NewHandler(log, store)

	mux.Put("/objects"+"/{objectID}", apiHandler.AddObject)
	mux.Get("/objects"+"/{objectID}", apiHandler.Object)

	// metrics handler
	mux.Handle("/metrics", promhttp.Handler())

	// healthchecks handlers
	healtzHandler := healthz.NewHandler(log, store)
	mux.Get("/probes/liveness", healtzHandler.Liveness)
	mux.Get("/probes/readiness", healtzHandler.Readiness)

	s := &http.Server{
		Addr:              net.JoinHostPort(set.Address, strconv.Itoa(set.Port)),
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	return &Service{
		logger:   serLog,
		server:   s,
		settings: set,
	}
}

// Run запускает http-сервер на прослушивание адреса и порта.
func (s *Service) Run() error {
	s.logger.Debug("Running server on", zap.String("address", s.server.Addr))

	return s.server.ListenAndServe() //nolint:wrapcheck
}

// Stop выключает http-сервер.
func (s *Service) Stop(ctx context.Context) error {
	s.logger.Debug("stopping http service")

	return s.server.Shutdown(ctx) //nolint:wrapcheck
}
