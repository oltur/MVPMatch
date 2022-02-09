package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/oltur/mvp-match/httputil"
	"github.com/oltur/mvp-match/model"
	"github.com/oltur/mvp-match/types"
	"net/http"
	"strconv"
	"strings"
)

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}

func (c *Controller) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gwtToken := ctx.GetHeader("Authorization")
		if len(gwtToken) == 0 {
			httputil.NewError(ctx, http.StatusUnauthorized, model.ErrUnauthorized)
			ctx.Abort()
			return
		}

		gwtToken = strings.Replace(gwtToken, "Bearer ", "", -1)

		userId, token, expires, err := c.validateGwt(gwtToken)
		if err != nil {
			err = model.ErrCannotValidateUserToken
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		ok, err := model.VerifyToken(userId, token, expires)
		if err != nil {
			err = model.ErrCannotValidateUserToken
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		if !ok {
			err = model.ErrAccessDenied
			httputil.NewError(ctx, http.StatusForbidden, err)
			ctx.Abort()
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func (c *Controller) getKey() (key string, err error) {
	// TODO: Read from env
	key = "xxxxxxxxxxx"
	return key, nil
}

func (c *Controller) createGwt(userId string, token string, expires int64) (res string, err error) {
	key, err := c.getKey()
	binKey := []byte(key)
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  userId,
		"token":   token,
		"expires": strconv.FormatInt(expires, 10),
	})

	// Sign and get the complete encoded token as a string using the secret
	res, err = jwtToken.SignedString(binKey)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (c *Controller) validateGwt(gwtToken string) (userId string, token string, expires int64, err error) {
	// sample token string taken from the New example
	key, err := c.getKey()
	binKey := []byte(key)

	t, err := jwt.Parse(gwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return binKey, nil
	})

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		i, err := strconv.ParseInt(claims["expires"].(string), 10, 64)
		if err != nil {
			return "", "", 0, nil
		}

		userId = claims["userId"].(string)
		token = claims["token"].(string)
		expires = i

		return userId, token, expires, nil
	} else {
		return "", "", 0, nil
	}
}

// Ping godoc
// @Summary      Ping
// @Description  pings
// @Tags         Tools
// @Accept       json
// @Produce      json
// @Success      200  {string}  string "Pong"
// @Router       /tools/ping [put]
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Pong")
}

func (c *Controller) getUserIdFromContext(ctx *gin.Context) (res types.Id, err error) {
	x, exists := ctx.Get("userId")
	if !exists {
		err = model.ErrUserNotFoundInContext
		return
	}
	s, ok := x.(string)
	if !ok {
		err = model.ErrUserNotFoundInContext
		return
	}
	res = types.Id(s)
	return
}
