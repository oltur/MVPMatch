package test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/oltur/mvp-match/controller"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
	"github.com/oltur/mvp-match/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuyOk(t *testing.T) {
	userId := types.Id("3")
	productId := types.Id("1")
	amountOfProducts := 2
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	p, err := model.ProductOne(productId)
	if err != nil {
		t.Fatal(err)
	}
	amountBefore := p.AmountAvailable
	model.UserResetDeposit(userId)
	coinValue := 100
	_, err = doTestDeposit(coinValue, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	data, err := doTestBuyOk(productId, amountOfProducts, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	if data.Total != 40 {
		t.Fatal("not the right total")
	}
	if data.ProductName != "Product #1" {
		t.Fatal("not the right product name")
	}
	if len(data.Change) != 2 || data.Change[0].Value != 50 || data.Change[1].Value != 10 {
		t.Fatal("not the right change")
	}
	p, err = model.ProductOne(productId)
	if err != nil {
		t.Fatal(err)
	}
	amountAfter := p.AmountAvailable

	assert.Equal(t, amountBefore-amountAfter, amountOfProducts)
}

func TestBuyFailedTwice(t *testing.T) {
	userId := types.Id("3")
	productId := types.Id("1")
	amountOfProducts := 2
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue := 100
	_, err = doTestDeposit(coinValue, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	// first time
	_, err = doTestBuyOk(productId, amountOfProducts, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	// second time (no more deposit)
	res, err := doTestBuyFail(productId, amountOfProducts, gwtToken, router)
	if res.Message != model.ErrNotEnoughDeposit.Error() {
		t.Fatal(err)
	}
}

func TestBuyFailedNotEnoughDeposit(t *testing.T) {
	userId := types.Id("3")
	productId := types.Id("1")
	amountOfProducts := 2
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue := 5
	_, err = doTestDeposit(coinValue, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	res, err := doTestBuyFail(productId, amountOfProducts, gwtToken, router)
	if res.Message != model.ErrNotEnoughDeposit.Error() {
		t.Fatal(err)
	}
}

func TestBuyFailedWrongRole(t *testing.T) {
	userId := types.Id("2")
	productId := types.Id("1")
	amountOfProducts := 2
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #2, Seller", "2")
	if err != nil {
		t.Fatal(err)
	}
	res, err := doTestBuyFail(productId, amountOfProducts, gwtToken, router)
	if res.Message != model.ErrInvalidBuyer.Error() {
		t.Fatal(err)
	}
}

func TestBuyFailedWrongProduct(t *testing.T) {
	userId := types.Id("3")
	productId := types.Id("999")
	amountOfProducts := 1
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	res, err := doTestBuyFail(productId, amountOfProducts, gwtToken, router)
	if res.Message != model.ErrNotFound.Error() {
		t.Fatal(err)
	}
}

func TestBuyFailedNotEnoughInventory(t *testing.T) {
	userId := types.Id("3")
	productId := types.Id("2")
	amountOfProducts := 2
	router, c := controller.SetupRouter()
	model.UserLogout(userId)
	gwtToken, _, err := c.DoLogin("User #3, Buyer", "3")
	if err != nil {
		t.Fatal(err)
	}
	model.UserResetDeposit(userId)
	coinValue := 100
	_, err = doTestDeposit(coinValue, gwtToken, router)
	if err != nil {
		t.Fatal(err)
	}
	res, err := doTestBuyFail(productId, amountOfProducts, gwtToken, router)
	if res.Message != model.ErrNotEnoughAmount.Error() {
		t.Fatal(err)
	}
}

// ------- implementation details ---------------

func doTestBuyOk(productId types.Id, amountOfProducts int, gwtToken string, router *gin.Engine) (res model.BuyResponse, err error) {
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/buy?productId=%s&amountOfProducts=%d", productId, amountOfProducts)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code < 200 || w.Code >= 300 {
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

func doTestBuyFail(productId types.Id, amountOfProducts int, gwtToken string, router *gin.Engine) (res httputil.HTTPError, err error) {
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/buy?productId=%s&amountOfProducts=%d", productId, amountOfProducts)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code >= 200 && w.Code < 300 {
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
