package test

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/oltur/mvp-match/controller"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoginOk(t *testing.T) {
	router, _ := controller.SetupRouter()
	w := httptest.NewRecorder()
	model.UserLogout("1")
	reqBody := `{"userName":"User #1, Seller", "password": "1"}`
	req, _ := http.NewRequest("POST", "/api/v1/user/login", strings.NewReader(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var data model.LoginResponse
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	in29Days := time.Now().Add(29 * 24 * time.Hour).UnixMilli()
	if data.TokenExpires < in29Days {
		t.Fatal("wrong tokenExpires")
	}
}

func TestLoginFailed(t *testing.T) {
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
