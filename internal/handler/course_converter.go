package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPICourse ドメインモデルをAPIモデルに変換する
func ToAPICourse(c *domain.Course) api.SubjectServiceCourse {
	return api.SubjectServiceCourse{
		Id:   c.ID.String(),
		Name: c.Name,
	}
}

// ToAPICourses ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPICourses(courses []domain.Course) []api.SubjectServiceCourse {
	result := make([]api.SubjectServiceCourse, len(courses))
	for i, c := range courses {
		result[i] = ToAPICourse(&c)
	}
	return result
}

// ToDomainCourseRequest APIモデルをドメインモデルに変換する
func ToDomainCourseRequest(req *api.SubjectServiceCourseRequest) *domain.CourseRequest {
	return &domain.CourseRequest{
		Name: req.Name,
	}
}
