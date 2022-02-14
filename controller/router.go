package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter() (*gin.Engine, *Controller) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 20 << 20 // 20 MiB

	// TODO: Change to specific CORS rules?
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Total-Count", "Authorization"}
	config.ExposeHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Total-Count", "Authorization"}
	r.Use(cors.New(config))

	c := NewController()

	v1 := r.Group("/api/v1")
	{
		deposit := v1.Group("/deposit")
		{
			deposit.Use(c.Auth())
			deposit.POST("", c.Deposit)
		}
		buy := v1.Group("/buy")
		{
			buy.Use(c.Auth())
			buy.POST("", c.Buy)
		}
		reset := v1.Group("/reset")
		{
			reset.Use(c.Auth())
			reset.POST("", c.Reset)
		}
		user := v1.Group("/user")
		{
			user.POST("", c.AddUser)
			user.POST("/login", c.Login)
			user.POST("/logout/all", c.LogoutAll)
			user.POST("/logout", c.Auth(), c.Logout)
			user.GET(":id", c.Auth(), c.ShowUser)
			user.GET("", c.Auth(), c.ListUsers)
			user.DELETE(":id", c.Auth(), c.DeleteUser)
			user.PATCH(":id", c.Auth(), c.UpdateUser)
		}
		product := v1.Group("/product")
		{
			product.GET(":id", c.ShowProduct)
			product.GET("", c.ListProducts)
			product.POST("", c.Auth(), c.AddProduct)
			product.DELETE(":id", c.Auth(), c.DeleteProduct)
			product.PATCH(":id", c.Auth(), c.UpdateProduct)
		}
		tools := v1.Group("/tools")
		{
			tools.GET("/ping", c.Ping)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r, c
}
