package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
)

// CourseRegistrationsV1List 履修情報を取得する
func (h *Handler) CourseRegistrationsV1List(c *gin.Context, params api.CourseRegistrationsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// CourseRegistrationsV1Create 履修情報を作成する
func (h *Handler) CourseRegistrationsV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// CourseRegistrationsV1Delete 履修情報を削除する
func (h *Handler) CourseRegistrationsV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacilityRoomsV1List 教室一覧を取得する
func (h *Handler) FacilityRoomsV1List(c *gin.Context, params api.FacilityRoomsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacilityRoomsV1Create 教室を作成する
func (h *Handler) FacilityRoomsV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacilityRoomsV1Detail 教室を詳細取得する
func (h *Handler) FacilityRoomsV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacilityRoomsV1Update 教室を更新する
func (h *Handler) FacilityRoomsV1Update(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacilityRoomsV1Delete 教室を削除する
func (h *Handler) FacilityRoomsV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// FacilityReservationsV1List 教室の予約一覧を取得する
func (h *Handler) FacilityReservationsV1List(c *gin.Context, id string, params api.FacilityReservationsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// TimetableItemsV1List 時間割を取得する
func (h *Handler) TimetableItemsV1List(c *gin.Context, params api.TimetableItemsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// TimetableItemsV1Create 時間割に追加する
func (h *Handler) TimetableItemsV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// TimetableItemsV1Delete 時間割を削除する
func (h *Handler) TimetableItemsV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// UserUsersV1Detail ユーザーを取得する
func (h *Handler) UserUsersV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
