package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// SubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectsV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	subjects, err := h.subjectService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.SubjectsV1List200JSONResponse{
		Subjects: ToAPISubjects(subjects),
	})
}

// SubjectsV1Detail 科目を詳細取得する
func (h *Handler) SubjectsV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	subject, err := h.subjectService.Detail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.SubjectsV1Detail200JSONResponse{
		Subject: ToAPISubject(subject),
	})
}

// SubjectsV1Create 科目を作成する
func (h *Handler) SubjectsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.subjectService.Create(c.Request.Context(), ToDomainSubjectRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.SubjectsV1Create201JSONResponse{
		Subject: ToAPISubject(subject),
	})
}

// SubjectsV1Update 科目を更新する
func (h *Handler) SubjectsV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.subjectService.Update(c.Request.Context(), id, ToDomainSubjectRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.SubjectsV1Update200JSONResponse{
		Subject: ToAPISubject(subject),
	})
}

// SubjectsV1Delete 科目を削除する
func (h *Handler) SubjectsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	if err := h.subjectService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
