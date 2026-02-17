package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPISubjectCategory ドメインモデルをAPIモデルに変換する
func ToAPISubjectCategory(c *domain.SubjectCategory) api.SubjectServiceSubjectCategory {
	return api.SubjectServiceSubjectCategory{
		Id:   c.ID,
		Name: c.Name,
	}
}

// ToAPISubjectCategories ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPISubjectCategories(categories []domain.SubjectCategory) []api.SubjectServiceSubjectCategory {
	result := make([]api.SubjectServiceSubjectCategory, len(categories))
	for i, c := range categories {
		result[i] = ToAPISubjectCategory(&c)
	}
	return result
}

// ToDomainSubjectCategoryRequest APIモデルをドメインモデルに変換する
func ToDomainSubjectCategoryRequest(req *api.SubjectServiceSubjectCategoryRequest) *domain.SubjectCategoryRequest {
	return &domain.SubjectCategoryRequest{
		Name: req.Name,
	}
}
