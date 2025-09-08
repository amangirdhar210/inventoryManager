package repository

import (
	"database/sql"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *sqliteRepository {
	return &sqliteRepository{
		db: db,
	}
}

func (repo *sqliteRepository) FindById(id string) (*domain.Product, error) {
	row := repo.db.QueryRow("SELECT id, name, price, quantity FROM products where id=?", id)

	var product domain.Product
	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, domain.ErrRepository
	}
	return &product, nil
}

func (repo *sqliteRepository) Save(product *domain.Product) error {
	statement, err := repo.db.Prepare("INSERT INTO products(id, name, price,quantity) VALUES(?,?,?,?)")
	if err != nil {
		return domain.ErrRepository
	}
	defer statement.Close()

	_, err = statement.Exec(product.Id, product.Name, product.Price, product.Quantity)
	if err != nil {
		return domain.ErrRepository
	}
	return nil
}

func (repo *sqliteRepository) Update(product *domain.Product) error {
	statement, err := repo.db.Prepare("UPDATE products SET name=?, price=?, quantity=? WHERE id =?")
	if err != nil {
		return domain.ErrRepository
	}
	defer statement.Close()
	_, err = statement.Exec(product.Name, product.Price, product.Quantity, product.Id)
	if err != nil {
		return domain.ErrRepository
	}
	return nil
}

func (repo *sqliteRepository) DeleteById(id string) error {
	statement, err := repo.db.Prepare("DELETE FROM products WHERE id =?")
	if err != nil {
		return domain.ErrRepository
	}
	defer statement.Close()

	res, err := statement.Exec(id)
	if err != nil {
		return domain.ErrRepository
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (repo *sqliteRepository) ListAll() ([]domain.Product, error) {
	rows, err := repo.db.Query("SELECT id, name, price, quantity FROM products")
	if err != nil {
		return nil, domain.ErrRepository
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, domain.ErrRepository
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, domain.ErrRepository
	}
	return products, nil
}

func (repo *sqliteRepository) FindByEmail(email string) (*domain.Manager, error) {
	row := repo.db.QueryRow("SELECT id, email, password FROM managers WHERE email = ?", email)

	manager := &domain.Manager{}
	err := row.Scan(&manager.Id, &manager.Email, &manager.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, domain.ErrRepository
	}
	return manager, nil
}
