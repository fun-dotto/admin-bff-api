package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FacultiesV1List 教員一覧を取得する
func (h *Handler) FacultiesV1List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacultiesV1Detail 教員を詳細取得する
func (h *Handler) FacultiesV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacultiesV1Create 教員を作成する
func (h *Handler) FacultiesV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacultiesV1Update 教員を更新する
func (h *Handler) FacultiesV1Update(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacultiesV1Delete 教員を削除する
func (h *Handler) FacultiesV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
