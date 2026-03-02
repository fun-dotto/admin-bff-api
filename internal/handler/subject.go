package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/subject_api"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
)

// SubjectSubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectSubjectsV1List(c *gin.Context, params api.SubjectSubjectsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	clientParams := &subject_api.SubjectsV1ListParams{}

	response, err := h.subjectClient.SubjectsV1ListWithResponse(c.Request.Context(), clientParams)
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

// SubjectSubjectsV1Detail 科目を詳細取得する
func (h *Handler) SubjectSubjectsV1Detail(c *gin.Context, id string) {
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

// SubjectSubjectsV1Upsert 科目を作成または更新する
func (h *Handler) SubjectSubjectsV1Upsert(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req subject_api.SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.subjectClient.SubjectsV1UpsertWithResponse(c.Request.Context(), req)
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

// SubjectSubjectsV1Delete 科目を削除する
func (h *Handler) SubjectSubjectsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.subjectClient.SubjectsV1DeleteWithResponse(c.Request.Context(), id)
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
