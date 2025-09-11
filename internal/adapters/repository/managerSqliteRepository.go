package repository

import (
	"database/sql"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
)

type managerRepository struct {
	db *sql.DB
}

func NewManagerRepository(db *sql.DB) *managerRepository {
	return &managerRepository{
		db: db,
	}
}

func (repo *managerRepository) FindByEmail(email string) (*domain.Manager, error) {
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
