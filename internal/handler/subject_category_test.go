package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
	"github.com/fun-dotto/api-template/internal/repository"
	"github.com/fun-dotto/api-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubjectCategoriesV1List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupMock    func() *repository.MockSubjectCategoryRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目群・科目区分一覧が取得できる",
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
					ListFunc: func(ctx context.Context) ([]domain.SubjectCategory, error) {
						return []domain.SubjectCategory{
							{ID: domain.SubjectCategoryID("1"), Name: "カテゴリ1"},
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategories []api.SubjectServiceSubjectCategory `json:"subjectCategories"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.SubjectCategories, 1)
				assert.Equal(t, "1", response.SubjectCategories[0].Id)
				assert.Equal(t, "カテゴリ1", response.SubjectCategories[0].Name)
			},
		},
		{
			name: "developerクレームのみでも一覧が取得できる",
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
					ListFunc: func(ctx context.Context) ([]domain.SubjectCategory, error) {
						return []domain.SubjectCategory{
							{ID: domain.SubjectCategoryID("1"), Name: "カテゴリ1"},
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategories []api.SubjectServiceSubjectCategory `json:"subjectCategories"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.SubjectCategories, 1)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.SubjectCategoriesV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectCategoriesV1Detail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockSubjectCategoryRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目群・科目区分詳細が取得できる",
			id:   "1",
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
					DetailFunc: func(ctx context.Context, id string) (*domain.SubjectCategory, error) {
						return &domain.SubjectCategory{ID: domain.SubjectCategoryID(id), Name: "カテゴリ" + id}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "1", response.SubjectCategory.Id)
				assert.Equal(t, "カテゴリ1", response.SubjectCategory.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.SubjectCategoriesV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectCategoriesV1Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		request      api.SubjectServiceSubjectCategoryRequest
		setupMock    func() *repository.MockSubjectCategoryRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:    "正常に科目群・科目区分を作成できる",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "新しいカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
					CreateFunc: func(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
						return &domain.SubjectCategory{ID: domain.SubjectCategoryID("created-id"), Name: req.Name}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "created-id", response.SubjectCategory.Id)
				assert.Equal(t, "新しいカテゴリ", response.SubjectCategory.Name)
			},
		},
		{
			name:    "developerクレームのみでも作成できる",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "developer経由のカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
					CreateFunc: func(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
						return &domain.SubjectCategory{ID: domain.SubjectCategoryID("created-id"), Name: req.Name}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "created-id", response.SubjectCategory.Id)
				assert.Equal(t, "developer経由のカテゴリ", response.SubjectCategory.Name)
			},
		},
		{
			name:    "認証トークンがない場合は401エラー",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "新しいカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			name:    "admin/developer以外のクレームのみのトークンでは403エラー",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "新しいカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/subject-categories", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.SubjectCategoriesV1Create(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectCategoriesV1Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		request      api.SubjectServiceSubjectCategoryRequest
		setupMock    func() *repository.MockSubjectCategoryRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:    "正常に科目群・科目区分を更新できる",
			id:      "1",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "更新されたカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
					UpdateFunc: func(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
						return &domain.SubjectCategory{ID: domain.SubjectCategoryID(id), Name: req.Name}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "1", response.SubjectCategory.Id)
				assert.Equal(t, "更新されたカテゴリ", response.SubjectCategory.Name)
			},
		},
		{
			name:    "認証トークンがない場合は401エラー",
			id:      "1",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "更新されたカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			name:    "admin/developer以外のクレームのみのトークンでは403エラー",
			id:      "1",
			request: api.SubjectServiceSubjectCategoryRequest{Name: "更新されたカテゴリ"},
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/subject-categories/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.SubjectCategoriesV1Update(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectCategoriesV1Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockSubjectCategoryRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目群・科目区分を削除できる",
			id:   "1",
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
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
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{
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
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			setupMock: func() *repository.MockSubjectCategoryRepository {
				return &repository.MockSubjectCategoryRepository{}
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
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.SubjectCategoriesV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
