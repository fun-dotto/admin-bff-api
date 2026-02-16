package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/external"
	"github.com/fun-dotto/api-template/internal/service"
)

type facultyRepository struct {
	client *subject_api.ClientWithResponses
}

// NewFacultyRepository 新規作成する
func NewFacultyRepository(client *subject_api.ClientWithResponses) service.FacultyRepository {
	return &facultyRepository{client: client}
}

// List 一覧を取得する
func (r *facultyRepository) List(ctx context.Context) ([]domain.Faculty, error) {
	response, err := r.client.FacultiesV1ListWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get faculties: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculties: status %d", response.StatusCode())
	}

	result := make([]domain.Faculty, len(response.JSON200.Faculties))
	for i, f := range response.JSON200.Faculties {
		result[i] = external.ToDomainFaculty(f)
	}

	return result, nil
}

// Detail 詳細を取得する
func (r *facultyRepository) Detail(ctx context.Context, id string) (*domain.Faculty, error) {
	response, err := r.client.FacultiesV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculty: status %d", response.StatusCode())
	}

	result := external.ToDomainFaculty(response.JSON200.Faculty)
	return &result, nil
}

// Create 新規作成する
func (r *facultyRepository) Create(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error) {
	body := external.ToExternalFacultyRequest(req)

	response, err := r.client.FacultiesV1CreateWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create faculty: %w", err)
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create faculty: status %d", response.StatusCode())
	}

	result := external.ToDomainFaculty(response.JSON201.Faculty)
	return &result, nil
}

// Update 更新する
func (r *facultyRepository) Update(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error) {
	body := external.ToExternalFacultyRequest(req)

	response, err := r.client.FacultiesV1UpdateWithResponse(ctx, id, body)
	if err != nil {
		return nil, fmt.Errorf("failed to update faculty: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update faculty: status %d", response.StatusCode())
	}

	result := external.ToDomainFaculty(response.JSON200.Faculty)
	return &result, nil
}

// Delete 削除する
func (r *facultyRepository) Delete(ctx context.Context, id string) error {
	response, err := r.client.FacultiesV1DeleteWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete faculty: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete faculty: status %d", response.StatusCode())
	}

	return nil
}
