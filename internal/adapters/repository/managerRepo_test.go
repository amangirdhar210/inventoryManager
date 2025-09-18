package repository

import (
	"errors"
	"log"

	//"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/stretchr/testify/require"
)

func Test_FindByEmail_WhenManagerIdExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error in setting up sql mock")
	}

	managerRepo := NewManagerRepository(db)
	email := "admin@test.go"
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow("id", email, "hashedpwd123")
	mock.ExpectQuery("SELECT id, email, password FROM managers WHERE email").WithArgs(email).WillReturnRows(rows)
	manager, err := managerRepo.FindByEmail(email)
	require.NoError(t, err)
	require.NotNil(t, manager)
	require.Equal(t, "id", manager.Id)
	require.Equal(t, email, manager.Email)
	require.Equal(t, "hashedpwd123", manager.Password)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func Test_FindByEmail_WhenNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error in setting up sql mock")
	}

	managerRepo := NewManagerRepository(db)
	email := "admin@test.go"
	rows := sqlmock.NewRows([]string{"id", "email", "password"})
	mock.ExpectQuery("SELECT id, email, password FROM managers WHERE email").WithArgs(email).WillReturnRows(rows)

	manager, err := managerRepo.FindByEmail(email)

	require.Equal(t, domain.ErrManagerNotFound, err)
	require.Nil(t, manager)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func Test_FindByEmail_WhenUnexpectedDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error in setting up sql mock")
	}

	managerRepo := NewManagerRepository(db)
	email := "admin@test.go"
	rows := sqlmock.NewRows([]string{"id", "email", "password"})
	mock.ExpectQuery("SELECT id, email, password FROM managers WHERE email").WithArgs(email).WillReturnRows(rows).WillReturnError(errors.New("Oh NO!!!!"))

	manager, err := managerRepo.FindByEmail(email)

	require.Equal(t, domain.ErrRepository, err)
	require.Nil(t, manager)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
