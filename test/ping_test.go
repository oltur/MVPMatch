package test

import (
	"github.com/go-playground/assert/v2"
	"github.com/oltur/mvp-match/controller"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	router, _ := controller.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tools/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `"Pong"`, w.Body.String())
}
