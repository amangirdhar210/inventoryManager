package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/gorilla/mux"
)

type mockInventoryService struct {
	AddProductFunc         func(name string, price float64, quantity int) (*domain.Product, error)
	GetProductFunc         func(id string) (*domain.Product, error)
	SellProductUnitsFunc   func(id string, quantity int) (*domain.Product, error)
	RestockProductFunc     func(id string, quantity int) (*domain.Product, error)
	DeleteProductFunc      func(id string) error
	UpdateProductPriceFunc func(id string, newPrice float64) error
	GetAllProductsFunc     func() ([]domain.Product, error)
	GetInventoryValueFunc  func() (float64, error)
}

func (m *mockInventoryService) AddProduct(name string, price float64, quantity int) (*domain.Product, error) {
	return m.AddProductFunc(name, price, quantity)
}
func (m *mockInventoryService) GetProduct(id string) (*domain.Product, error) {
	return m.GetProductFunc(id)
}
func (m *mockInventoryService) SellProductUnits(id string, quantity int) (*domain.Product, error) {
	return m.SellProductUnitsFunc(id, quantity)
}
func (m *mockInventoryService) RestockProduct(id string, quantity int) (*domain.Product, error) {
	return m.RestockProductFunc(id, quantity)
}
func (m *mockInventoryService) DeleteProduct(id string) error {
	return m.DeleteProductFunc(id)
}
func (m *mockInventoryService) UpdateProductPrice(id string, newPrice float64) error {
	return m.UpdateProductPriceFunc(id, newPrice)
}
func (m *mockInventoryService) GetAllProducts() ([]domain.Product, error) {
	return m.GetAllProductsFunc()
}
func (m *mockInventoryService) GetInventoryValue() (float64, error) {
	return m.GetInventoryValueFunc()
}

func TestHTTPHandler_AddProduct(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		setupMock      func(*mockInventoryService)
		wantStatusCode int
		wantBody       string
	}{
		{
			"success",
			`{"name":"Test Laptop","price":1500.50,"quantity":10}`,
			func(m *mockInventoryService) {
				m.AddProductFunc = func(name string, price float64, quantity int) (*domain.Product, error) {
					return &domain.Product{Id: "new-id", Name: name, Price: price, Quantity: quantity}, nil
				}
			},
			http.StatusCreated,
			`"Id":"new-id"`,
		},
		{
			"fail_invalid_body",
			`{"name":"Test Laptop"`,
			func(m *mockInventoryService) {},
			http.StatusBadRequest,
			`Invalid Request Body`,
		},
		{
			"fail_service_error",
			`{"name":"Test Laptop","price":1500.50,"quantity":10}`,
			func(m *mockInventoryService) {
				m.AddProductFunc = func(name string, price float64, quantity int) (*domain.Product, error) {
					return nil, errors.New("service failure")
				}
			},
			http.StatusInternalServerError,
			`service failure`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockInventoryService{}
			tt.setupMock(mockService)
			handler := NewHTTPHandler(mockService)

			req := httptest.NewRequest("POST", "/products", strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()

			handler.AddProduct(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestHTTPHandler_GetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		setupMock      func(*mockInventoryService)
		wantStatusCode int
		wantBody       string
	}{
		{
			"success",
			"prod-123",
			func(m *mockInventoryService) {
				m.GetProductFunc = func(id string) (*domain.Product, error) {
					return &domain.Product{Id: "prod-123", Name: "Found Product"}, nil
				}
			},
			http.StatusOK,
			`"Id":"prod-123"`,
		},
		{
			"fail_not_found",
			"prod-456",
			func(m *mockInventoryService) {
				m.GetProductFunc = func(id string) (*domain.Product, error) {
					return nil, errors.New("not found")
				}
			},
			http.StatusNotFound,
			`not found`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockInventoryService{}
			tt.setupMock(mockService)
			handler := NewHTTPHandler(mockService)

			req := httptest.NewRequest("GET", "/products/"+tt.productID, nil)
			rr := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})

			handler.GetProduct(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestHTTPHandler_SellProductUnits(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		reqBody        string
		setupMock      func(*mockInventoryService)
		wantStatusCode int
		wantBody       string
	}{
		{
			"success",
			"prod-123",
			`{"quantity": 5}`,
			func(m *mockInventoryService) {
				m.SellProductUnitsFunc = func(id string, quantity int) (*domain.Product, error) {
					return &domain.Product{}, nil
				}
			},
			http.StatusOK,
			"",
		},
		{
			"fail_insufficient_stock",
			"prod-123",
			`{"quantity": 50}`,
			func(m *mockInventoryService) {
				m.SellProductUnitsFunc = func(id string, quantity int) (*domain.Product, error) {
					return nil, errors.New("insufficient stock")
				}
			},
			http.StatusBadRequest,
			"insufficient stock",
		},
		{
			"fail_product_not_found",
			"prod-456",
			`{"quantity": 5}`,
			func(m *mockInventoryService) {
				m.SellProductUnitsFunc = func(id string, quantity int) (*domain.Product, error) {
					return nil, errors.New("product not found")
				}
			},
			http.StatusNotFound,
			"product not found",
		},
		{
			"fail_invalid_body",
			"prod-123",
			`{"quantity":}`,
			func(m *mockInventoryService) {},
			http.StatusBadRequest,
			"Invalid request body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockInventoryService{}
			tt.setupMock(mockService)
			handler := NewHTTPHandler(mockService)

			req := httptest.NewRequest("POST", "/products/"+tt.productID+"/sell", strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})

			handler.SellProductUnits(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if tt.wantBody != "" && !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestHTTPHandler_UpdateProductPrice(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		reqBody        string
		setupMock      func(*mockInventoryService)
		wantStatusCode int
		wantBody       string
	}{
		{
			"success",
			"prod-123",
			`{"price": 99.99}`,
			func(m *mockInventoryService) {
				m.UpdateProductPriceFunc = func(id string, newPrice float64) error {
					return nil
				}
			},
			http.StatusOK,
			"",
		},
		{
			"fail_product_not_found",
			"prod-456",
			`{"price": 99.99}`,
			func(m *mockInventoryService) {
				m.UpdateProductPriceFunc = func(id string, newPrice float64) error {
					return domain.ErrProductNotFound
				}
			},
			http.StatusNotFound,
			"product not found",
		},
		{
			"fail_invalid_price",
			"prod-123",
			`{"price": -10}`,
			func(m *mockInventoryService) {
				m.UpdateProductPriceFunc = func(id string, newPrice float64) error {
					return errors.New("price must be greater than zero")
				}
			},
			http.StatusBadRequest,
			"price must be greater than zero",
		},
		{
			"fail_internal_server_error",
			"prod-123",
			`{"price": 99.99}`,
			func(m *mockInventoryService) {
				m.UpdateProductPriceFunc = func(id string, newPrice float64) error {
					return errors.New("database update failed")
				}
			},
			http.StatusInternalServerError,
			"internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockInventoryService{}
			tt.setupMock(mockService)
			handler := NewHTTPHandler(mockService)

			req := httptest.NewRequest("PUT", "/products/"+tt.productID+"/price", strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})

			handler.UpdateProductPrice(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if tt.wantBody != "" && !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestHTTPHandler_GetAllProducts(t *testing.T) {
	t.Run("success_with_products", func(t *testing.T) {
		mockService := &mockInventoryService{
			GetAllProductsFunc: func() ([]domain.Product, error) {
				return []domain.Product{
					{Id: "p1", Name: "Product 1"},
					{Id: "p2", Name: "Product 2"},
				}, nil
			},
		}
		handler := NewHTTPHandler(mockService)
		req := httptest.NewRequest("GET", "/products", nil)
		rr := httptest.NewRecorder()

		handler.GetAllProducts(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
		var products []domain.Product
		json.NewDecoder(rr.Body).Decode(&products)
		if len(products) != 2 {
			t.Errorf("expected 2 products, got %d", len(products))
		}
	})

	t.Run("fail_service_error", func(t *testing.T) {
		mockService := &mockInventoryService{
			GetAllProductsFunc: func() ([]domain.Product, error) {
				return nil, errors.New("db error")
			},
		}
		handler := NewHTTPHandler(mockService)
		req := httptest.NewRequest("GET", "/products", nil)
		rr := httptest.NewRecorder()

		handler.GetAllProducts(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusInternalServerError)
		}
	})
}

func TestHTTPHandler_GetInventoryValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockInventoryService{
			GetInventoryValueFunc: func() (float64, error) {
				return 1234.56, nil
			},
		}
		handler := NewHTTPHandler(mockService)
		req := httptest.NewRequest("GET", "/inventory/value", nil)
		rr := httptest.NewRecorder()

		handler.GetInventoryValue(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
		if !strings.Contains(rr.Body.String(), "1234.56") {
			t.Errorf("body does not contain expected value, got %s", rr.Body.String())
		}
	})
}

func TestHTTPHandler_DeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockInventoryService{
			DeleteProductFunc: func(id string) error {
				return nil
			},
		}
		handler := NewHTTPHandler(mockService)
		req := httptest.NewRequest("DELETE", "/products/prod-123", nil)
		rr := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": "prod-123"})

		handler.DeleteProduct(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
	})
}

func TestHTTPHandler_RestockProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockInventoryService{
			RestockProductFunc: func(id string, quantity int) (*domain.Product, error) {
				return &domain.Product{}, nil
			},
		}
		handler := NewHTTPHandler(mockService)
		reqBody := `{"quantity": 100}`
		req := httptest.NewRequest("POST", "/products/prod-123/restock", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": "prod-123"})

		handler.RestockProduct(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
	})
}
