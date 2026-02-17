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

func TestSubjectCategoriesV1List(t *testing.T) {
	tests := []struct {
		name               string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に科目群・科目区分一覧が取得できる",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")

				categories, ok := response["subjectCategories"].([]interface{})
				assert.True(t, ok, "subjectCategoriesフィールドが配列ではありません")
				assert.NotEmpty(t, categories, "科目群・科目区分が空です")
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
				categories, ok := response["subjectCategories"].([]interface{})
				assert.True(t, ok, "subjectCategoriesフィールドが配列ではありません")
				assert.NotEmpty(t, categories, "科目群・科目区分が空です")
				assert.Len(t, categories, 1, "MockRepositoryは1件返すはずです")
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
			name:           "科目群・科目区分のフィールドが正しく返される",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategories []api.SubjectServiceSubjectCategory `json:"subjectCategories"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.SubjectCategories, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", response.SubjectCategories[0].Id)
				assert.Equal(t, "カテゴリ1", response.SubjectCategories[0].Name)
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
			mockRepo := repository.NewMockSubjectCategoryRepository()
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
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
	tests := []struct {
		name           string
		id             string
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に科目群・科目区分詳細が取得できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.SubjectCategory.Id)
				assert.Equal(t, "カテゴリ1", response.SubjectCategory.Name)
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
			mockRepo := repository.NewMockSubjectCategoryRepository()
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
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
	tests := []struct {
		name               string
		request            api.SubjectServiceSubjectCategoryRequest
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目群・科目区分を作成できる",
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "新しいカテゴリ",
			},
			withAdminClaim: true,
			wantCode:       http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.SubjectCategory.Id)
				assert.Equal(t, "新しいカテゴリ", response.SubjectCategory.Name)
			},
		},
		{
			name: "developerクレームのみでも作成できる",
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "developer経由のカテゴリ",
			},
			withDeveloperClaim: true,
			wantCode:           http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.SubjectCategory.Id)
				assert.Equal(t, "developer経由のカテゴリ", response.SubjectCategory.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "新しいカテゴリ",
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
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "新しいカテゴリ",
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
			mockRepo := repository.NewMockSubjectCategoryRepository()
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)
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
	tests := []struct {
		name           string
		id             string
		request        api.SubjectServiceSubjectCategoryRequest
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目群・科目区分を更新できる",
			id:   "1",
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "更新されたカテゴリ",
			},
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					SubjectCategory api.SubjectServiceSubjectCategory `json:"subjectCategory"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.SubjectCategory.Id)
				assert.Equal(t, "更新されたカテゴリ", response.SubjectCategory.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "更新されたカテゴリ",
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
			request: api.SubjectServiceSubjectCategoryRequest{
				Name: "更新されたカテゴリ",
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
			mockRepo := repository.NewMockSubjectCategoryRepository()
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
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
			name:           "正常に科目群・科目区分を削除できる",
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
			mockRepo := repository.NewMockSubjectCategoryRepository()
			categoryService := service.NewSubjectCategoryService(mockRepo)
			h := handler.NewHandler().WithSubjectCategoryService(categoryService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.SubjectCategoriesV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
