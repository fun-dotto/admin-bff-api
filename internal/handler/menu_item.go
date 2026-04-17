package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
)

func (h *Handler) MenuItemsV1List(c *gin.Context, _ api.MenuItemsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
