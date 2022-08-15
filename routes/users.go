package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/golang-jwtAuth-project/controllers"
	middleware "github.com/golang-jwtAuth-project/middlewares"
)

func UserRouters(incomingRouter *gin.Engine) {
	incomingRouter.Use(middleware.Authentication())
	incomingRouter.GET("/users/", controller.Get())
	incomingRouter.PATCH("/users/:id", controller.Update())
	incomingRouter.DELETE("/users/:id", controller.Delete())
	incomingRouter.GET("users/:id", controller.GetById())
}
