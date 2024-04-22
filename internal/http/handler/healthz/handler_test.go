package healthz

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"st-test/internal/http/handler/healthz/mocks"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	store := mocks.NewStore(t)
	require.NotNil(t, store)

	h := NewHandler(log, store)
	require.NotNil(t, h)
}

func TestHandler_Liveness(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	cases := []struct {
		name         string
		prepareStore func(store *mocks.Store)
		checkResult  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			prepareStore: func(store *mocks.Store) {
				store.EXPECT().Check().Once().Return(nil)
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "error",
			prepareStore: func(store *mocks.Store) {
				store.EXPECT().Check().Once().Return(errors.New("some error"))
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewStore(t)
			require.NotNil(t, store)
			if tc.prepareStore != nil {
				tc.prepareStore(store)
			}

			h := NewHandler(log, store)
			require.NotNil(t, h)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "foo/bar", http.NoBody)

			h.Liveness(rr, req)
			tc.checkResult(t, rr)
		})
	}
}

func TestHandler_Readiness(t *testing.T) {
	t.Parallel()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	require.NotNil(t, log)

	cases := []struct {
		name         string
		prepareStore func(store *mocks.Store)
		checkResult  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			prepareStore: func(store *mocks.Store) {
				store.EXPECT().Check().Once().Return(nil)
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "error",
			prepareStore: func(store *mocks.Store) {
				store.EXPECT().Check().Once().Return(errors.New("some error"))
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewStore(t)
			require.NotNil(t, store)
			if tc.prepareStore != nil {
				tc.prepareStore(store)
			}

			h := NewHandler(log, store)
			require.NotNil(t, h)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "foo/bar", http.NoBody)

			h.Readiness(rr, req)
			tc.checkResult(t, rr)
		})
	}
}
