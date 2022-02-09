package controller

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	c := NewController()

	v1 := r.Group("/api/v1")
	{
		deposit := v1.Group("/deposit")
		{
			deposit.Use(c.Auth())
			deposit.PUT("", c.Deposit)
		}
		buy := v1.Group("/buy")
		{
			buy.Use(c.Auth())
			buy.PUT("", c.Buy)
		}
		reset := v1.Group("/reset")
		{
			reset.Use(c.Auth())
			reset.PUT("", c.Reset)
		}
		user := v1.Group("/user")
		{
			user.POST("", c.AddUser)
			user.PUT("/login", c.Login)
			user.PUT("/logout/all", c.LogoutAll)
			x := user.PUT("/logout", c.Logout)
			x.Use(c.Auth())
			x = user.GET(":id", c.ShowUser)
			x.Use(c.Auth())
			x = user.GET("", c.ListUsers)
			x.Use(c.Auth())
			x = user.DELETE(":id", c.DeleteUser)
			x.Use(c.Auth())
			x = user.PATCH(":id", c.UpdateUser)
			x.Use(c.Auth())
		}
		product := v1.Group("/product")
		{
			product.GET(":id", c.ShowProduct)
			product.GET("", c.ListProducts)
			x := product.POST("", c.AddProduct)
			x.Use(c.Auth())
			x = product.DELETE(":id", c.DeleteProduct)
			x.Use(c.Auth())
			x = product.PATCH(":id", c.UpdateProduct)
			x.Use(c.Auth())
		}
		tools := v1.Group("/tools")
		{
			tools.GET("/ping", c.Ping)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
