package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CoursesV1List コース一覧を取得する
func (h *Handler) CoursesV1List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// CoursesV1Detail コースを詳細取得する
func (h *Handler) CoursesV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// CoursesV1Create コースを作成する
func (h *Handler) CoursesV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// CoursesV1Update コースを更新する
func (h *Handler) CoursesV1Update(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// CoursesV1Delete コースを削除する
func (h *Handler) CoursesV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
