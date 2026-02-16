package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/external"
	"github.com/fun-dotto/api-template/internal/service"
)

type courseRepository struct {
	client *subject_api.ClientWithResponses
}

// NewCourseRepository 新規作成する
func NewCourseRepository(client *subject_api.ClientWithResponses) service.CourseRepository {
	return &courseRepository{client: client}
}

// List 一覧を取得する
func (r *courseRepository) List(ctx context.Context) ([]domain.Course, error) {
	response, err := r.client.CoursesV1ListWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get courses: status %d", response.StatusCode())
	}

	result := make([]domain.Course, len(response.JSON200.Courses))
	for i, c := range response.JSON200.Courses {
		result[i] = external.ToDomainCourse(c)
	}

	return result, nil
}

// Detail 詳細を取得する
func (r *courseRepository) Detail(ctx context.Context, id string) (*domain.Course, error) {
	response, err := r.client.CoursesV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get course: status %d", response.StatusCode())
	}

	result := external.ToDomainCourse(response.JSON200.Course)
	return &result, nil
}

// Create 新規作成する
func (r *courseRepository) Create(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error) {
	body := external.ToExternalCourseRequest(req)

	response, err := r.client.CoursesV1CreateWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create course: status %d", response.StatusCode())
	}

	result := external.ToDomainCourse(response.JSON201.Course)
	return &result, nil
}

// Update 更新する
func (r *courseRepository) Update(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error) {
	body := external.ToExternalCourseRequest(req)

	response, err := r.client.CoursesV1UpdateWithResponse(ctx, id, body)
	if err != nil {
		return nil, fmt.Errorf("failed to update course: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update course: status %d", response.StatusCode())
	}

	result := external.ToDomainCourse(response.JSON200.Course)
	return &result, nil
}

// Delete 削除する
func (r *courseRepository) Delete(ctx context.Context, id string) error {
	response, err := r.client.CoursesV1DeleteWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete course: status %d", response.StatusCode())
	}

	return nil
}
