package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// SubjectCategoriesV1List 科目群・科目区分一覧を取得する
func (h *Handler) SubjectCategoriesV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	categories, err := h.subjectCategoryService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.SubjectCategoriesV1List200JSONResponse{
		SubjectCategories: ToAPISubjectCategories(categories),
	})
}

// SubjectCategoriesV1Detail 科目群・科目区分を詳細取得する
func (h *Handler) SubjectCategoriesV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	category, err := h.subjectCategoryService.Detail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.SubjectCategoriesV1Detail200JSONResponse{
		SubjectCategory: ToAPISubjectCategory(category),
	})
}

// SubjectCategoriesV1Create 科目群・科目区分を作成する
func (h *Handler) SubjectCategoriesV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceSubjectCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.subjectCategoryService.Create(c.Request.Context(), ToDomainSubjectCategoryRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.SubjectCategoriesV1Create201JSONResponse{
		SubjectCategory: ToAPISubjectCategory(category),
	})
}

// SubjectCategoriesV1Update 科目群・科目区分を更新する
func (h *Handler) SubjectCategoriesV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceSubjectCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.subjectCategoryService.Update(c.Request.Context(), id, ToDomainSubjectCategoryRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.SubjectCategoriesV1Update200JSONResponse{
		SubjectCategory: ToAPISubjectCategory(category),
	})
}

// SubjectCategoriesV1Delete 科目群・科目区分を削除する
func (h *Handler) SubjectCategoriesV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	if err := h.subjectCategoryService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
