package controller

import (
	"github.com/oltur/mvp-match/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
)

// Deposit godoc
// @Summary      Deposit money
// @Description  Deposit a coin of given value for current Buyer user
// @Tags         Vending Machine
// @Accept       json
// @Produce      json
// @Param        coinValue     query     int     false  "Coin value"       Enums(5, 10, 20, 50, 100)
// @Success      200  {object}  model.DepositResponse
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /deposit [post]
func (c *Controller) Deposit(ctx *gin.Context) {
	var s string
	var err error
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	s = ctx.Query("coinValue")
	coin := &model.Coin{}
	coin.Value, err = strconv.Atoi(s)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = coin.Validation()
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	if user.Role != model.UserRoleBuyer {
		err = model.ErrInvalidBuyer
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user.Deposit = user.Deposit + coin.Value

	err = model.UserSave(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := &model.DepositResponse{Deposit: user.Deposit}

	ctx.JSON(http.StatusOK, res)
}

// Buy godoc
// @Summary      Buy product
// @Description  Buy given amount of given product for current Buyer user
// @Tags         Vending Machine
// @Accept       json
// @Produce      json
// @Param        productId   query      string  true  "Product ID"
// @Param        amountOfProducts     query     int     false  "Amount of products"
// @Success      200  {object}  model.BuyResponse
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /buy [post]
func (c *Controller) Buy(ctx *gin.Context) {
	var s string
	var err error
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	s = ctx.Query("productId")
	productId := types.Id(s)
	s = ctx.Query("amountOfProducts")
	amountOfProducts, err := strconv.Atoi(s)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// load and validate data
	user, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	if user.Role != model.UserRoleBuyer {
		err = model.ErrInvalidBuyer
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	product, err := model.ProductOne(productId)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	totalCost := product.Cost * amountOfProducts

	if user.Deposit < totalCost {
		err = model.ErrNotEnoughDeposit
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if product.AmountAvailable < amountOfProducts {
		err = model.ErrNotEnoughAmount
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// check change
	totalChange := user.Deposit - totalCost
	change, err := c.calculateChange(totalChange)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// buy!
	product.AmountAvailable = product.AmountAvailable - amountOfProducts
	user.Deposit = 0

	err = model.ProductSave(product)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = model.UserSave(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := &model.BuyResponse{Total: totalCost, ProductName: product.ProductName, Change: change}

	ctx.JSON(http.StatusOK, res)

}

// Reset godoc
// @Summary      Reset deposit
// @Description  Reset current user deposit
// @Tags         Vending Machine
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.BuyResponse
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /reset [post]
func (c *Controller) Reset(ctx *gin.Context) {
	var s string
	var err error
	x, exists := ctx.Get("userId")
	if !exists {
		err = model.ErrUserNotFoundInContext
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	s, ok := x.(string)
	if !ok {
		err = model.ErrUserNotFoundInContext
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	userId := types.Id(s)

	// load and validate data
	user, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	if user.Role != model.UserRoleBuyer {
		err = model.ErrInvalidBuyer
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	change, err := c.calculateChange(user.Deposit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = model.UserResetDeposit(user.ID)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := &model.ResetResponse{Change: change}

	ctx.JSON(http.StatusOK, res)

}

// --------------- implementation details -------------

func (c Controller) calculateChange(totalChange int) (res []*model.Coin, err error) {
	res = make([]*model.Coin, 0, 1)
	for totalChange > 0 {
		coin := &model.Coin{}
		if totalChange >= 100 {
			coin.Value = 100
		} else if totalChange >= 50 {
			coin.Value = 50
		} else if totalChange >= 20 {
			coin.Value = 20
		} else if totalChange >= 10 {
			coin.Value = 10
		} else if totalChange >= 5 {
			coin.Value = 5
		} else {
			err = model.ErrInvalidCost
			return
		}
		err = coin.Validation()
		if err != nil {
			return
		}
		res = append(res, coin)
		totalChange = totalChange - coin.Value
	}
	return
}
