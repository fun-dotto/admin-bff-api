package external

import (
	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToDomainSubjectCategory 外部API形式からドメイン形式に変換する
func ToDomainSubjectCategory(c subject_api.SubjectCategory) domain.SubjectCategory {
	return domain.SubjectCategory{
		ID:   domain.SubjectCategoryID(c.Id),
		Name: c.Name,
	}
}

// ToExternalSubjectCategoryRequest ドメイン形式から外部API形式に変換する
func ToExternalSubjectCategoryRequest(req *domain.SubjectCategoryRequest) subject_api.SubjectCategoryRequest {
	return subject_api.SubjectCategoryRequest{
		Name: req.Name,
	}
}
