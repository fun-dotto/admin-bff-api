package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fun-dotto/admin-bff-api/generated/external/subject_api"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
)

// SubjectsV1List 科目一覧を取得する
func (h *Handler) SubjectsV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	params := &subject_api.SubjectsV1ListParams{
		Q:                       "",
		Grade:                   []subject_api.DottoFoundationV1Grade{},
		Courses:                 []subject_api.DottoFoundationV1Course{},
		Class:                   []subject_api.DottoFoundationV1Class{},
		Classification:          []subject_api.DottoFoundationV1SubjectClassification{},
		Semester:                []subject_api.DottoFoundationV1CourseSemester{},
		RequirementType:         []subject_api.DottoFoundationV1SubjectRequirementType{},
		CalturalSubjectCategory: []subject_api.DottoFoundationV1CulturalSubjectCategory{},
	}

	response, err := h.subjectClient.SubjectsV1ListWithResponse(c.Request.Context(), params)
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

	response, err := h.subjectClient.SubjectsV1UpsertWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, response.JSON200)
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

	if response.StatusCode() != http.StatusNoContent {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.Status(http.StatusNoContent)
}
