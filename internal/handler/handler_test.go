package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	firebaseauth "firebase.google.com/go/v4/auth"
	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/admin-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestSubjectsV1List_UsesRenamedQueryParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"subjects":[]}`))
	}))
	defer server.Close()

	h := newTestHandler(t, server.URL)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/subjects", nil)
	setAdminClaim(c)

	q := "math"
	year := 2026
	grades := []api.DottoFoundationV1Grade{api.B1}
	courses := []api.DottoFoundationV1Course{api.InformationSystem}
	classes := []api.DottoFoundationV1Class{api.A}
	classifications := []api.DottoFoundationV1SubjectClassification{api.Cultural}
	semesters := []api.DottoFoundationV1CourseSemester{api.Q1, api.Q2}
	requirementTypes := []api.DottoFoundationV1SubjectRequirementType{api.Optional}
	categories := []api.DottoFoundationV1CulturalSubjectCategory{api.Society}

	h.SubjectsV1List(c, api.SubjectsV1ListParams{
		Q:                         &q,
		Grades:                    &grades,
		Courses:                   &courses,
		Classes:                   &classes,
		Classifications:           &classifications,
		Year:                      &year,
		Semesters:                 &semesters,
		RequirementTypes:          &requirementTypes,
		CulturalSubjectCategories: &categories,
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if gotQuery.Get("grades") == "" || gotQuery.Get("grade") != "" {
		t.Fatalf("query grades = %q, grade = %q", gotQuery.Get("grades"), gotQuery.Get("grade"))
	}
	if gotQuery.Get("classes") == "" || gotQuery.Get("class") != "" {
		t.Fatalf("query classes = %q, class = %q", gotQuery.Get("classes"), gotQuery.Get("class"))
	}
	if gotQuery.Get("classifications") == "" || gotQuery.Get("classification") != "" {
		t.Fatalf("query classifications = %q, classification = %q", gotQuery.Get("classifications"), gotQuery.Get("classification"))
	}
	if gotQuery.Get("semesters") == "" || gotQuery.Get("semester") != "" {
		t.Fatalf("query semesters = %q, semester = %q", gotQuery.Get("semesters"), gotQuery.Get("semester"))
	}
	if gotQuery.Get("requirementTypes") == "" || gotQuery.Get("requirementType") != "" {
		t.Fatalf("query requirementTypes = %q, requirementType = %q", gotQuery.Get("requirementTypes"), gotQuery.Get("requirementType"))
	}
	if gotQuery.Get("culturalSubjectCategories") == "" || gotQuery.Get("culturalSubjectCategory") != "" {
		t.Fatalf("query culturalSubjectCategories = %q, culturalSubjectCategory = %q", gotQuery.Get("culturalSubjectCategories"), gotQuery.Get("culturalSubjectCategory"))
	}
}

func TestPersonalCalendarItemsV1List_ProxiesAcademicAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	wantDate := time.Date(2026, 3, 26, 9, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("userId"); got != "user-1" {
			t.Fatalf("userId = %q, want %q", got, "user-1")
		}
		if got := query.Get("dates"); got == "" {
			t.Fatal("dates query parameter is empty")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"personalCalendarItems":[{"date":"2026-03-26T09:00:00Z","slot":{"dayOfWeek":"Thursday","period":"first"},"timetableItem":{"id":"item-1","subject":{"id":"subject-1","name":"Algorithms"},"slot":{"dayOfWeek":"Thursday","period":"first"},"rooms":[]}}]}`))
	}))
	defer server.Close()

	h := newTestHandler(t, server.URL)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/personalCalendarItems", nil)
	setAdminClaim(c)

	h.PersonalCalendarItemsV1List(c, api.PersonalCalendarItemsV1ListParams{
		UserId: "user-1",
		Dates:  []time.Time{wantDate},
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		PersonalCalendarItems []struct {
			Date string `json:"date"`
		} `json:"personalCalendarItems"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.PersonalCalendarItems) != 1 || body.PersonalCalendarItems[0].Date != "2026-03-26T09:00:00Z" {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
}

func newTestHandler(t *testing.T, baseURL string) *Handler {
	t.Helper()

	academicClient, err := academic_api.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatalf("new academic client: %v", err)
	}
	announcementClient, err := announcement_api.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatalf("new announcement client: %v", err)
	}

	return NewHandler(academicClient, announcementClient)
}

func setAdminClaim(c *gin.Context) {
	c.Set(middleware.FirebaseTokenContextKey, &firebaseauth.Token{
		Claims: map[string]interface{}{
			"admin": true,
		},
	})
}
