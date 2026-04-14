package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
)

// NotificationV1List 通知一覧を取得する
func (h *Handler) NotificationV1List(c *gin.Context, _ api.NotificationV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// NotificationV1Create 通知を作成する
func (h *Handler) NotificationV1Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// NotificationV1Update 通知を更新する
func (h *Handler) NotificationV1Update(c *gin.Context, _ string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// NotificationV1Delete 通知を削除する
func (h *Handler) NotificationV1Delete(c *gin.Context, _ string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
