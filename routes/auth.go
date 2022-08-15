package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/golang-jwtAuth-project/controllers"
)

func AuthRouter(incomingRouter *gin.Engine) {
	incomingRouter.POST("/signup", controller.Signup())
	incomingRouter.POST("/login", controller.Login())

}
