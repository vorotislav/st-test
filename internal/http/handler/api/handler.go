// Package api описывает обработчик для запросов api.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	httpErr "st-test/internal/http/handler/handlererrors"
	"st-test/internal/http/handler/responder"
	"st-test/internal/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	// expiresHeader заголовок для времени жизни объекта.
	expiresHeader = "X-EXPIRES"
)

// Storage описывает методы хранилища для сохранения и получения объектов.
//
//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=Storage --with-expecter=true
type Storage interface {
	SaveObject(ctx context.Context, item models.Item) (int, error)
	GetObject(ctx context.Context, id int) (models.Item, error)
}

// Handler http-обработчик запросов.
type Handler struct {
	log   *zap.Logger
	store Storage
}

// NewHandler конструктор для Handler.
func NewHandler(log *zap.Logger, store Storage) *Handler {
	return &Handler{
		log:   log.Named("object handler"),
		store: store,
	}
}

// AddObject метод обработки PUT запросов.
func (h *Handler) AddObject(w http.ResponseWriter, r *http.Request) {
	// получаем ID объекта из пути запроса
	objectID, err := strconv.Atoi(chi.URLParam(r, "objectID"))
	if err != nil {
		h.log.Error("failed get object id", zap.Error(err))

		responder.JSON(w, httpErr.NewInvalidInput("failed get object id", err.Error()))

		return
	}

	// вычитываем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("failed read body", zap.Error(err))

		responder.JSON(w, httpErr.NewInvalidInput("failed read body", err.Error()))

		return
	}

	defer r.Body.Close()

	// проверяем на валидность json
	err = validate(body)
	if err != nil {
		h.log.Error("failed check body on json", zap.Error(err))

		responder.JSON(w, httpErr.NewInvalidInput("failed check body on json", err.Error()))

		return
	}

	var duration time.Duration

	// проверяем заголовок с описанием времени жизни объекта
	if expires := r.Header.Get(expiresHeader); expires != "" {
		duration, err = time.ParseDuration(expires)
		if err != nil {
			h.log.Error("failed parse expires duration", zap.Error(err))
		}
	}

	// сохраняем объект в хранилище
	resObjectID, err := h.store.SaveObject(r.Context(), models.Item{
		ID:      objectID,
		Body:    body,
		Expires: duration,
	})
	if err != nil {
		h.log.Error("failed save object", zap.Error(err))

		responder.JSON(w, httpErr.NewInternalError("failed save object", err.Error()))

		return
	}

	// если из хранилища вернулся id = 0, значит объект уже был в хранилище и только обновили информацию
	if resObjectID == 0 {
		h.log.Info("update object successful")

		w.WriteHeader(http.StatusNoContent)

		return
	}

	h.log.Info("save object successful")

	w.WriteHeader(http.StatusOK)
}

// Object возвращает объект из хранилища.
func (h *Handler) Object(w http.ResponseWriter, r *http.Request) {
	// получаем ID объекта из пути запроса
	objectID, err := strconv.Atoi(chi.URLParam(r, "objectID"))
	if err != nil {
		h.log.Error("failed get object id", zap.Error(err))

		responder.JSON(w, httpErr.NewInvalidInput("failed get object id", err.Error()))

		return
	}

	// получаем объект из хранилища.
	item, err := h.store.GetObject(r.Context(), objectID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			responder.JSON(w, httpErr.NewNotFoundError("failed get object"))

			return
		}

		responder.JSON(w, httpErr.NewInternalError("failed get object", err.Error()))

		return
	}

	responder.JSON(w, item)
}

func validate(raw []byte) error {
	var js json.RawMessage
	return json.Unmarshal(raw, &js) //nolint:wrapcheck,nlreturn
}
