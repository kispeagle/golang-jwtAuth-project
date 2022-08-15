package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwtAuth-project/routes"
)

func main() {

	router := gin.New()

	router.Use(gin.Logger())

	routes.AuthRouter(router)
	routes.UserRouters(router)

	router.Run("localhost:8001")
}
