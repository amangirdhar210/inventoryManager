package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	productsTableSQL := `
    CREATE TABLE products (
        id TEXT NOT NULL PRIMARY KEY,
        name TEXT,
        price REAL,
        quantity INTEGER
    );`
	if _, err := db.Exec(productsTableSQL); err != nil {
		t.Fatalf("Failed to create products table: %v", err)
	}

	managersTableSQL := `
    CREATE TABLE managers (
        id TEXT NOT NULL PRIMARY KEY,
        email TEXT UNIQUE,
        password TEXT
    );`
	if _, err := db.Exec(managersTableSQL); err != nil {
		t.Fatalf("Failed to create managers table: %v", err)
	}

	return db
}

func TestSqliteRepository_SaveAndFindById(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteRepository(db)
	product, _ := domain.CreateNewProduct("Test Keyboard", 99.99, 50)

	if err := repo.Save(product); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	found, err := repo.FindById(product.Id)
	if err != nil {
		t.Fatalf("FindById() returned an unexpected error: %v", err)
	}

	if product.Id != found.Id || product.Name != found.Name {
		t.Errorf("FindById() got = %+v, want %+v", found, product)
	}
}

func TestSqliteRepository_FindById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteRepository(db)

	_, err := repo.FindById("non-existent-id")
	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("expected error %v, got %v", domain.ErrProductNotFound, err)
	}
}

func TestSqliteRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteRepository(db)
	product, _ := domain.CreateNewProduct("Old Name", 10.0, 10)
	repo.Save(product)

	product.Name = "New Name"
	product.Price = 25.50
	product.Quantity = 100

	if err := repo.Update(product); err != nil {
		t.Fatalf("Update() returned an unexpected error: %v", err)
	}

	updated, _ := repo.FindById(product.Id)
	if updated.Name != "New Name" || updated.Price != 25.50 || updated.Quantity != 100 {
		t.Errorf("Update() failed. got = %+v, want %+v", updated, product)
	}
}

func TestSqliteRepository_DeleteById(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteRepository(db)
	product, _ := domain.CreateNewProduct("ToDelete", 1.0, 1)
	repo.Save(product)

	t.Run("success", func(t *testing.T) {
		err := repo.DeleteById(product.Id)
		if err != nil {
			t.Fatalf("DeleteById() returned an unexpected error: %v", err)
		}

		_, err = repo.FindById(product.Id)
		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Errorf("expected error %v after delete, but got %v", domain.ErrProductNotFound, err)
		}
	})

	t.Run("fail_not_found", func(t *testing.T) {
		err := repo.DeleteById("non-existent-id")
		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Errorf("expected error %v for non-existent product, but got %v", domain.ErrProductNotFound, err)
		}
	})
}

func TestSqliteRepository_ListAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteRepository(db)

	t.Run("list_all_with_products", func(t *testing.T) {
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

	t.Run("list_all_empty", func(t *testing.T) {
		db.Exec("DELETE FROM products")
		products, err := repo.ListAll()
		if err != nil {
			t.Fatalf("ListAll() on empty table returned an error: %v", err)
		}
		if len(products) != 0 {
			t.Errorf("ListAll() on empty table should return 0 products, got %d", len(products))
		}
	})
}

func TestSqliteRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteRepository(db)

	manager := &domain.Manager{
		Id:       uuid.NewString(),
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	db.Exec("INSERT INTO managers (id, email, password) VALUES (?, ?, ?)", manager.Id, manager.Email, manager.Password)

	t.Run("success", func(t *testing.T) {
		found, err := repo.FindByEmail("test@example.com")
		if err != nil {
			t.Fatalf("FindByEmail() returned an unexpected error: %v", err)
		}
		if found.Email != manager.Email || found.Id != manager.Id {
			t.Errorf("FindByEmail() got = %+v, want %+v", found, manager)
		}
	})

	t.Run("fail_not_found", func(t *testing.T) {
		_, err := repo.FindByEmail("notfound@example.com")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("expected error %v, got %v", domain.ErrInvalidCredentials, err)
		}
	})
}

func TestSqliteRepository_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteRepository(db)
	db.Close()

	t.Run("FindById_db_error", func(t *testing.T) {
		_, err := repo.FindById("any-id")
		if !errors.Is(err, domain.ErrRepository) {
			t.Errorf("expected ErrRepository, got %v", err)
		}
	})

	t.Run("Save_db_error", func(t *testing.T) {
		err := repo.Save(&domain.Product{})
		if !errors.Is(err, domain.ErrRepository) {
			t.Errorf("expected ErrRepository, got %v", err)
		}
	})

	t.Run("FindByEmail_db_error", func(t *testing.T) {
		_, err := repo.FindByEmail("any@email.com")
		if !errors.Is(err, domain.ErrRepository) {
			t.Errorf("expected ErrRepository, got %v", err)
		}
	})
}
