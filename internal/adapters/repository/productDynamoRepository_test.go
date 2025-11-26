package repository

import (
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

type mockProductRepository struct {
	products map[string]*domain.Product
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (m *mockProductRepository) FindById(id string) (*domain.Product, error) {
	product, exists := m.products[id]
	if !exists {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (m *mockProductRepository) Save(product *domain.Product) error {
	m.products[product.Id] = product
	return nil
}

func (m *mockProductRepository) Update(product *domain.Product) error {
	if _, exists := m.products[product.Id]; !exists {
		return domain.ErrProductNotFound
	}
	m.products[product.Id] = product
	return nil
}

func (m *mockProductRepository) DeleteById(id string) error {
	if _, exists := m.products[id]; !exists {
		return domain.ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *mockProductRepository) ListAll() ([]domain.Product, error) {
	products := make([]domain.Product, 0, len(m.products))
	for _, product := range m.products {
		products = append(products, *product)
	}
	return products, nil
}

func TestMockProductRepository_FindById(t *testing.T) {
	repo := newMockProductRepository()
	product := &domain.Product{
		Id:       "test-id",
		Name:     "Test Product",
		Price:    99.99,
		Quantity: 10,
	}
	repo.Save(product)

	found, err := repo.FindById("test-id")
	assert.NoError(t, err)
	assert.Equal(t, product.Id, found.Id)
	assert.Equal(t, product.Name, found.Name)

	_, err = repo.FindById("non-existent")
	assert.Equal(t, domain.ErrProductNotFound, err)
}

func TestMockProductRepository_Save(t *testing.T) {
	repo := newMockProductRepository()
	product := &domain.Product{
		Id:       "test-id",
		Name:     "Test Product",
		Price:    99.99,
		Quantity: 10,
	}

	err := repo.Save(product)
	assert.NoError(t, err)

	saved, _ := repo.FindById("test-id")
	assert.Equal(t, product.Name, saved.Name)
}

func TestMockProductRepository_Update(t *testing.T) {
	repo := newMockProductRepository()
	product := &domain.Product{
		Id:       "test-id",
		Name:     "Test Product",
		Price:    99.99,
		Quantity: 10,
	}
	repo.Save(product)

	product.Price = 149.99
	err := repo.Update(product)
	assert.NoError(t, err)

	updated, _ := repo.FindById("test-id")
	assert.Equal(t, 149.99, updated.Price)

	nonExistent := &domain.Product{Id: "fake"}
	err = repo.Update(nonExistent)
	assert.Equal(t, domain.ErrProductNotFound, err)
}

func TestMockProductRepository_DeleteById(t *testing.T) {
	repo := newMockProductRepository()
	product := &domain.Product{
		Id:       "test-id",
		Name:     "Test Product",
		Price:    99.99,
		Quantity: 10,
	}
	repo.Save(product)

	err := repo.DeleteById("test-id")
	assert.NoError(t, err)

	_, err = repo.FindById("test-id")
	assert.Equal(t, domain.ErrProductNotFound, err)

	err = repo.DeleteById("non-existent")
	assert.Equal(t, domain.ErrProductNotFound, err)
}

func TestMockProductRepository_ListAll(t *testing.T) {
	repo := newMockProductRepository()

	products, err := repo.ListAll()
	assert.NoError(t, err)
	assert.Len(t, products, 0)

	repo.Save(&domain.Product{Id: "1", Name: "Product 1", Price: 10, Quantity: 5})
	repo.Save(&domain.Product{Id: "2", Name: "Product 2", Price: 20, Quantity: 10})

	products, err = repo.ListAll()
	assert.NoError(t, err)
	assert.Len(t, products, 2)
}
