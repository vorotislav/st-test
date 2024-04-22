package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"st-test/internal/http/handler/api/mocks"
	"st-test/internal/models"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	store := mocks.NewStorage(t)
	require.NotNil(t, store)

	h := NewHandler(log, store)
	require.NotNil(t, h)
}

func TestHandler_AddObject(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	cases := []struct {
		name         string
		prepareStore func(store *mocks.Storage)
		giveRequest  func() *http.Request
		checkResult  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "invalid object id format",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "invalid")

				body := []byte(`some body`)
				req, _ := http.NewRequest(http.MethodPut, "foo/bar", bytes.NewBuffer(body))
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Contains(t, rr.Body.String(), "failed get object id")
			},
		},
		{
			name: "invalid object body format",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				body := []byte(`some body`)
				req, _ := http.NewRequest(http.MethodPut, "foo/bar", bytes.NewBuffer(body))
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Contains(t, rr.Body.String(), "failed check body on json")
			},
		},
		{
			name: "store error",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				body := []byte(`{"some":"body"}`)
				req, _ := http.NewRequest(http.MethodPut, "foo/bar", bytes.NewBuffer(body))
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			prepareStore: func(store *mocks.Storage) {
				store.EXPECT().SaveObject(mock.Anything, mock.AnythingOfType("models.Item")).
					Once().
					Return(0, errors.New("some error"))
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "failed save object")
			},
		},
		{
			name: "successful update",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				body := []byte(`{"some":"body"}`)
				req, _ := http.NewRequest(http.MethodPut, "foo/bar", bytes.NewBuffer(body))
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			prepareStore: func(store *mocks.Storage) {
				store.EXPECT().SaveObject(mock.Anything, mock.AnythingOfType("models.Item")).
					Once().
					Return(0, nil)
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, rr.Code)
			},
		},
		{
			name: "successful save",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				body := []byte(`{"some":"body"}`)
				req, _ := http.NewRequest(http.MethodPut, "foo/bar", bytes.NewBuffer(body))
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			prepareStore: func(store *mocks.Storage) {
				store.EXPECT().SaveObject(mock.Anything, mock.AnythingOfType("models.Item")).
					Once().
					Return(1, nil)
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewStorage(t)
			require.NotNil(t, store)
			if tc.prepareStore != nil {
				tc.prepareStore(store)
			}

			h := &Handler{
				log:   log,
				store: store,
			}

			var (
				req = tc.giveRequest()
				rr  = httptest.NewRecorder()
			)

			h.AddObject(rr, req)
			tc.checkResult(t, rr)
		})
	}
}

func TestHandler_Object(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	cases := []struct {
		name         string
		prepareStore func(store *mocks.Storage)
		giveRequest  func() *http.Request
		checkResult  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "invalid object id format",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "invalid")

				req, _ := http.NewRequest(http.MethodGet, "foo/bar", http.NoBody)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Contains(t, rr.Body.String(), "failed get object id")
			},
		},
		{
			name: "store error",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				req, _ := http.NewRequest(http.MethodGet, "foo/bar", http.NoBody)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			prepareStore: func(store *mocks.Storage) {
				store.EXPECT().GetObject(mock.Anything, mock.AnythingOfType("int")).
					Once().
					Return(models.Item{}, errors.New("some error"))
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "failed get object")
			},
		},
		{
			name: "not found",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				req, _ := http.NewRequest(http.MethodGet, "foo/bar", http.NoBody)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			prepareStore: func(store *mocks.Storage) {
				store.EXPECT().GetObject(mock.Anything, mock.AnythingOfType("int")).
					Once().
					Return(models.Item{}, models.ErrNotFound)
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rr.Code)
				assert.Contains(t, rr.Body.String(), "failed get object")
			},
		},
		{
			name: "success",
			giveRequest: func() *http.Request {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("objectID", "1")

				req, _ := http.NewRequest(http.MethodGet, "foo/bar", http.NoBody)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				return req
			},
			prepareStore: func(store *mocks.Storage) {
				store.EXPECT().GetObject(mock.Anything, mock.AnythingOfType("int")).
					Once().
					Return(models.Item{
						Body: []byte(`{"some":"body"}`),
					}, nil)
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Contains(t, rr.Body.String(), "{\"some\":\"body\"}")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewStorage(t)
			require.NotNil(t, store)
			if tc.prepareStore != nil {
				tc.prepareStore(store)
			}

			h := &Handler{
				log:   log,
				store: store,
			}

			var (
				req = tc.giveRequest()
				rr  = httptest.NewRecorder()
			)

			h.Object(rr, req)
			tc.checkResult(t, rr)
		})
	}
}
