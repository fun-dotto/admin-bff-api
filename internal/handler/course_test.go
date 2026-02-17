package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/handler"
	"github.com/fun-dotto/api-template/internal/repository"
	"github.com/fun-dotto/api-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoursesV1List(t *testing.T) {
	tests := []struct {
		name               string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常にコース一覧が取得できる",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")

				courses, ok := response["courses"].([]interface{})
				assert.True(t, ok, "coursesフィールドが配列ではありません")
				assert.NotEmpty(t, courses, "コースが空です")
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
				courses, ok := response["courses"].([]interface{})
				assert.True(t, ok, "coursesフィールドが配列ではありません")
				assert.NotEmpty(t, courses, "コースが空です")
				assert.Len(t, courses, 1, "MockRepositoryは1件返すはずです")
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
			name:           "コースのフィールドが正しく返される",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Courses []api.SubjectServiceCourse `json:"courses"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Courses, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", response.Courses[0].Id)
				assert.Equal(t, "コース1", response.Courses[0].Name)
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
		{
			name:         "admin/developer以外のクレームのみのトークンでは403エラー",
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
			mockRepo := repository.NewMockCourseRepository()
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
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
	tests := []struct {
		name           string
		id             string
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常にコース詳細が取得できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Course.Id)
				assert.Equal(t, "コース1", response.Course.Name)
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
		{
			name:         "admin/developer以外のクレームのみのトークンでは403エラー",
			id:           "1",
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
			mockRepo := repository.NewMockCourseRepository()
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
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
	tests := []struct {
		name               string
		request            api.SubjectServiceCourseRequest
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコースを作成できる",
			request: api.SubjectServiceCourseRequest{
				Name: "新しいコース",
			},
			withAdminClaim: true,
			wantCode:       http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Course.Id)
				assert.Equal(t, "新しいコース", response.Course.Name)
			},
		},
		{
			name: "developerクレームのみでも作成できる",
			request: api.SubjectServiceCourseRequest{
				Name: "developer経由のコース",
			},
			withDeveloperClaim: true,
			wantCode:           http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Course.Id)
				assert.Equal(t, "developer経由のコース", response.Course.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			request: api.SubjectServiceCourseRequest{
				Name: "新しいコース",
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
			request: api.SubjectServiceCourseRequest{
				Name: "新しいコース",
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
			mockRepo := repository.NewMockCourseRepository()
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
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
	tests := []struct {
		name           string
		id             string
		request        api.SubjectServiceCourseRequest
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコースを更新できる",
			id:   "1",
			request: api.SubjectServiceCourseRequest{
				Name: "更新されたコース",
			},
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Course api.SubjectServiceCourse `json:"course"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Course.Id)
				assert.Equal(t, "更新されたコース", response.Course.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			request: api.SubjectServiceCourseRequest{
				Name: "更新されたコース",
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
			request: api.SubjectServiceCourseRequest{
				Name: "更新されたコース",
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
			mockRepo := repository.NewMockCourseRepository()
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
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
	tests := []struct {
		name               string
		id                 string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常にコースを削除できる",
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
		{
			name:         "admin/developer以外のクレームのみのトークンでは403エラー",
			id:           "1",
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
			mockRepo := repository.NewMockCourseRepository()
			courseService := service.NewCourseService(mockRepo)
			h := handler.NewHandler().WithCourseService(courseService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.CoursesV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
