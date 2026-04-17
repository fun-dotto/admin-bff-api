package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
)

func (h *Handler) FacultyRoomsV1List(c *gin.Context, _ api.FacultyRoomsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *Handler) FacultyRoomsV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *Handler) FacultyRoomsV1Delete(c *gin.Context, _ string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
