package test

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/oltur/mvp-match/controller"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductOk(t *testing.T) {
	router, _ := controller.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/product/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var data model.Product
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.ID != "1" || data.ProductName != "Product #1" {
		t.Fatal("not the right product")
	}
}

func TestGetProductFailedDoesNotExist(t *testing.T) {
	router, _ := controller.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/product/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrNotFound.Error() {
		t.Fatal("not the right message")
	}
}
