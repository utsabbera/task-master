package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	t.Run("POST /tasks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handler := NewMockHandler(mockCtrl)
		router := NewRouter(handler)
		rw := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/tasks", nil)
		require.NoError(t, err)

		handler.EXPECT().Create(rw, req)

		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	t.Run("GET /tasks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handler := NewMockHandler(mockCtrl)
		router := NewRouter(handler)
		rw := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
		require.NoError(t, err)

		handler.EXPECT().List(rw, req)

		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	t.Run("GET /tasks/{id}", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handler := NewMockHandler(mockCtrl)
		router := NewRouter(handler)
		rw := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/tasks/123", nil)
		require.NoError(t, err)

		handler.EXPECT().Get(rw, req)

		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	t.Run("PATCH /tasks/{id}", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handler := NewMockHandler(mockCtrl)
		router := NewRouter(handler)
		rw := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPatch, "/tasks/123", nil)
		require.NoError(t, err)

		handler.EXPECT().Update(rw, req)

		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	t.Run("DELETE /tasks/{id}", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handler := NewMockHandler(mockCtrl)
		router := NewRouter(handler)
		rw := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodDelete, "/tasks/123", nil)
		require.NoError(t, err)

		handler.EXPECT().Delete(rw, req)

		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	t.Run("POST /chat", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handler := NewMockHandler(mockCtrl)
		router := NewRouter(handler)
		rw := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/chat", nil)
		require.NoError(t, err)

		handler.EXPECT().Chat(rw, req)

		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}
