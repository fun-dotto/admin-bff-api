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

	clientParams := &subject_api.SubjectsV1ListParams{
		Q:              params.Q,
		Grade:          convertSlicePtr[api.DottoFoundationV1Grade, subject_api.DottoFoundationV1Grade](params.Grade),
		Courses:        convertSlicePtr[api.DottoFoundationV1Course, subject_api.DottoFoundationV1Course](params.Courses),
		Class:          convertSlicePtr[api.DottoFoundationV1Class, subject_api.DottoFoundationV1Class](params.Class),
		Classification: convertSlicePtr[api.DottoFoundationV1SubjectClassification, subject_api.DottoFoundationV1SubjectClassification](params.Classification),
		Year:           params.Year,
	}

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

func convertSlicePtr[From, To ~string](src *[]From) *[]To {
	if src == nil {
		return nil
	}
	result := make([]To, len(*src))
	for i, v := range *src {
		result[i] = To(v)
	}
	return &result
}
