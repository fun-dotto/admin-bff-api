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

func TestSubjectsV1List(t *testing.T) {
	tests := []struct {
		name               string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に科目一覧が取得できる",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")

				subjects, ok := response["subjects"].([]interface{})
				assert.True(t, ok, "subjectsフィールドが配列ではありません")
				assert.NotEmpty(t, subjects, "科目が空です")
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
				subjects, ok := response["subjects"].([]interface{})
				assert.True(t, ok, "subjectsフィールドが配列ではありません")
				assert.NotEmpty(t, subjects, "科目が空です")
				assert.Len(t, subjects, 1, "MockRepositoryは1件返すはずです")
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

				subjects, ok := response["subjects"].([]interface{})
				assert.True(t, ok, "subjectsフィールドが配列ではありません")
				assert.Len(t, subjects, 1, "MockRepositoryは1件返すはずです")
			},
		},
		{
			name:           "科目のフィールドが正しく返される",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subjects []api.SubjectServiceSubject `json:"subjects"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Subjects, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", response.Subjects[0].Id)
				assert.Equal(t, "科目1", response.Subjects[0].Name)
				assert.Equal(t, api.H1, response.Subjects[0].Semester)
				assert.Equal(t, "syllabus-1", response.Subjects[0].SyllabusId)
				assert.NotEmpty(t, response.Subjects[0].Faculty.Id, "Facultyが設定されていること")
				assert.NotEmpty(t, response.Subjects[0].DayOfWeekTimetableSlots, "DayOfWeekTimetableSlotsが設定されていること")
				assert.NotEmpty(t, response.Subjects[0].EligibleAttributes, "EligibleAttributesが設定されていること")
				assert.NotEmpty(t, response.Subjects[0].Requirements, "Requirementsが設定されていること")
				assert.NotEmpty(t, response.Subjects[0].Categories, "Categoriesが設定されていること")
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
			mockRepo := repository.NewMockSubjectRepository()
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.SubjectsV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectsV1Detail(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		withAdminClaim bool
		customClaims   map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に科目詳細が取得できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Subject.Id)
				assert.Equal(t, "科目1", response.Subject.Name)
				assert.Equal(t, api.H1, response.Subject.Semester)
				assert.Equal(t, "syllabus-1", response.Subject.SyllabusId)
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
			mockRepo := repository.NewMockSubjectRepository()
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.SubjectsV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectsV1Create(t *testing.T) {
	classA := api.A

	tests := []struct {
		name               string
		request            api.SubjectServiceSubjectRequest
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目を作成できる",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "新しい科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H1,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes: []api.SubjectServiceSubjectTargetClass{
					{
						Grade: api.B1,
						Class: &classA,
					},
				},
				Requirements: []api.SubjectServiceSubjectRequirementRequest{
					{
						CourseId:        "course-1",
						RequirementType: api.Required,
					},
				},
				CategoryIds: []string{"category-1"},
				SyllabusId:  "syllabus-new",
			},
			withAdminClaim: true,
			wantCode:       http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Subject.Id)
				assert.Equal(t, "新しい科目", response.Subject.Name)
				assert.Equal(t, api.H1, response.Subject.Semester)
				assert.Equal(t, "syllabus-new", response.Subject.SyllabusId)
			},
		},
		{
			name: "developerクレームのみでも作成できる",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "developer経由の科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H2,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes: []api.SubjectServiceSubjectTargetClass{
					{
						Grade: api.B2,
						Class: &classA,
					},
				},
				Requirements: []api.SubjectServiceSubjectRequirementRequest{
					{
						CourseId:        "course-1",
						RequirementType: api.Optional,
					},
				},
				CategoryIds: []string{"category-1"},
				SyllabusId:  "syllabus-dev",
			},
			withDeveloperClaim: true,
			wantCode:           http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.Subject.Id)
				assert.Equal(t, "developer経由の科目", response.Subject.Name)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "新しい科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H1,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{},
				CategoryIds:               []string{},
				SyllabusId:                "syllabus-new",
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
			request: api.SubjectServiceSubjectRequest{
				Name:                      "新しい科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H1,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{},
				CategoryIds:               []string{},
				SyllabusId:                "syllabus-new",
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
			mockRepo := repository.NewMockSubjectRepository()
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)
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
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/subjects", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.SubjectsV1Create(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectsV1Update(t *testing.T) {
	classB := api.B

	tests := []struct {
		name           string
		id             string
		request        api.SubjectServiceSubjectRequest
		withAdminClaim bool
		customClaims   map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目を更新できる",
			id:   "1",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "更新された科目",
				FacultyId:                 "faculty-2",
				Semester:                  api.H2,
				DayOfWeekTimetableSlotIds: []string{"slot-2"},
				EligibleAttributes: []api.SubjectServiceSubjectTargetClass{
					{
						Grade: api.B2,
						Class: &classB,
					},
				},
				Requirements: []api.SubjectServiceSubjectRequirementRequest{
					{
						CourseId:        "course-2",
						RequirementType: api.OptionalRequired,
					},
				},
				CategoryIds: []string{"category-2"},
				SyllabusId:  "syllabus-updated",
			},
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.Subject.Id)
				assert.Equal(t, "更新された科目", response.Subject.Name)
				assert.Equal(t, api.H2, response.Subject.Semester)
				assert.Equal(t, "syllabus-updated", response.Subject.SyllabusId)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "更新された科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H1,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{},
				CategoryIds:               []string{},
				SyllabusId:                "syllabus-updated",
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
			request: api.SubjectServiceSubjectRequest{
				Name:                      "更新された科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H1,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{},
				CategoryIds:               []string{},
				SyllabusId:                "syllabus-updated",
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
			mockRepo := repository.NewMockSubjectRepository()
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)
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
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/subjects/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.SubjectsV1Update(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestSubjectsV1Delete(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{} // 指定時はこのクレームでトークンをセット（403検証用）
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に科目を削除できる",
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
			mockRepo := repository.NewMockSubjectRepository()
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.SubjectsV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
