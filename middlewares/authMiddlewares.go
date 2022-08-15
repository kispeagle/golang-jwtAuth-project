package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	helper "github.com/golang-jwtAuth-project/helpers"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Token is invalid"})
			return
		}

		claims, err := helper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Token is invalid"})
			c.Abort()
			return
		}

		c.Set("id", claims.Id)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("phone", claims.Phone)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}
