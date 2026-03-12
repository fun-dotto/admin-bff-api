package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
)

// SubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectsV1List(c *gin.Context, params api.SubjectsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	clientParams := &academic_api.SubjectsV1ListParams{
		Q:                       params.Q,
		Grade:                   convertSlicePtr[api.DottoFoundationV1Grade, academic_api.DottoFoundationV1Grade](params.Grade),
		Courses:                 convertSlicePtr[api.DottoFoundationV1Course, academic_api.DottoFoundationV1Course](params.Courses),
		Class:                   convertSlicePtr[api.DottoFoundationV1Class, academic_api.DottoFoundationV1Class](params.Class),
		Classification:          convertSlicePtr[api.DottoFoundationV1SubjectClassification, academic_api.DottoFoundationV1SubjectClassification](params.Classification),
		Year:                    params.Year,
		Semester:                convertSlicePtr[api.DottoFoundationV1CourseSemester, academic_api.DottoFoundationV1CourseSemester](params.Semester),
		RequirementType:         convertSlicePtr[api.DottoFoundationV1SubjectRequirementType, academic_api.DottoFoundationV1SubjectRequirementType](params.RequirementType),
		CulturalSubjectCategory: convertSlicePtr[api.DottoFoundationV1CulturalSubjectCategory, academic_api.DottoFoundationV1CulturalSubjectCategory](params.CulturalSubjectCategory),
	}

	response, err := h.academicClient.SubjectsV1ListWithResponse(c.Request.Context(), clientParams)
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

	response, err := h.academicClient.SubjectsV1DetailWithResponse(c.Request.Context(), id)
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

// SubjectsV1Upsert 科目を作成または更新する
func (h *Handler) SubjectsV1Upsert(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.SubjectsV1UpsertWithResponse(c.Request.Context(), req)
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

	response, err := h.academicClient.SubjectsV1DeleteWithResponse(c.Request.Context(), id)
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
