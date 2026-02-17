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

func TestFacultiesV1List(t *testing.T) {
	tests := []struct {
		name               string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に教員一覧が取得できる",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")

				faculties, ok := response["faculties"].([]interface{})
				assert.True(t, ok, "facultiesフィールドが配列ではありません")
				assert.NotEmpty(t, faculties, "教員が空です")
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
				faculties, ok := response["faculties"].([]interface{})
				assert.True(t, ok, "facultiesフィールドが配列ではありません")
				assert.NotEmpty(t, faculties, "教員が空です")
				assert.Len(t, faculties, 1, "MockRepositoryは1件返すはずです")
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
			name:           "教員のフィールドが正しく返される",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Faculties []api.SubjectServiceFaculty `json:"faculties"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Faculties, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", response.Faculties[0].Id)
				assert.Equal(t, "教員1", response.Faculties[0].Name)
				assert.Equal(t, "faculty1@example.com", response.Faculties[0].Email)
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
			mockRepo := repository.NewMockFacultyRepository()
			facultyService := service.NewFacultyService(mockRepo)
			h := handler.NewHandler().WithFacultyService(facultyService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.FacultiesV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestFacultiesV1Detail(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に教員詳細が取得できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Faculty api.SubjectServiceFaculty `json:"faculty"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Faculty.Id)
				assert.Equal(t, "教員1", response.Faculty.Name)
				assert.Equal(t, "faculty1@example.com", response.Faculty.Email)
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
			mockRepo := repository.NewMockFacultyRepository()
			facultyService := service.NewFacultyService(mockRepo)
			h := handler.NewHandler().WithFacultyService(facultyService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.FacultiesV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestFacultiesV1Create(t *testing.T) {
	tests := []struct {
		name               string
		request            api.SubjectServiceFacultyRequest
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に教員を作成できる",
			request: api.SubjectServiceFacultyRequest{
				Name:  "新しい教員",
				Email: "newfaculty@example.com",
			},
			withAdminClaim: true,
			wantCode:       http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Faculty api.SubjectServiceFaculty `json:"faculty"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Faculty.Id)
				assert.Equal(t, "新しい教員", response.Faculty.Name)
				assert.Equal(t, "newfaculty@example.com", response.Faculty.Email)
			},
		},
		{
			name: "developerクレームのみでも作成できる",
			request: api.SubjectServiceFacultyRequest{
				Name:  "developer経由の教員",
				Email: "devfaculty@example.com",
			},
			withDeveloperClaim: true,
			wantCode:           http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Faculty api.SubjectServiceFaculty `json:"faculty"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Faculty.Id)
				assert.Equal(t, "developer経由の教員", response.Faculty.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			request: api.SubjectServiceFacultyRequest{
				Name:  "新しい教員",
				Email: "newfaculty@example.com",
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
			request: api.SubjectServiceFacultyRequest{
				Name:  "新しい教員",
				Email: "newfaculty@example.com",
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
			mockRepo := repository.NewMockFacultyRepository()
			facultyService := service.NewFacultyService(mockRepo)
			h := handler.NewHandler().WithFacultyService(facultyService)
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
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/faculties", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.FacultiesV1Create(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestFacultiesV1Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		request        api.SubjectServiceFacultyRequest
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に教員を更新できる",
			id:   "1",
			request: api.SubjectServiceFacultyRequest{
				Name:  "更新された教員",
				Email: "updated@example.com",
			},
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Faculty api.SubjectServiceFaculty `json:"faculty"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Faculty.Id)
				assert.Equal(t, "更新された教員", response.Faculty.Name)
				assert.Equal(t, "updated@example.com", response.Faculty.Email)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			request: api.SubjectServiceFacultyRequest{
				Name:  "更新された教員",
				Email: "updated@example.com",
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
			request: api.SubjectServiceFacultyRequest{
				Name:  "更新された教員",
				Email: "updated@example.com",
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
			mockRepo := repository.NewMockFacultyRepository()
			facultyService := service.NewFacultyService(mockRepo)
			h := handler.NewHandler().WithFacultyService(facultyService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/faculties/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.FacultiesV1Update(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestFacultiesV1Delete(t *testing.T) {
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
			name:           "正常に教員を削除できる",
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
			mockRepo := repository.NewMockFacultyRepository()
			facultyService := service.NewFacultyService(mockRepo)
			h := handler.NewHandler().WithFacultyService(facultyService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.FacultiesV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
