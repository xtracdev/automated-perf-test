package services

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartUiModeSuccesfully(t *testing.T) {

	r := chi.NewRouter()

	r.Mount("/", getIndexPage())

	assert.IsType(t, &chi.Mux{}, r)

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, htmlType, resp.Header().Get(contentTypeHeader))

}

func TestStartUiModeWithInvalidURL(t *testing.T) {

	r := chi.NewRouter()

	r.Mount("/", getIndexPage())

	assert.IsType(t, &chi.Mux{}, r)

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/xxx", nil)
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
