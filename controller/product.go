package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
	"github.com/oltur/mvp-match/types"
	"github.com/rs/xid"
	"net/http"
)

// ShowProduct godoc
// @Summary      Show a product
// @Description  get string by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  model.Product
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /product/{id} [get]
func (c *Controller) ShowProduct(ctx *gin.Context) {
	s := ctx.Param("id")
	id := types.Id(s)
	product, err := model.ProductOne(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// ListProducts godoc
// @Summary      List products
// @Description  get products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        q    query     string  false  "name search by q"  Format(email)
// @Success      200  {array}   model.Product
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /product [get]
func (c *Controller) ListProducts(ctx *gin.Context) {
	q := ctx.Query("q")
	products, err := model.ProductsAll(q)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// AddProduct godoc
// @Summary      Add product
// @Description  Add new product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        product  body      model.AddProductReq  true  "Add product request"
// @Success      200      {object}  model.Product
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /product [post]
func (c *Controller) AddProduct(ctx *gin.Context) {
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// check Seller role
	user, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if user.Role != model.UserRoleSeller {
		err = model.ErrInvalidSeller
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	var req model.AddProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := req.Validation(); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	product := &model.Product{
		ID:              types.Id(xid.New().String()),
		ProductName:     req.ProductName,
		SellerId:        userId,
		AmountAvailable: req.AmountAvailable,
		Cost:            req.Cost,
	}
	res, err := model.ProductInsert(product)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update by json product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        product  body      model.Product  true  "Update product info"
// @Success      200      {object}  model.UpdateProductRequest
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/{id} [patch]
func (c *Controller) UpdateProduct(ctx *gin.Context) {
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// check Seller role
	user, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if user.Role != model.UserRoleSeller {
		err = model.ErrInvalidSeller
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	var updateProductReq model.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&updateProductReq); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// update
	err = model.ProductUpdate(&updateProductReq)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, updateProductReq)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete by product ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success 	 204  {string} string "Ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/{id} [delete]
func (c *Controller) DeleteProduct(ctx *gin.Context) {
	s := ctx.Param("id")
	id := types.Id(s)

	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// check Seller role
	user, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if user.Role != model.UserRoleSeller {
		err = model.ErrInvalidSeller
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	// check product ownership
	product, err := model.ProductOne(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	if product.SellerId != userId {
		err = model.ErrWrongSeller
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	// delete
	err = model.ProductDelete(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusNoContent, "Ok")
}
