package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// CoursesV1List コース一覧を取得する
func (h *Handler) CoursesV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.CoursesV1ListWithResponse(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, response.JSON200)
}

// CoursesV1Detail コースを詳細取得する
func (h *Handler) CoursesV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.CoursesV1DetailWithResponse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, response.JSON200)
}

// CoursesV1Create コースを作成する
func (h *Handler) CoursesV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req subject_api.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.subjectClient.CoursesV1CreateWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON201 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, response.JSON201)
}

// CoursesV1Update コースを更新する
func (h *Handler) CoursesV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req subject_api.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.subjectClient.CoursesV1UpdateWithResponse(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, response.JSON200)
}

// CoursesV1Delete コースを削除する
func (h *Handler) CoursesV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.CoursesV1DeleteWithResponse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.StatusCode() != http.StatusNoContent {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.Status(http.StatusNoContent)
}
