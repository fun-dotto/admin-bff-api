package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/external"
	"github.com/fun-dotto/api-template/internal/service"
)

type subjectCategoryRepository struct {
	client *subject_api.ClientWithResponses
}

// NewSubjectCategoryRepository 新規作成する
func NewSubjectCategoryRepository(client *subject_api.ClientWithResponses) service.SubjectCategoryRepository {
	return &subjectCategoryRepository{client: client}
}

// List 一覧を取得する
func (r *subjectCategoryRepository) List(ctx context.Context) ([]domain.SubjectCategory, error) {
	response, err := r.client.SubjectCategoriesV1ListWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject categories: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subject categories: status %d", response.StatusCode())
	}

	result := make([]domain.SubjectCategory, len(response.JSON200.SubjectCategories))
	for i, c := range response.JSON200.SubjectCategories {
		result[i] = external.ToDomainSubjectCategory(c)
	}

	return result, nil
}

// Detail 詳細を取得する
func (r *subjectCategoryRepository) Detail(ctx context.Context, id string) (*domain.SubjectCategory, error) {
	response, err := r.client.SubjectCategoriesV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject category: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subject category: status %d", response.StatusCode())
	}

	result := external.ToDomainSubjectCategory(response.JSON200.SubjectCategory)
	return &result, nil
}

// Create 新規作成する
func (r *subjectCategoryRepository) Create(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	body := external.ToExternalSubjectCategoryRequest(req)

	response, err := r.client.SubjectCategoriesV1CreateWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create subject category: %w", err)
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create subject category: status %d", response.StatusCode())
	}

	result := external.ToDomainSubjectCategory(response.JSON201.SubjectCategory)
	return &result, nil
}

// Update 更新する
func (r *subjectCategoryRepository) Update(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	body := external.ToExternalSubjectCategoryRequest(req)

	response, err := r.client.SubjectCategoriesV1UpdateWithResponse(ctx, id, body)
	if err != nil {
		return nil, fmt.Errorf("failed to update subject category: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update subject category: status %d", response.StatusCode())
	}

	result := external.ToDomainSubjectCategory(response.JSON200.SubjectCategory)
	return &result, nil
}

// Delete 削除する
func (r *subjectCategoryRepository) Delete(ctx context.Context, id string) error {
	response, err := r.client.SubjectCategoriesV1DeleteWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete subject category: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete subject category: status %d", response.StatusCode())
	}

	return nil
}
