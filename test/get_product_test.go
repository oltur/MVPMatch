package test

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/oltur/mvp-match/controller"
	"github.com/oltur/mvp-match/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProduct(t *testing.T) {
	router := controller.SetupRouter()

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
	if data.ID != "1" {
		t.Fatal("not that product")
	}
}
