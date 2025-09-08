package repository

import (
	"database/sql"
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	createTableSQL := `
    CREATE TABLE products (
        id TEXT NOT NULL PRIMARY KEY,
        name TEXT,
        price REAL,
        quantity INTEGER
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create products table: %v", err)
	}

	return db
}

func TestSqliteRepository_SaveAndFindById(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	product, _ := domain.CreateNewProduct("Test Keyboard", 99.99, 50)

	err := repo.Save(product)
	if err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	foundProduct, err := repo.FindById(product.Id)
	if err != nil {
		t.Fatalf("FindById() returned an unexpected error: %v", err)
	}

	if product.Id != foundProduct.Id || product.Name != foundProduct.Name || product.Price != foundProduct.Price || product.Quantity != foundProduct.Quantity {
		t.Errorf("FindById() got = %+v, want %+v", foundProduct, product)
	}
}

func TestSqliteRepository_FindById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteRepository(db)

	_, err := repo.FindById("non-existent-id")
	if err == nil {
		t.Fatal("FindById() expected an error for non-existent product, but got nil")
	}
}

func TestSqliteRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	product, _ := domain.CreateNewProduct("Old Name", 10.0, 10)
	repo.Save(product)

	product.Name = "New Updated Name"
	product.Price = 25.50
	product.Quantity = 100

	err := repo.Update(product)
	if err != nil {
		t.Fatalf("Update() returned an unexpected error: %v", err)
	}

	updatedProduct, err := repo.FindById(product.Id)
	if err != nil {
		t.Fatalf("FindById() after update returned an unexpected error: %v", err)
	}

	if product.Id != updatedProduct.Id || product.Name != updatedProduct.Name || product.Price != updatedProduct.Price || product.Quantity != updatedProduct.Quantity {
		t.Errorf("Update() failed. got = %+v, want %+v", updatedProduct, product)
	}
}

func TestSqliteRepository_DeleteById(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	product, _ := domain.CreateNewProduct("ToDelete", 1.0, 1)
	repo.Save(product)

	err := repo.DeleteById(product.Id)
	if err != nil {
		t.Fatalf("DeleteById() returned an unexpected error: %v", err)
	}

	_, err = repo.FindById(product.Id)
	if err == nil {
		t.Fatal("FindById() after delete should have returned an error, but got nil")
	}
}

func TestSqliteRepository_ListAll(t *testing.T) {
	t.Run("list_all_empty", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()
		repo := NewSQLiteRepository(db)

		products, err := repo.ListAll()
		if err != nil {
			t.Fatalf("ListAll() on empty table returned an error: %v", err)
		}
		if len(products) != 0 {
			t.Errorf("ListAll() on empty table should return 0 products, got %d", len(products))
		}
	})

	t.Run("list_all_with_products", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()
		repo := NewSQLiteRepository(db)

		p1, _ := domain.CreateNewProduct("Product 1", 10, 1)
		p2, _ := domain.CreateNewProduct("Product 2", 20, 2)
		repo.Save(p1)
		repo.Save(p2)

		products, err := repo.ListAll()
		if err != nil {
			t.Fatalf("ListAll() returned an error: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("ListAll() should return 2 products, got %d", len(products))
		}
	})
}
