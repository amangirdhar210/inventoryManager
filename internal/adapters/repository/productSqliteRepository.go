package repository

import (
	"database/sql"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (repo *productRepository) FindById(id string) (*domain.Product, error) {
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

func (repo *productRepository) Save(product *domain.Product) error {
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

func (repo *productRepository) Update(product *domain.Product) error {
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

func (repo *productRepository) DeleteById(id string) error {
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

func (repo *productRepository) ListAll() ([]domain.Product, error) {
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
