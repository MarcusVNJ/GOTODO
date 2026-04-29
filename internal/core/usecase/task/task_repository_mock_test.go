package usecase_test

import (
	"context"

	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Save(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) ExistByID(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context, statusFilter string, minPriority int) ([]*models.Task, error) {
	args := m.Called(ctx, statusFilter, minPriority)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByStatus(ctx context.Context, status enums.Status) ([]*models.Task, error) {
	args := m.Called(ctx, status)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}
