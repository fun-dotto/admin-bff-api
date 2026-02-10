package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubjectCategoriesV1List 科目群・科目区分一覧を取得する
func (h *Handler) SubjectCategoriesV1List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectCategoriesV1Detail 科目群・科目区分を詳細取得する
func (h *Handler) SubjectCategoriesV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectCategoriesV1Create 科目群・科目区分を作成する
func (h *Handler) SubjectCategoriesV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectCategoriesV1Update 科目群・科目区分を更新する
func (h *Handler) SubjectCategoriesV1Update(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// SubjectCategoriesV1Delete 科目群・科目区分を削除する
func (h *Handler) SubjectCategoriesV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
