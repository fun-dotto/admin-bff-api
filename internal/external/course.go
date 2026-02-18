package external

import (
	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToDomainCourse 外部API形式からドメイン形式に変換する
func ToDomainCourse(c subject_api.Course) domain.Course {
	return domain.Course{
		ID:   domain.CourseID(c.Id),
		Name: c.Name,
	}
}

// ToExternalCourseRequest ドメイン形式から外部API形式に変換する
func ToExternalCourseRequest(req *domain.CourseRequest) subject_api.CourseRequest {
	return subject_api.CourseRequest{
		Name: req.Name,
	}
}
