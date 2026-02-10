package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/handler"
	"github.com/fun-dotto/api-template/internal/middleware"
	"github.com/fun-dotto/api-template/internal/repository"
	"github.com/fun-dotto/api-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestContext Firebase認証をモックしたテストコンテキストを作成する
func setupTestContext(withAdminClaim bool) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Requestを初期化
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	if withAdminClaim {
		// Firebaseトークンをモック
		token := &auth.Token{
			Claims: map[string]interface{}{
				"admin": true,
			},
		}
		c.Set(middleware.FirebaseTokenContextKey, token)
	}

	return w, c
}

// setupTestContextWithClaims 指定したクレームでFirebaseトークンをモックしたテストコンテキストを作成する（403など権限不足の検証用）
func setupTestContextWithClaims(claims map[string]interface{}) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	token := &auth.Token{Claims: claims}
	c.Set(middleware.FirebaseTokenContextKey, token)
	return w, c
}

func TestAnnouncementsV1List(t *testing.T) {
	tests := []struct {
		name               string
		withAdminClaim     bool
		withDeveloperClaim bool
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常にお知らせ一覧が取得できる",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")

				announcements, ok := response["announcements"].([]interface{})
				assert.True(t, ok, "announcementsフィールドが配列ではありません")
				assert.NotEmpty(t, announcements, "アナウンスメントが空です")
			},
		},
		{
			name:               "developerクレームのみでも一覧が取得できる",
			withDeveloperClaim: true,
			wantCode:           http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				announcements, ok := response["announcements"].([]interface{})
				assert.True(t, ok, "announcementsフィールドが配列ではありません")
				assert.NotEmpty(t, announcements, "アナウンスメントが空です")
				assert.Len(t, announcements, 1, "MockRepositoryは1件返すはずです")
			},
		},
		{
			name:           "Content-Typeがapplication/jsonである",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			},
		},
		{
			name:           "レスポンスが正しい構造である",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				announcements, ok := response["announcements"].([]interface{})
				assert.True(t, ok, "announcementsフィールドが配列ではありません")
				assert.Len(t, announcements, 1, "MockRepositoryは1件返すはずです")
			},
		},
		{
			name:           "お知らせのフィールドが正しく返される",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcements []api.AnnouncementServiceAnnouncement `json:"announcements"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Announcements, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", response.Announcements[0].Id)
				assert.Equal(t, "お知らせ1", response.Announcements[0].Title)
				assert.Equal(t, "https://example.com/1", response.Announcements[0].Url)
				assert.NotNil(t, response.Announcements[0].AvailableFrom)
				assert.NotNil(t, response.Announcements[0].AvailableUntil)
			},
		},
		{
			name:           "認証トークンがない場合は401エラー",
			withAdminClaim: false,
			wantCode:       http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler(announcementService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.AnnouncementsV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestAnnouncementsV1Detail(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		withAdminClaim bool
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常にお知らせ詳細が取得できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcement api.AnnouncementServiceAnnouncement `json:"announcement"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Announcement.Id)
				assert.Equal(t, "お知らせ1", response.Announcement.Title)
				assert.Equal(t, "https://example.com/1", response.Announcement.Url)
			},
		},
		{
			name:           "認証トークンがない場合は401エラー",
			id:             "1",
			withAdminClaim: false,
			wantCode:       http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler(announcementService)
			w, c := setupTestContext(tt.withAdminClaim)

			h.AnnouncementsV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestAnnouncementsV1Create(t *testing.T) {
	now := time.Now()
	until := now.Add(24 * time.Hour)

	tests := []struct {
		name               string
		request            api.AnnouncementServiceAnnouncementRequest
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせを作成できる",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "新しいお知らせ",
				Url:            "https://example.com/new",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			withAdminClaim: true,
			wantCode:       http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcement api.AnnouncementServiceAnnouncement `json:"announcement"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Announcement.Id)
				assert.Equal(t, "新しいお知らせ", response.Announcement.Title)
				assert.Equal(t, "https://example.com/new", response.Announcement.Url)
			},
		},
		{
			name: "developerクレームのみでも作成できる",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "developer経由のお知らせ",
				Url:            "https://example.com/developer",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			withDeveloperClaim: true,
			wantCode:           http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcement api.AnnouncementServiceAnnouncement `json:"announcement"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Announcement.Id)
				assert.Equal(t, "developer経由のお知らせ", response.Announcement.Title)
				assert.Equal(t, "https://example.com/developer", response.Announcement.Url)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "新しいお知らせ",
				Url:            "https://example.com/new",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			withAdminClaim: false,
			wantCode:       http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
		{
			name: "admin/developer以外のクレームのみのトークンでは403エラー",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "新しいお知らせ",
				Url:            "https://example.com/new",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			customClaims: map[string]interface{}{"user": true},
			wantCode:     http.StatusForbidden,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Insufficient permissions", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler(announcementService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			// リクエストボディを設定
			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/announcements", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.AnnouncementsV1Create(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestAnnouncementsV1Update(t *testing.T) {
	now := time.Now()
	until := now.Add(24 * time.Hour)

	tests := []struct {
		name           string
		id             string
		request        api.AnnouncementServiceAnnouncementRequest
		withAdminClaim bool
		customClaims   map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせを更新できる",
			id:   "1",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "更新されたお知らせ",
				Url:            "https://example.com/updated",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcement api.AnnouncementServiceAnnouncement `json:"announcement"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Announcement.Id)
				assert.Equal(t, "更新されたお知らせ", response.Announcement.Title)
				assert.Equal(t, "https://example.com/updated", response.Announcement.Url)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "更新されたお知らせ",
				Url:            "https://example.com/updated",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			withAdminClaim: false,
			wantCode:       http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
		{
			name: "admin/developer以外のクレームのみのトークンでは403エラー",
			id:   "1",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "更新されたお知らせ",
				Url:            "https://example.com/updated",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			customClaims: map[string]interface{}{"user": true},
			wantCode:     http.StatusForbidden,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Insufficient permissions", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler(announcementService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			// リクエストボディを設定
			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/announcements/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.AnnouncementsV1Update(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestAnnouncementsV1Delete(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		withAdminClaim     bool
		withDeveloperClaim bool
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常にお知らせを削除できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusNoContent,
			validate:       nil,
		},
		{
			name:               "developerクレームのみでも削除できる",
			id:                 "1",
			withDeveloperClaim: true,
			wantCode:           http.StatusNoContent,
			validate:           nil,
		},
		{
			name:           "認証トークンがない場合は401エラー",
			id:             "1",
			withAdminClaim: false,
			wantCode:       http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler(announcementService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.AnnouncementsV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
