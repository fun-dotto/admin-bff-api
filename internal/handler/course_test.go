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

func TestCoursesV1List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupMock    func() *repository.MockCourseRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコース一覧が取得できる",
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
					ListFunc: func(ctx context.Context) ([]domain.Course, error) {
						return []domain.Course{
							{ID: domain.CourseID("1"), Name: "コース1"},
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Courses []api.SubjectServiceCourse `json:"courses"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Courses, 1)
				assert.Equal(t, "1", response.Courses[0].Id)
				assert.Equal(t, "コース1", response.Courses[0].Name)
			},
		},
		{
			name: "developerクレームのみでも一覧が取得できる",
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
					ListFunc: func(ctx context.Context) ([]domain.Course, error) {
						return []domain.Course{
							{ID: domain.CourseID("1"), Name: "コース1"},
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Courses []api.SubjectServiceCourse `json:"courses"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Courses, 1)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.CoursesV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestCoursesV1Detail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockCourseRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコース詳細が取得できる",
			id:   "1",
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
					DetailFunc: func(ctx context.Context, id string) (*domain.Course, error) {
						return &domain.Course{ID: domain.CourseID(id), Name: "コース" + id}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "1", response.Course.Id)
				assert.Equal(t, "コース1", response.Course.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.CoursesV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestCoursesV1Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		request      api.SubjectServiceCourseRequest
		setupMock    func() *repository.MockCourseRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:    "正常にコースを作成できる",
			request: api.SubjectServiceCourseRequest{Name: "新しいコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
					CreateFunc: func(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error) {
						return &domain.Course{ID: domain.CourseID("created-id"), Name: req.Name}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "created-id", response.Course.Id)
				assert.Equal(t, "新しいコース", response.Course.Name)
			},
		},
		{
			name:    "developerクレームのみでも作成できる",
			request: api.SubjectServiceCourseRequest{Name: "developer経由のコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
					CreateFunc: func(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error) {
						return &domain.Course{ID: domain.CourseID("created-id"), Name: req.Name}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "created-id", response.Course.Id)
				assert.Equal(t, "developer経由のコース", response.Course.Name)
			},
		},
		{
			name:    "認証トークンがない場合は401エラー",
			request: api.SubjectServiceCourseRequest{Name: "新しいコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			request: api.SubjectServiceCourseRequest{Name: "新しいコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/courses", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.CoursesV1Create(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestCoursesV1Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		request      api.SubjectServiceCourseRequest
		setupMock    func() *repository.MockCourseRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:    "正常にコースを更新できる",
			id:      "1",
			request: api.SubjectServiceCourseRequest{Name: "更新されたコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
					UpdateFunc: func(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error) {
						return &domain.Course{ID: domain.CourseID(id), Name: req.Name}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "1", response.Course.Id)
				assert.Equal(t, "更新されたコース", response.Course.Name)
			},
		},
		{
			name:    "認証トークンがない場合は401エラー",
			id:      "1",
			request: api.SubjectServiceCourseRequest{Name: "更新されたコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			request: api.SubjectServiceCourseRequest{Name: "更新されたコース"},
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/courses/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.CoursesV1Update(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestCoursesV1Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockCourseRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコースを削除できる",
			id:   "1",
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
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
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{
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
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			setupMock: func() *repository.MockCourseRepository {
				return &repository.MockCourseRepository{}
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
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.CoursesV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
