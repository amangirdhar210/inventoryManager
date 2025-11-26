package repository

import (
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

type mockManagerRepository struct {
	managers map[string]*domain.Manager
}

func newMockManagerRepository() *mockManagerRepository {
	return &mockManagerRepository{
		managers: make(map[string]*domain.Manager),
	}
}

func (m *mockManagerRepository) FindByEmail(email string) (*domain.Manager, error) {
	for _, manager := range m.managers {
		if manager.Email == email {
			return manager, nil
		}
	}
	return nil, domain.ErrManagerNotFound
}

func (m *mockManagerRepository) Save(manager *domain.Manager) error {
	m.managers[manager.Id] = manager
	return nil
}

func TestMockManagerRepository_FindByEmail(t *testing.T) {
	repo := newMockManagerRepository()
	manager := &domain.Manager{
		Id:       "test-id",
		Email:    "test@example.com",
		Password: "hashed-password",
	}
	repo.Save(manager)

	found, err := repo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, manager.Email, found.Email)

	_, err = repo.FindByEmail("nonexistent@example.com")
	assert.Equal(t, domain.ErrManagerNotFound, err)
}

func TestMockManagerRepository_Save(t *testing.T) {
	repo := newMockManagerRepository()
	manager := &domain.Manager{
		Id:       "test-id",
		Email:    "test@example.com",
		Password: "hashed-password",
	}

	err := repo.Save(manager)
	assert.NoError(t, err)

	saved, _ := repo.FindByEmail("test@example.com")
	assert.Equal(t, manager.Email, saved.Email)
}
