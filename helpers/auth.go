package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckRole(c *gin.Context, role string) error {
	userRole := c.GetString("user_type")

	if userRole != role {
		return errors.New("Unauthorized to access this resource")
	}

	return nil
}

func MatchId(c *gin.Context, id string) error {
	userId := c.GetString("id")
	userRole := c.GetString("user_type")
	if userRole != "ADMIN" && userId != id {
		return errors.New("Unauthorized to access this resource")
	}

	return nil
}
