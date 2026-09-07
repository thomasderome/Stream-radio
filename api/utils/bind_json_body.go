package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BodyBinder(c *gin.Context, T any) error {
	if err := c.ShouldBindJSON(&T); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}
