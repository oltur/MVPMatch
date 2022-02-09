package test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oltur/mvp-match/controller"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
	"github.com/oltur/mvp-match/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDepositOk(t *testing.T) {
	userId := types.Id("3")
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue := 5
	data, err := doTestDeposit(coinValue, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	if data.Deposit != coinValue {
		t.Fatal("not the right deposit")
	}
}

func TestDepositOkTwice(t *testing.T) {
	userId := types.Id("3")
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue1 := 5
	data, err := doTestDeposit(coinValue1, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	coinValue2 := 20
	data, err = doTestDeposit(coinValue2, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	if data.Deposit != coinValue1+coinValue2 {
		t.Fatal("not the right deposit")
	}
}

func TestDepositFailedWrongCoinValue(t *testing.T) {
	userId := types.Id("3")
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue := 4
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/deposit?coinValue=%d", coinValue)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatal(err)
		return
	}
	body := w.Body.String()
	res := &httputil.HTTPError{}
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		t.Fatal(err)
		return
	}
	if res.Message != model.ErrInvalidCoinValue.Error() {
		t.Fatal(err)
		return
	}
}

func TestDepositFailedWrongRole(t *testing.T) {
	userId := types.Id("2")
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #2, Seller", "2")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue := 5
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/deposit?coinValue=%d", coinValue)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatal(err)
		return
	}
	body := w.Body.String()
	res := &httputil.HTTPError{}
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		t.Fatal(err)
		return
	}
	if res.Message != model.ErrInvalidBuyer.Error() {
		t.Fatal(err)
		return
	}
}

// ------- implementation details ---------------

func doTestDeposit(coinValue int, gwtToken string, router *gin.Engine) (res model.DepositResponse, err error) {
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/deposit?coinValue=%d", coinValue)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		err = fmt.Errorf("unexpected http code")
		return
	}
	body := w.Body.String()
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return
	}
	return
}
