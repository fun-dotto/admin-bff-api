package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
)

// CourseRegistrationsV1List 履修情報を取得する
func (h *Handler) CourseRegistrationsV1List(c *gin.Context, params api.CourseRegistrationsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.CourseRegistrationsV1ListWithResponse(c.Request.Context(), &academic_api.CourseRegistrationsV1ListParams{
		UserId:    params.UserId,
		Year:      params.Year,
		Semesters: convertSlice[api.DottoFoundationV1CourseSemester, academic_api.DottoFoundationV1CourseSemester](params.Semesters),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registrations": response.JSON200.CourseRegistrations,
	})
}

// CourseRegistrationsV1Create 履修情報を作成する
func (h *Handler) CourseRegistrationsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.CourseRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.CourseRegistrationsV1CreateWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON201 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"registration": response.JSON201.CourseRegistration,
	})
}

// CourseRegistrationsV1Delete 履修情報を削除する
func (h *Handler) CourseRegistrationsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.CourseRegistrationsV1DeleteWithResponse(c.Request.Context(), id)
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

// RoomsV1List 教室一覧を取得する
func (h *Handler) RoomsV1List(c *gin.Context, params api.RoomsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.RoomsV1ListWithResponse(c.Request.Context(), &academic_api.RoomsV1ListParams{
		Floor: convertSlicePtr[api.DottoFoundationV1Floor, academic_api.DottoFoundationV1Floor](params.Floor),
	})
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

// RoomsV1Create 教室を作成する
func (h *Handler) RoomsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.RoomsV1CreateWithResponse(c.Request.Context(), req)
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

// RoomsV1Detail 教室を詳細取得する
func (h *Handler) RoomsV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.RoomsV1DetailWithResponse(c.Request.Context(), id)
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

// RoomsV1Update 教室を更新する
func (h *Handler) RoomsV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.RoomsV1UpdateWithResponse(c.Request.Context(), id, req)
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

// RoomsV1Delete 教室を削除する
func (h *Handler) RoomsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.RoomsV1DeleteWithResponse(c.Request.Context(), id)
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

// ReservationsV1List 教室の予約一覧を取得する
func (h *Handler) ReservationsV1List(c *gin.Context, id string, params api.ReservationsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.ReservationsV1ListWithResponse(c.Request.Context(), id, &academic_api.ReservationsV1ListParams{
		From:  params.From,
		Until: params.Until,
	})
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

// TimetableItemsV1List 時間割を取得する
func (h *Handler) TimetableItemsV1List(c *gin.Context, params api.TimetableItemsV1ListParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.TimetableItemsV1ListWithResponse(c.Request.Context(), &academic_api.TimetableItemsV1ListParams{
		Year:      params.Year,
		Semesters: convertSlice[api.DottoFoundationV1CourseSemester, academic_api.DottoFoundationV1CourseSemester](params.Semesters),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON200 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": response.JSON200.TimetableItems,
	})
}

// TimetableItemsV1Create 時間割に追加する
func (h *Handler) TimetableItemsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req academic_api.TimetableItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.academicClient.TimetableItemsV1CreateWithResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.JSON201 == nil {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"item": response.JSON201.TimetableItem,
	})
}

// TimetableItemsV1Delete 時間割を削除する
func (h *Handler) TimetableItemsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.TimetableItemsV1DeleteWithResponse(c.Request.Context(), id)
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

// UsersV1Detail ユーザーを取得する
func (h *Handler) UsersV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.userClient.UsersV1DetailWithResponse(c.Request.Context(), id)
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

func convertSlice[From, To ~string](src []From) []To {
	result := make([]To, len(src))
	for i, v := range src {
		result[i] = To(v)
	}
	return result
}
