package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DayOfWeekTimetableSlotsV1List 曜日・時限一覧を取得する
func (h *Handler) DayOfWeekTimetableSlotsV1List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DayOfWeekTimetableSlotsV1Detail 曜日・時限を詳細取得する
func (h *Handler) DayOfWeekTimetableSlotsV1Detail(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DayOfWeekTimetableSlotsV1Create 曜日・時限を作成する
func (h *Handler) DayOfWeekTimetableSlotsV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DayOfWeekTimetableSlotsV1Update 曜日・時限を更新する
func (h *Handler) DayOfWeekTimetableSlotsV1Update(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DayOfWeekTimetableSlotsV1Delete 曜日・時限を削除する
func (h *Handler) DayOfWeekTimetableSlotsV1Delete(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
