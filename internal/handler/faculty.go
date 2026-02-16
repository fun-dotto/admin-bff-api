package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// FacultiesV1List 教員一覧を取得する
func (h *Handler) FacultiesV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	faculties, err := h.facultyService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.FacultiesV1List200JSONResponse{
		Faculties: ToAPIFaculties(faculties),
	})
}

// FacultiesV1Detail 教員を詳細取得する
func (h *Handler) FacultiesV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	faculty, err := h.facultyService.Detail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.FacultiesV1Detail200JSONResponse{
		Faculty: ToAPIFaculty(faculty),
	})
}

// FacultiesV1Create 教員を作成する
func (h *Handler) FacultiesV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceFacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	faculty, err := h.facultyService.Create(c.Request.Context(), ToDomainFacultyRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.FacultiesV1Create201JSONResponse{
		Faculty: ToAPIFaculty(faculty),
	})
}

// FacultiesV1Update 教員を更新する
func (h *Handler) FacultiesV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceFacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	faculty, err := h.facultyService.Update(c.Request.Context(), id, ToDomainFacultyRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.FacultiesV1Update200JSONResponse{
		Faculty: ToAPIFaculty(faculty),
	})
}

// FacultiesV1Delete 教員を削除する
func (h *Handler) FacultiesV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	if err := h.facultyService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
