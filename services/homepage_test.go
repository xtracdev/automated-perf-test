package services

import (

	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"github.com/go-chi/chi"
)

func TestStartUiMode(t *testing.T) {

	r := chi.NewRouter()

	r.Mount("/", StartUiMode())

	assert.IsType(t, &chi.Mux{}, r)

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t,htmlType, resp.Header().Get(contentTypeHeader))
}
