package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/external"
	"github.com/fun-dotto/api-template/internal/service"
)

type subjectRepository struct {
	client *subject_api.ClientWithResponses
}

// NewSubjectRepository 新規作成する
func NewSubjectRepository(client *subject_api.ClientWithResponses) service.SubjectRepository {
	return &subjectRepository{client: client}
}

// List 一覧を取得する
func (r *subjectRepository) List(ctx context.Context) ([]domain.Subject, error) {
	response, err := r.client.SubjectsV1ListWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subjects: status %d", response.StatusCode())
	}

	result := make([]domain.Subject, len(response.JSON200.Subjects))
	for i, s := range response.JSON200.Subjects {
		result[i] = external.ToDomainSubject(s)
	}

	return result, nil
}

// Detail 詳細を取得する
func (r *subjectRepository) Detail(ctx context.Context, id string) (*domain.Subject, error) {
	response, err := r.client.SubjectsV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subject: status %d", response.StatusCode())
	}

	result := external.ToDomainSubject(response.JSON200.Subject)
	return &result, nil
}

// Create 新規作成する
func (r *subjectRepository) Create(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error) {
	body := external.ToExternalSubjectRequest(req)

	response, err := r.client.SubjectsV1CreateWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create subject: %w", err)
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create subject: status %d", response.StatusCode())
	}

	result := external.ToDomainSubject(response.JSON201.Subject)
	return &result, nil
}

// Update 更新する
func (r *subjectRepository) Update(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error) {
	body := external.ToExternalSubjectRequest(req)

	response, err := r.client.SubjectsV1UpdateWithResponse(ctx, id, body)
	if err != nil {
		return nil, fmt.Errorf("failed to update subject: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update subject: status %d", response.StatusCode())
	}

	result := external.ToDomainSubject(response.JSON200.Subject)
	return &result, nil
}

// Delete 削除する
func (r *subjectRepository) Delete(ctx context.Context, id string) error {
	response, err := r.client.SubjectsV1DeleteWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete subject: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete subject: status %d", response.StatusCode())
	}

	return nil
}
