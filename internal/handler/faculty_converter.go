package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPIFaculty ドメインモデルをAPIモデルに変換する
func ToAPIFaculty(f *domain.Faculty) api.SubjectServiceFaculty {
	return api.SubjectServiceFaculty{
		Id:    f.ID.String(),
		Name:  f.Name,
		Email: f.Email,
	}
}

// ToAPIFaculties ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPIFaculties(faculties []domain.Faculty) []api.SubjectServiceFaculty {
	result := make([]api.SubjectServiceFaculty, len(faculties))
	for i, f := range faculties {
		result[i] = ToAPIFaculty(&f)
	}
	return result
}

// ToDomainFacultyRequest APIモデルをドメインモデルに変換する
func ToDomainFacultyRequest(req *api.SubjectServiceFacultyRequest) *domain.FacultyRequest {
	return &domain.FacultyRequest{
		Name:  req.Name,
		Email: req.Email,
	}
}
