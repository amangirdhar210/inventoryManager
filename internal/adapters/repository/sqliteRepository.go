package repository

import (
	"database/sql"
	"fmt"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/amangirdhar210/inventory-manager/internal/core/ports"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) ports.ProductRepository {
	return &sqliteRepository{
		db: db,
	}
}
func (productRepo *sqliteRepository) FindById(id string) (*domain.Product, error) {
	row := productRepo.db.QueryRow("SELECT id, name, price, quantity FROM products where id=?", id)

	var product domain.Product
	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %s not found", id)
		}
		return nil, err
	}
	return &product, nil
}

func (productRepo *sqliteRepository) Save(product *domain.Product) error {
	statement, err := productRepo.db.Prepare("INSERT INTO products(id, name, price,quantity) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(product.Id, product.Name, product.Price, product.Quantity)
	return err

}

func (productRepo *sqliteRepository) Update(product *domain.Product) error {
	statement, err := productRepo.db.Prepare("UPDATE products SET name=?, price=?, quantity=? WHERE id =?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(product.Name, product.Price, product.Quantity, product.Id)
	return err
}

func (productRepo *sqliteRepository) DeleteById(id string) error {
	statement, err := productRepo.db.Prepare("DELETE FROM products WHERE id =?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(id)
	return err
}

func (productRepo *sqliteRepository) ListAll() ([]domain.Product, error) {
	rows, err := productRepo.db.Query("SELECT id, name, price, quantity FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
