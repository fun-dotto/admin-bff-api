package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
	"github.com/fun-dotto/api-template/internal/middleware"
	"github.com/fun-dotto/api-template/internal/repository"
	"github.com/fun-dotto/api-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestContext(withAdminClaim bool) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	if withAdminClaim {
		token := &auth.Token{
			Claims: map[string]interface{}{
				"admin": true,
			},
		}
		c.Set(middleware.FirebaseTokenContextKey, token)
	}

	return w, c
}

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
	gin.SetMode(gin.TestMode)
	now := time.Now()
	until := now.Add(24 * time.Hour)

	tests := []struct {
		name         string
		setupMock    func() *repository.MockAnnouncementRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせ一覧が取得できる",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					ListFunc: func(ctx context.Context) ([]domain.Announcement, error) {
						return []domain.Announcement{
							{
								ID:             "1",
								Title:          "お知らせ1",
								URL:            "https://example.com/1",
								AvailableFrom:  now,
								AvailableUntil: &until,
							},
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcements []api.AnnouncementServiceAnnouncement `json:"announcements"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Len(t, response.Announcements, 1)
				assert.Equal(t, "1", response.Announcements[0].Id)
				assert.Equal(t, "お知らせ1", response.Announcements[0].Title)
				assert.Equal(t, "https://example.com/1", response.Announcements[0].Url)
			},
		},
		{
			name: "developerクレームのみでも一覧が取得できる",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					ListFunc: func(ctx context.Context) ([]domain.Announcement, error) {
						return []domain.Announcement{
							{
								ID:             "1",
								Title:          "お知らせ1",
								URL:            "https://example.com/1",
								AvailableFrom:  now,
								AvailableUntil: &until,
							},
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcements []api.AnnouncementServiceAnnouncement `json:"announcements"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Announcements, 1)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
			},
			customClaims: nil,
			wantCode:     http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
		{
			name: "admin/developer以外のクレームのみのトークンでは403エラー",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
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
			mockRepo := tt.setupMock()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler().WithAnnouncementService(announcementService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
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
	gin.SetMode(gin.TestMode)
	now := time.Now()
	until := now.Add(24 * time.Hour)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockAnnouncementRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせ詳細が取得できる",
			id:   "1",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					DetailFunc: func(ctx context.Context, id string) (*domain.Announcement, error) {
						return &domain.Announcement{
							ID:             id,
							Title:          "お知らせ" + id,
							URL:            "https://example.com/" + id,
							AvailableFrom:  now,
							AvailableUntil: &until,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
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
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
			},
			customClaims: nil,
			wantCode:     http.StatusUnauthorized,
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
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
			mockRepo := tt.setupMock()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler().WithAnnouncementService(announcementService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.AnnouncementsV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestAnnouncementsV1Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	until := now.Add(24 * time.Hour)

	tests := []struct {
		name         string
		request      api.AnnouncementServiceAnnouncementRequest
		setupMock    func() *repository.MockAnnouncementRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせを作成できる",
			request: api.AnnouncementServiceAnnouncementRequest{
				Title:          "新しいお知らせ",
				Url:            "https://example.com/new",
				AvailableFrom:  now,
				AvailableUntil: &until,
			},
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					CreateFunc: func(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
						return &domain.Announcement{
							ID:             "created-id",
							Title:          req.Title,
							URL:            req.URL,
							AvailableFrom:  req.AvailableFrom,
							AvailableUntil: req.AvailableUntil,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusCreated,
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					CreateFunc: func(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
						return &domain.Announcement{
							ID:             "created-id",
							Title:          req.Title,
							URL:            req.URL,
							AvailableFrom:  req.AvailableFrom,
							AvailableUntil: req.AvailableUntil,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Announcement api.AnnouncementServiceAnnouncement `json:"announcement"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "created-id", response.Announcement.Id)
				assert.Equal(t, "developer経由のお知らせ", response.Announcement.Title)
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
			},
			customClaims: nil,
			wantCode:     http.StatusUnauthorized,
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
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
			mockRepo := tt.setupMock()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler().WithAnnouncementService(announcementService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

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
	gin.SetMode(gin.TestMode)
	now := time.Now()
	until := now.Add(24 * time.Hour)

	tests := []struct {
		name         string
		id           string
		request      api.AnnouncementServiceAnnouncementRequest
		setupMock    func() *repository.MockAnnouncementRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					UpdateFunc: func(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
						return &domain.Announcement{
							ID:             id,
							Title:          req.Title,
							URL:            req.URL,
							AvailableFrom:  req.AvailableFrom,
							AvailableUntil: req.AvailableUntil,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
			},
			customClaims: nil,
			wantCode:     http.StatusUnauthorized,
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
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
			mockRepo := tt.setupMock()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler().WithAnnouncementService(announcementService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

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
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockAnnouncementRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせを削除できる",
			id:   "1",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					DeleteFunc: func(ctx context.Context, id string) error {
						return nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusNoContent,
			validate:     nil,
		},
		{
			name: "developerクレームのみでも削除できる",
			id:   "1",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					DeleteFunc: func(ctx context.Context, id string) error {
						return nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusNoContent,
			validate:     nil,
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
			},
			customClaims: nil,
			wantCode:     http.StatusUnauthorized,
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
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{}
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
			mockRepo := tt.setupMock()
			announcementService := service.NewAnnouncementService(mockRepo)
			h := handler.NewHandler().WithAnnouncementService(announcementService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.AnnouncementsV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
