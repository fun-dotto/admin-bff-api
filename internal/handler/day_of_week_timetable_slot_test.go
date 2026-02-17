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

func TestDayOfWeekTimetableSlotsV1List(t *testing.T) {
	tests := []struct {
		name               string
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に曜日・時限一覧が取得できる",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")

				slots, ok := response["dayOfWeekTimetableSlots"].([]interface{})
				assert.True(t, ok, "dayOfWeekTimetableSlotsフィールドが配列ではありません")
				assert.NotEmpty(t, slots, "曜日・時限が空です")
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
				slots, ok := response["dayOfWeekTimetableSlots"].([]interface{})
				assert.True(t, ok, "dayOfWeekTimetableSlotsフィールドが配列ではありません")
				assert.NotEmpty(t, slots, "曜日・時限が空です")
				assert.Len(t, slots, 1, "MockRepositoryは1件返すはずです")
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
			name:           "曜日・時限のフィールドが正しく返される",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					DayOfWeekTimetableSlots []api.SubjectServiceDayOfWeekTimetableSlot `json:"dayOfWeekTimetableSlots"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.DayOfWeekTimetableSlots, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", response.DayOfWeekTimetableSlots[0].Id)
				assert.Equal(t, api.Monday, response.DayOfWeekTimetableSlots[0].DayOfWeek)
				assert.Equal(t, api.Slot1, response.DayOfWeekTimetableSlots[0].TimetableSlot)
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
			mockRepo := repository.NewMockDayOfWeekTimetableSlotRepository()
			slotService := service.NewDayOfWeekTimetableSlotService(mockRepo)
			h := handler.NewHandler().WithDayOfWeekTimetableSlotService(slotService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.DayOfWeekTimetableSlotsV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestDayOfWeekTimetableSlotsV1Detail(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "正常に曜日・時限詳細が取得できる",
			id:             "1",
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					DayOfWeekTimetableSlot api.SubjectServiceDayOfWeekTimetableSlot `json:"dayOfWeekTimetableSlot"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.DayOfWeekTimetableSlot.Id)
				assert.Equal(t, api.Monday, response.DayOfWeekTimetableSlot.DayOfWeek)
				assert.Equal(t, api.Slot1, response.DayOfWeekTimetableSlot.TimetableSlot)
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
			mockRepo := repository.NewMockDayOfWeekTimetableSlotRepository()
			slotService := service.NewDayOfWeekTimetableSlotService(mockRepo)
			h := handler.NewHandler().WithDayOfWeekTimetableSlotService(slotService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.DayOfWeekTimetableSlotsV1Detail(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestDayOfWeekTimetableSlotsV1Create(t *testing.T) {
	tests := []struct {
		name               string
		request            api.SubjectServiceDayOfWeekTimetableSlotRequest
		withAdminClaim     bool
		withDeveloperClaim bool
		customClaims       map[string]interface{}
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に曜日・時限を作成できる",
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Tuesday,
				TimetableSlot: api.Slot2,
			},
			withAdminClaim: true,
			wantCode:       http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					DayOfWeekTimetableSlot api.SubjectServiceDayOfWeekTimetableSlot `json:"dayOfWeekTimetableSlot"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.DayOfWeekTimetableSlot.Id)
				assert.Equal(t, api.Tuesday, response.DayOfWeekTimetableSlot.DayOfWeek)
				assert.Equal(t, api.Slot2, response.DayOfWeekTimetableSlot.TimetableSlot)
			},
		},
		{
			name: "developerクレームのみでも作成できる",
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Wednesday,
				TimetableSlot: api.Slot3,
			},
			withDeveloperClaim: true,
			wantCode:           http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					DayOfWeekTimetableSlot api.SubjectServiceDayOfWeekTimetableSlot `json:"dayOfWeekTimetableSlot"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "created-id", response.DayOfWeekTimetableSlot.Id)
				assert.Equal(t, api.Wednesday, response.DayOfWeekTimetableSlot.DayOfWeek)
				assert.Equal(t, api.Slot3, response.DayOfWeekTimetableSlot.TimetableSlot)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Monday,
				TimetableSlot: api.Slot1,
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
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Monday,
				TimetableSlot: api.Slot1,
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
			mockRepo := repository.NewMockDayOfWeekTimetableSlotRepository()
			slotService := service.NewDayOfWeekTimetableSlotService(mockRepo)
			h := handler.NewHandler().WithDayOfWeekTimetableSlotService(slotService)
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
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/day-of-week-timetable-slots", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.DayOfWeekTimetableSlotsV1Create(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestDayOfWeekTimetableSlotsV1Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		request        api.SubjectServiceDayOfWeekTimetableSlotRequest
		withAdminClaim bool
		customClaims   map[string]interface{}
		wantCode       int
		validate       func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に曜日・時限を更新できる",
			id:   "1",
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Friday,
				TimetableSlot: api.Slot5,
			},
			withAdminClaim: true,
			wantCode:       http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					DayOfWeekTimetableSlot api.SubjectServiceDayOfWeekTimetableSlot `json:"dayOfWeekTimetableSlot"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.Equal(t, "1", response.DayOfWeekTimetableSlot.Id)
				assert.Equal(t, api.Friday, response.DayOfWeekTimetableSlot.DayOfWeek)
				assert.Equal(t, api.Slot5, response.DayOfWeekTimetableSlot.TimetableSlot)
			},
		},
		{
			name: "認証トークンがない場合は401エラー",
			id:   "1",
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Monday,
				TimetableSlot: api.Slot1,
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
			request: api.SubjectServiceDayOfWeekTimetableSlotRequest{
				DayOfWeek:     api.Monday,
				TimetableSlot: api.Slot1,
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
			mockRepo := repository.NewMockDayOfWeekTimetableSlotRepository()
			slotService := service.NewDayOfWeekTimetableSlotService(mockRepo)
			h := handler.NewHandler().WithDayOfWeekTimetableSlotService(slotService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "リクエストボディのJSONエンコードに失敗しました")
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/day-of-week-timetable-slots/"+tt.id, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.DayOfWeekTimetableSlotsV1Update(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestDayOfWeekTimetableSlotsV1Delete(t *testing.T) {
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
			name:           "正常に曜日・時限を削除できる",
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
			mockRepo := repository.NewMockDayOfWeekTimetableSlotRepository()
			slotService := service.NewDayOfWeekTimetableSlotService(mockRepo)
			h := handler.NewHandler().WithDayOfWeekTimetableSlotService(slotService)
			var w *httptest.ResponseRecorder
			var c *gin.Context
			if tt.customClaims != nil {
				w, c = setupTestContextWithClaims(tt.customClaims)
			} else if tt.withDeveloperClaim {
				w, c = setupTestContextWithClaims(map[string]interface{}{"developer": true})
			} else {
				w, c = setupTestContext(tt.withAdminClaim)
			}

			h.DayOfWeekTimetableSlotsV1Delete(c, tt.id)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
