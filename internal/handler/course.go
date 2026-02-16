package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// CoursesV1List コース一覧を取得する
func (h *Handler) CoursesV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	courses, err := h.courseService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.CoursesV1List200JSONResponse{
		Courses: ToAPICourses(courses),
	})
}

// CoursesV1Detail コースを詳細取得する
func (h *Handler) CoursesV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	course, err := h.courseService.Detail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.CoursesV1Detail200JSONResponse{
		Course: ToAPICourse(course),
	})
}

// CoursesV1Create コースを作成する
func (h *Handler) CoursesV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.Create(c.Request.Context(), ToDomainCourseRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.CoursesV1Create201JSONResponse{
		Course: ToAPICourse(course),
	})
}

// CoursesV1Update コースを更新する
func (h *Handler) CoursesV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.Update(c.Request.Context(), id, ToDomainCourseRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.CoursesV1Update200JSONResponse{
		Course: ToAPICourse(course),
	})
}

// CoursesV1Delete コースを削除する
func (h *Handler) CoursesV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	if err := h.courseService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
