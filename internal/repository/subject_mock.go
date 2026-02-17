package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type mockSubjectRepository struct{}

// NewMockSubjectRepository モックリポジトリを作成する
func NewMockSubjectRepository() service.SubjectRepository {
	return &mockSubjectRepository{}
}

// List 一覧を取得する（モック）
func (r *mockSubjectRepository) List(ctx context.Context) ([]domain.Subject, error) {
	classA := domain.ClassA
	return []domain.Subject{
		{
			ID:   "1",
			Name: "科目1",
			Faculty: domain.Faculty{
				ID:    "faculty-1",
				Name:  "教員1",
				Email: "faculty1@example.com",
			},
			Semester: domain.CourseSemesterH1,
			DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{
				{
					ID:            "slot-1",
					DayOfWeek:     domain.DayOfWeekMonday,
					TimetableSlot: domain.TimetableSlotSlot1,
				},
			},
			EligibleAttributes: []domain.SubjectTargetClass{
				{
					Grade: domain.GradeB1,
					Class: &classA,
				},
			},
			Requirements: []domain.SubjectRequirement{
				{
					Course: domain.Course{
						ID:   "course-1",
						Name: "コース1",
					},
					RequirementType: domain.SubjectRequirementTypeRequired,
				},
			},
			Categories: []domain.SubjectCategory{
				{
					ID:   "category-1",
					Name: "カテゴリ1",
				},
			},
			SyllabusID: "syllabus-1",
		},
	}, nil
}

// Detail 詳細を取得する（モック）
func (r *mockSubjectRepository) Detail(ctx context.Context, id string) (*domain.Subject, error) {
	classA := domain.ClassA
	return &domain.Subject{
		ID:   id,
		Name: "科目" + id,
		Faculty: domain.Faculty{
			ID:    "faculty-1",
			Name:  "教員1",
			Email: "faculty1@example.com",
		},
		Semester: domain.CourseSemesterH1,
		DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{
			{
				ID:            "slot-1",
				DayOfWeek:     domain.DayOfWeekMonday,
				TimetableSlot: domain.TimetableSlotSlot1,
			},
		},
		EligibleAttributes: []domain.SubjectTargetClass{
			{
				Grade: domain.GradeB1,
				Class: &classA,
			},
		},
		Requirements: []domain.SubjectRequirement{
			{
				Course: domain.Course{
					ID:   "course-1",
					Name: "コース1",
				},
				RequirementType: domain.SubjectRequirementTypeRequired,
			},
		},
		Categories: []domain.SubjectCategory{
			{
				ID:   "category-1",
				Name: "カテゴリ1",
			},
		},
		SyllabusID: "syllabus-1",
	}, nil
}

// Create 新規作成する（モック）
func (r *mockSubjectRepository) Create(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error) {
	return &domain.Subject{
		ID:   "created-id",
		Name: req.Name,
		Faculty: domain.Faculty{
			ID:    req.FacultyID,
			Name:  "教員1",
			Email: "faculty1@example.com",
		},
		Semester: req.Semester,
		DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{
			{
				ID:            "slot-1",
				DayOfWeek:     domain.DayOfWeekMonday,
				TimetableSlot: domain.TimetableSlotSlot1,
			},
		},
		EligibleAttributes: req.EligibleAttributes,
		Requirements: []domain.SubjectRequirement{
			{
				Course: domain.Course{
					ID:   "course-1",
					Name: "コース1",
				},
				RequirementType: domain.SubjectRequirementTypeRequired,
			},
		},
		Categories: []domain.SubjectCategory{
			{
				ID:   "category-1",
				Name: "カテゴリ1",
			},
		},
		SyllabusID: req.SyllabusID,
	}, nil
}

// Update 更新する（モック）
func (r *mockSubjectRepository) Update(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error) {
	return &domain.Subject{
		ID:   id,
		Name: req.Name,
		Faculty: domain.Faculty{
			ID:    req.FacultyID,
			Name:  "教員1",
			Email: "faculty1@example.com",
		},
		Semester: req.Semester,
		DayOfWeekTimetableSlots: []domain.DayOfWeekTimetableSlot{
			{
				ID:            "slot-1",
				DayOfWeek:     domain.DayOfWeekMonday,
				TimetableSlot: domain.TimetableSlotSlot1,
			},
		},
		EligibleAttributes: req.EligibleAttributes,
		Requirements: []domain.SubjectRequirement{
			{
				Course: domain.Course{
					ID:   "course-1",
					Name: "コース1",
				},
				RequirementType: domain.SubjectRequirementTypeRequired,
			},
		},
		Categories: []domain.SubjectCategory{
			{
				ID:   "category-1",
				Name: "カテゴリ1",
			},
		},
		SyllabusID: req.SyllabusID,
	}, nil
}

// Delete 削除する（モック）
func (r *mockSubjectRepository) Delete(ctx context.Context, id string) error {
	return nil
}
