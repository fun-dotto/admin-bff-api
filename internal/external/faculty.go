package external

import (
	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToDomainFaculty 外部API形式からドメイン形式に変換する
func ToDomainFaculty(f subject_api.Faculty) domain.Faculty {
	return domain.Faculty{
		ID:    f.Id,
		Name:  f.Name,
		Email: f.Email,
	}
}

// ToExternalFacultyRequest ドメイン形式から外部API形式に変換する
func ToExternalFacultyRequest(req *domain.FacultyRequest) subject_api.FacultyRequest {
	return subject_api.FacultyRequest{
		Name:  req.Name,
		Email: req.Email,
	}
}
