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

func TestGetUserOk(t *testing.T) {
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()
	model.UserLogout("1")
	gwtToken, _, err := c.DoLogin("User #1, Seller", "1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/user/1", nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var data model.User
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.ID != "1" || data.UserName != "User #1, Seller" {
		t.Fatal("not the right user")
	}
}

func TestGetUserFailedDoesNotExistNotAdmin(t *testing.T) {
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()
	model.UserLogout("1")
	gwtToken, _, err := c.DoLogin("User #1, Seller", "1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/user/999", nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrAccessDenied.Error() {
		t.Fatal("not the right message")
	}
}

func TestGetUserFailedAccessDeniedNotAdmin(t *testing.T) {
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()
	model.UserLogout("1")
	gwtToken, _, err := c.DoLogin("User #1, Seller", "1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/user/2", nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrAccessDenied.Error() {
		t.Fatal("not the right message")
	}
}

func TestGetUserFailedDoesNotExistAdmin(t *testing.T) {
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()
	model.UserLogout("4")
	gwtToken, _, err := c.DoLogin("User #4, Admin", "4")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/user/999", nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrNotFound.Error() {
		t.Fatal("not the right message")
	}
}
