package main

import (
	"github.com/oltur/mvp-match/controller"
	_ "github.com/oltur/mvp-match/docs"
)

// @title           MVP Match test task
// @version         0.1
// @description     This is a MVP Match test task, based on celler example.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url NA
// http://www.swagger.io/support
// @contact.email  olturua@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	r, _ := controller.SetupRouter()
	r.Run(":8081")
}
