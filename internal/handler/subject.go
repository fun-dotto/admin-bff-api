package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// SubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectsV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.SubjectsV1ListWithResponse(c.Request.Context())
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

// SubjectsV1Detail 科目を詳細取得する
func (h *Handler) SubjectsV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.SubjectsV1DetailWithResponse(c.Request.Context(), id)
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

// SubjectsV1Create 科目を作成する
func (h *Handler) SubjectsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req subject_api.SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.subjectClient.SubjectsV1CreateWithResponse(c.Request.Context(), req)
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

// SubjectsV1Update 科目を更新する
func (h *Handler) SubjectsV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req subject_api.SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.subjectClient.SubjectsV1UpdateWithResponse(c.Request.Context(), id, req)
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

// SubjectsV1Delete 科目を削除する
func (h *Handler) SubjectsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.SubjectsV1DeleteWithResponse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.StatusCode() == http.StatusNoContent {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
}
