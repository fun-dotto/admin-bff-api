package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectsV1List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectsV1Detail 科目を詳細取得する
func (h *Handler) SubjectsV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectsV1Create 科目を作成する
func (h *Handler) SubjectsV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectsV1Update 科目を更新する
func (h *Handler) SubjectsV1Update(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectsV1Delete 科目を削除する
func (h *Handler) SubjectsV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
