package utilities

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJsonErrorResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid data",
	})
}
