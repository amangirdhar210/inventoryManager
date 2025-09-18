package repository

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/stretchr/testify/require"
)

func Test_FindById_WhenValidId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error in seting up mock db ")
	}

	productRepo := NewProductRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	rows.AddRow("valid-id", "test-product", 50.50, 20)

	mock.ExpectQuery("SELECT id, name, price, quantity FROM products where id=?").WithArgs("valid-id").WillReturnRows(rows)

	product, err := productRepo.FindById("valid-id")
	require.NoError(t, err)
	require.Equal(t, "valid-id", product.Id)
	require.Equal(t, "test-product", product.Name)
	require.Equal(t, 50.50, product.Price)
	require.Equal(t, 20, product.Quantity)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_FindById_WhenIdNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error in seting up mock db ")
	}

	productRepo := NewProductRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})

	mock.ExpectQuery("SELECT id, name, price, quantity FROM products where id=?").WithArgs("not-existing-id").WillReturnRows(rows)

	product, err := productRepo.FindById("not-existing-id")
	require.Nil(t, product)
	require.Error(t, err)
	require.Equal(t, domain.ErrProductNotFound, err)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_FindById_WhenUnexpectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error in seting up mock db ")
	}

	productRepo := NewProductRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})

	mock.ExpectQuery("SELECT id, name, price, quantity FROM products where id=?").WithArgs("test-id").
		WillReturnError(errors.New("Unexpected DB error")).WillReturnRows(rows)

	product, err := productRepo.FindById("test-id")

	require.Nil(t, product)
	require.Error(t, err)
	require.Equal(t, domain.ErrRepository, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func Test_SaveSuccessCase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("The mock db setup failed")
	}

	productRepo := NewProductRepository(db)
	sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	mockProduct := &domain.Product{
		Id:       "test-id",
		Name:     "test-product",
		Price:    10.10,
		Quantity: 5,
	}
	mock.ExpectPrepare("INSERT INTO products\\(id, name, price,quantity\\) VALUES\\(\\?,\\?,\\?,\\?\\)").
		ExpectExec().
		WithArgs(mockProduct.Id, mockProduct.Name, mockProduct.Price, mockProduct.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = productRepo.Save(mockProduct)

	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func Test_Save_ErrorCaseExecuting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("The mock db setup failed")
	}

	productRepo := NewProductRepository(db)
	sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	mockProduct := &domain.Product{
		Id:       "test-id",
		Name:     "test-product",
		Price:    10.10,
		Quantity: 5,
	}
	mock.ExpectPrepare("INSERT INTO products\\(id, name, price,quantity\\) VALUES\\(\\?,\\?,\\?,\\?\\)").
		ExpectExec().
		WithArgs(mockProduct.Id, mockProduct.Name, mockProduct.Price, mockProduct.Quantity).
		WillReturnError(errors.New("Some error occured"))

	err = productRepo.Save(mockProduct)

	require.Error(t, err)
	require.Equal(t, domain.ErrRepository, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func Test_Save_ErrorPreparingStatement(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("The mock db setup failed")
	}

	productRepo := NewProductRepository(db)
	sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	mockProduct := &domain.Product{
		Id:       "test-id",
		Name:     "test-product",
		Price:    10.10,
		Quantity: 5,
	}
	mock.ExpectPrepare("INSERT INTO products\\(id, name, price,quantity\\) VALUES\\(\\?,\\?,\\?,\\?\\)").
		WillReturnError(errors.New("Some error occured"))

	err = productRepo.Save(mockProduct)

	require.Error(t, err)
	require.Equal(t, domain.ErrRepository, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func Test_Update_ErrorPreparingStatement(t *testing.T) {
	db, mock, err := sqlmock.New()

	productRepo := NewProductRepository(db)
	if err != nil {
		log.Fatal("error setting up repo", err)
	}
	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE products SET name=?, price=?, quantity=? WHERE id =?")).
		WillReturnError(errors.New("error preparing statement"))
	UpdatedProduct := domain.Product{
		Id:       "existing-id",
		Name:     "product-name",
		Price:    200,
		Quantity: 5,
	}
	err = productRepo.Update(&UpdatedProduct)
	require.Error(t, err)
	require.Equal(t, domain.ErrRepository, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_Update_ErrorExecutingStatement(t *testing.T) {
	db, mock, err := sqlmock.New()

	productRepo := NewProductRepository(db)
	if err != nil {
		log.Fatal("error setting up repo", err)

	}
	UpdatedProduct := domain.Product{
		Id:       "existing-id",
		Name:     "product-name",
		Price:    200,
		Quantity: 5,
	}
	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE products SET name=?, price=?, quantity=? WHERE id =?")).ExpectExec().
		WithArgs(UpdatedProduct.Name, UpdatedProduct.Price, UpdatedProduct.Quantity, UpdatedProduct.Id)

	err = productRepo.Update(&UpdatedProduct)
	require.Error(t, err)
	require.Equal(t, domain.ErrRepository, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_DeleteById_ErrorPreparing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	productRepo := NewProductRepository(db)
	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	rows.AddRow("delete-id", "product-to-be-deleted", 100.02, 10)

	mock.ExpectPrepare("DELETE FROM products WHERE id =?").
		WillReturnError(errors.New("Error occured preparing"))

	err = productRepo.DeleteById("delete-id")
	require.Error(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_DeleteById_ErrorExecuting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	productRepo := NewProductRepository(db)
	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	rows.AddRow("delete-id", "product-to-be-deleted", 100.02, 10)

	mock.ExpectPrepare("DELETE FROM products WHERE id =?").
		ExpectExec().WithArgs("delete-id").WillReturnError(errors.New("error executing the querry"))

	err = productRepo.DeleteById("delete-id")
	require.Error(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_ListAll_SuccessCase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	productRepo := NewProductRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})
	rows.AddRow("TestId1", "testProduct1", 10.10, 10).
		AddRow("TestId2", "testProduct2", 20.10, 20).
		AddRow("TestId3", "testProduct3", 30.10, 30).
		AddRow("TestId4", "testProduct4", 40.10, 40).
		AddRow("TestId5", "testProduct5", 50.10, 50)

	mock.ExpectQuery("SELECT id, name, price, quantity FROM products").WillReturnRows(rows)
	products, err := productRepo.ListAll()

	require.NoError(t, err)
	require.Equal(t, 5, len(products))

}

func Test_ListAllWhenNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("The sql mock setup failed")
	}
	productRepo := NewProductRepository(db)
	rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity"})

	mock.ExpectQuery("SELECT id, name, price, quantity FROM products").WillReturnRows(rows)

	products, err := productRepo.ListAll()

	require.Equal(t, 0, len(products))
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func Test_ListAllWhenError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	productRepo := NewProductRepository(db)
	mock.ExpectQuery("SELECT id, name, price, quantity FROM products").WillReturnError(errors.New("error from db"))

	products, err := productRepo.ListAll()

	require.Nil(t, products)
	require.Error(t, err)
	require.Equal(t, domain.ErrRepository, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
