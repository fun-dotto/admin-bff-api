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

func createMockSubject(id string) *domain.Subject {
	classA := domain.ClassA
	return &domain.Subject{
		ID:   domain.SubjectID(id),
		Name: "科目" + id,
		Faculty: domain.Faculty{
			ID:    domain.FacultyID("faculty-1"),
			Name:  "教員1",
			Email: "faculty1@example.com",
		},
		Semester: domain.CourseSemesterH1,
		DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{
			{ID: "slot-1", DayOfWeek: domain.DayOfWeekMonday, TimetableSlot: domain.TimetableSlotSlot1},
		},
		EligibleAttributes: []domain.SubjectTargetClass{
			{Grade: domain.GradeB1, Class: &classA},
		},
		Requirements: []domain.SubjectRequirement{
			{Course: domain.Course{ID: domain.CourseID("course-1"), Name: "コース1"}, RequirementType: domain.SubjectRequirementTypeRequired},
		},
		Categories: []domain.SubjectCategory{
			{ID: "category-1", Name: "カテゴリ1"},
		},
		SyllabusID: "syllabus-1",
	}
}

func TestSubjectsV1List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupMock    func() *repository.MockSubjectRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目一覧が取得できる",
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
					ListFunc: func(ctx context.Context) ([]domain.Subject, error) {
						return []domain.Subject{*createMockSubject("1")}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subjects []api.SubjectServiceSubject `json:"subjects"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Subjects, 1)
				assert.Equal(t, "1", response.Subjects[0].Id)
				assert.Equal(t, "科目1", response.Subjects[0].Name)
				assert.Equal(t, api.H1, response.Subjects[0].Semester)
				assert.Equal(t, "syllabus-1", response.Subjects[0].SyllabusId)
			},
		},
		{
			name: "developerクレームのみでも一覧が取得できる",
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
					ListFunc: func(ctx context.Context) ([]domain.Subject, error) {
						return []domain.Subject{*createMockSubject("1")}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subjects []api.SubjectServiceSubject `json:"subjects"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Subjects, 1)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
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
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockSubjectRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目詳細が取得できる",
			id:   "1",
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
					DetailFunc: func(ctx context.Context, id string) (*domain.Subject, error) {
						return createMockSubject(id), nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "1", response.Subject.Id)
				assert.Equal(t, "科目1", response.Subject.Name)
				assert.Equal(t, api.H1, response.Subject.Semester)
				assert.Equal(t, "syllabus-1", response.Subject.SyllabusId)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
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
	gin.SetMode(gin.TestMode)
	classA := api.A

	tests := []struct {
		name         string
		request      api.SubjectServiceSubjectRequest
		setupMock    func() *repository.MockSubjectRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目を作成できる",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "新しい科目",
				FacultyId:                 "faculty-1",
				Semester:                  api.H1,
				DayOfWeekTimetableSlotIds: []string{"slot-1"},
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{{Grade: api.B1, Class: &classA}},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{{CourseId: "course-1", RequirementType: api.Required}},
				CategoryIds:               []string{"category-1"},
				SyllabusId:                "syllabus-new",
			},
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
					CreateFunc: func(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error) {
						classA := domain.ClassA
						return &domain.Subject{
							ID:                      domain.SubjectID("created-id"),
							Name:                    req.Name,
							Faculty:                 domain.Faculty{ID: domain.FacultyID(req.FacultyID), Name: "教員1", Email: "faculty1@example.com"},
							Semester:                req.Semester,
							DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{{ID: domain.DayOfWeekTimetableSlotID("slot-1"), DayOfWeek: domain.DayOfWeekMonday, TimetableSlot: domain.TimetableSlotSlot1}},
							EligibleAttributes:      []domain.SubjectTargetClass{{Grade: domain.GradeB1, Class: &classA}},
							Requirements:            []domain.SubjectRequirement{{Course: domain.Course{ID: domain.CourseID("course-1"), Name: "コース1"}, RequirementType: domain.SubjectRequirementTypeRequired}},
							Categories:              []domain.SubjectCategory{{ID: domain.SubjectCategoryID("category-1"), Name: "カテゴリ1"}},
							SyllabusID:              req.SyllabusID,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
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
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{{Grade: api.B2, Class: &classA}},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{{CourseId: "course-1", RequirementType: api.Optional}},
				CategoryIds:               []string{"category-1"},
				SyllabusId:                "syllabus-dev",
			},
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
					CreateFunc: func(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error) {
						classA := domain.ClassA
						return &domain.Subject{
							ID:                      domain.SubjectID("created-id"),
							Name:                    req.Name,
							Faculty:                 domain.Faculty{ID: domain.FacultyID(req.FacultyID), Name: "教員1", Email: "faculty1@example.com"},
							Semester:                req.Semester,
							DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{{ID: domain.DayOfWeekTimetableSlotID("slot-1"), DayOfWeek: domain.DayOfWeekMonday, TimetableSlot: domain.TimetableSlotSlot1}},
							EligibleAttributes:      []domain.SubjectTargetClass{{Grade: domain.GradeB2, Class: &classA}},
							Requirements:            []domain.SubjectRequirement{{Course: domain.Course{ID: domain.CourseID("course-1"), Name: "コース1"}, RequirementType: domain.SubjectRequirementTypeOptional}},
							Categories:              []domain.SubjectCategory{{ID: domain.SubjectCategoryID("category-1"), Name: "カテゴリ1"}},
							SyllabusID:              req.SyllabusID,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"developer": true},
			wantCode:     http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
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
	gin.SetMode(gin.TestMode)
	classB := api.B

	tests := []struct {
		name         string
		id           string
		request      api.SubjectServiceSubjectRequest
		setupMock    func() *repository.MockSubjectRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目を更新できる",
			id:   "1",
			request: api.SubjectServiceSubjectRequest{
				Name:                      "更新された科目",
				FacultyId:                 "faculty-2",
				Semester:                  api.H2,
				DayOfWeekTimetableSlotIds: []string{"slot-2"},
				EligibleAttributes:        []api.SubjectServiceSubjectTargetClass{{Grade: api.B2, Class: &classB}},
				Requirements:              []api.SubjectServiceSubjectRequirementRequest{{CourseId: "course-2", RequirementType: api.OptionalRequired}},
				CategoryIds:               []string{"category-2"},
				SyllabusId:                "syllabus-updated",
			},
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
					UpdateFunc: func(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error) {
						classB := domain.ClassB
						return &domain.Subject{
							ID:                      domain.SubjectID(id),
							Name:                    req.Name,
							Faculty:                 domain.Faculty{ID: domain.FacultyID(req.FacultyID), Name: "教員2", Email: "faculty2@example.com"},
							Semester:                req.Semester,
							DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{{ID: domain.DayOfWeekTimetableSlotID("slot-2"), DayOfWeek: domain.DayOfWeekTuesday, TimetableSlot: domain.TimetableSlotSlot2}},
							EligibleAttributes:      []domain.SubjectTargetClass{{Grade: domain.GradeB2, Class: &classB}},
							Requirements:            []domain.SubjectRequirement{{Course: domain.Course{ID: domain.CourseID("course-2"), Name: "コース2"}, RequirementType: domain.SubjectRequirementTypeOptionalRequired}},
							Categories:              []domain.SubjectCategory{{ID: domain.SubjectCategoryID("category-2"), Name: "カテゴリ2"}},
							SyllabusID:              req.SyllabusID,
						}, nil
					},
				}
			},
			customClaims: map[string]interface{}{"admin": true},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Subject api.SubjectServiceSubject `json:"subject"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)
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
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		setupMock    func() *repository.MockSubjectRepository
		customClaims map[string]interface{}
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に科目を削除できる",
			id:   "1",
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			setupMock: func() *repository.MockSubjectRepository {
				return &repository.MockSubjectRepository{}
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
			subjectService := service.NewSubjectService(mockRepo)
			h := handler.NewHandler().WithSubjectService(subjectService)

			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(false)
			}

			h.SubjectsV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
