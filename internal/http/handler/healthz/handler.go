// Package healthz описывает обработки для запросов проверки здоровья сервиса.
package healthz

import (
	"net/http"

	"go.uber.org/zap"
)

// store описывает метод проверки доступности хранилища.
//
//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=store --exported --with-expecter=true
type store interface {
	Check() error
}

// Handler http-обработчик запросов.
type Handler struct {
	log   *zap.Logger
	store store
}

// NewHandler конструктор для Handler.
func NewHandler(log *zap.Logger, store store) *Handler {
	return &Handler{
		log:   log,
		store: store,
	}
}

// Liveness метод обработки запроса на получение liveness-пробы.
func (h *Handler) Liveness(w http.ResponseWriter, _ *http.Request) {
	if err := h.store.Check(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}

// Readiness метод обработки запроса на получение readiness-пробы.
func (h *Handler) Readiness(w http.ResponseWriter, _ *http.Request) {
	if err := h.store.Check(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}
