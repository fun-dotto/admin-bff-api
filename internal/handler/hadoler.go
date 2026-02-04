package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnnouncementsV1List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
