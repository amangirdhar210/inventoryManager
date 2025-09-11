package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type mockInventoryService struct {
	AddProductFunc         func(name string, price float64, quantity int) (*domain.Product, error)
	GetProductFunc         func(id string) (*domain.Product, error)
	SellProductUnitsFunc   func(id string, quantity int) error
	RestockProductFunc     func(id string, quantity int) error
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
func (m *mockInventoryService) SellProductUnits(id string, quantity int) error {
	return m.SellProductUnitsFunc(id, quantity)
}
func (m *mockInventoryService) RestockProduct(id string, quantity int) error {
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

type mockAuthService struct {
	LoginFunc func(email, password string) (string, error)
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	return m.LoginFunc(email, password)
}

func getTestToken() string {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(config.JWTSecretKey))
	return signedToken
}

func newTestRouter(handler *HTTPHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/logout", handler.Logout).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(handler.AuthMiddleware)
	apiRouter.HandleFunc("/products", handler.AddProduct).Methods("POST")
	apiRouter.HandleFunc("/products", handler.GetAllProducts).Methods("GET")
	apiRouter.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
	apiRouter.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")
	apiRouter.HandleFunc("/products/{id}/sell", handler.SellProductUnits).Methods("POST")
	apiRouter.HandleFunc("/products/{id}/restock", handler.RestockProduct).Methods("POST")
	apiRouter.HandleFunc("/products/{id}/price", handler.UpdateProductPrice).Methods("PUT")
	apiRouter.HandleFunc("/inventory/value", handler.GetInventoryValue).Methods("GET")

	return router
}

func TestHTTPHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		setupMock      func(*mockAuthService)
		wantStatusCode int
		wantBody       string
	}{
		{
			name:    "success",
			reqBody: `{"email":"test@example.com","password":"password123"}`,
			setupMock: func(m *mockAuthService) {
				m.LoginFunc = func(email, password string) (string, error) {
					return "fake-jwt-token", nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `"token":"fake-jwt-token"`,
		},
		{
			name:    "fail_invalid_credentials",
			reqBody: `{"email":"wrong@example.com","password":"wrong"}`,
			setupMock: func(m *mockAuthService) {
				m.LoginFunc = func(email, password string) (string, error) {
					return "", domain.ErrInvalidCredentials
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       domain.ErrInvalidCredentials.Error(),
		},
		{
			name:           "fail_invalid_body",
			reqBody:        `{"email":"bad"`,
			setupMock:      func(m *mockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Invalid request body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := &mockAuthService{}
			tt.setupMock(mockAuth)
			handler := NewHTTPHandler(nil, mockAuth)
			router := newTestRouter(handler)

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestHTTPHandler_ProductEndpoints_Auth(t *testing.T) {
	mockInventory := &mockInventoryService{
		GetProductFunc: func(id string) (*domain.Product, error) {
			return &domain.Product{Id: "prod-123"}, nil
		},
	}
	handler := NewHTTPHandler(mockInventory, nil)
	router := newTestRouter(handler)

	t.Run("success_with_valid_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/products/prod-123", nil)
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("fail_without_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/products/prod-123", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("fail_with_invalid_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/products/prod-123", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})
}

func TestHTTPHandler_GetProduct(t *testing.T) {
	t.Run("fail_not_found", func(t *testing.T) {
		mockInventory := &mockInventoryService{
			GetProductFunc: func(id string) (*domain.Product, error) {
				return nil, domain.ErrProductNotFound
			},
		}
		handler := NewHTTPHandler(mockInventory, nil)
		router := newTestRouter(handler)

		req := httptest.NewRequest("GET", "/api/products/prod-456", nil)
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusNotFound)
		}
		if !strings.Contains(rr.Body.String(), domain.ErrProductNotFound.Error()) {
			t.Errorf("body does not contain %q", domain.ErrProductNotFound.Error())
		}
	})
}

func TestHTTPHandler_SellProductUnits(t *testing.T) {
	t.Run("fail_insufficient_stock", func(t *testing.T) {
		mockInventory := &mockInventoryService{
			SellProductUnitsFunc: func(id string, quantity int) error {
				return domain.ErrInsufficientStock
			},
		}
		handler := NewHTTPHandler(mockInventory, nil)
		router := newTestRouter(handler)

		reqBody := `{"quantity": 50}`
		req := httptest.NewRequest("POST", "/api/products/prod-123/sell", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
		}
		if !strings.Contains(rr.Body.String(), domain.ErrInsufficientStock.Error()) {
			t.Errorf("body does not contain %q", domain.ErrInsufficientStock.Error())
		}
	})
}

func TestHTTPHandler_AddProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockInventory := &mockInventoryService{
			AddProductFunc: func(name string, price float64, quantity int) (*domain.Product, error) {
				return &domain.Product{Id: "new-id", Name: name, Price: price, Quantity: quantity}, nil
			},
		}
		handler := NewHTTPHandler(mockInventory, nil)
		router := newTestRouter(handler)

		reqBody := `{"name":"Test Laptop","price":1500.50,"quantity":10}`
		req := httptest.NewRequest("POST", "/api/products", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusCreated)
		}
		var product domain.Product
		json.NewDecoder(rr.Body).Decode(&product)
		if product.Id != "new-id" {
			t.Errorf("expected product id to be 'new-id', got %s", product.Id)
		}
	})
}

func TestHTTPHandler_GetAllProducts(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mockInventoryService)
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "success_with_products",
			setupMock: func(m *mockInventoryService) {
				m.GetAllProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{
						{Id: "p1", Name: "Product 1"},
						{Id: "p2", Name: "Product 2"},
					}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `"Id":"p2"`,
		},
		{
			name: "success_no_products",
			setupMock: func(m *mockInventoryService) {
				m.GetAllProductsFunc = func() ([]domain.Product, error) {
					return []domain.Product{}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `[]`,
		},
		{
			name: "fail_service_error",
			setupMock: func(m *mockInventoryService) {
				m.GetAllProductsFunc = func() ([]domain.Product, error) {
					return nil, domain.ErrRepository
				}
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockInventoryService{}
			tt.setupMock(mockService)
			handler := NewHTTPHandler(mockService, nil)
			router := newTestRouter(handler)

			req := httptest.NewRequest("GET", "/api/products", nil)
			req.Header.Set("Authorization", "Bearer "+getTestToken())
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestHTTPHandler_DeleteProduct(t *testing.T) {
	mockService := &mockInventoryService{
		DeleteProductFunc: func(id string) error {
			if id == "prod-123" {
				return nil
			}
			return domain.ErrProductNotFound
		},
	}
	handler := NewHTTPHandler(mockService, nil)
	router := newTestRouter(handler)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/products/prod-123", nil)
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
		if !strings.Contains(rr.Body.String(), "deleted successfully") {
			t.Errorf("body does not contain success message")
		}
	})

	t.Run("fail_not_found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/products/prod-456", nil)
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestHTTPHandler_GetInventoryValue(t *testing.T) {
	mockService := &mockInventoryService{
		GetInventoryValueFunc: func() (float64, error) {
			return 1234.56, nil
		},
	}
	handler := NewHTTPHandler(mockService, nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest("GET", "/api/inventory/value", nil)
	req.Header.Set("Authorization", "Bearer "+getTestToken())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "1234.56") {
		t.Errorf("body does not contain expected value, got %s", rr.Body.String())
	}
}

func TestHTTPHandler_RestockProduct(t *testing.T) {
	mockService := &mockInventoryService{
		RestockProductFunc: func(id string, quantity int) error {
			return nil
		},
	}
	handler := NewHTTPHandler(mockService, nil)
	router := newTestRouter(handler)

	t.Run("success", func(t *testing.T) {
		reqBody := `{"quantity": 50}`
		req := httptest.NewRequest("POST", "/api/products/prod-123/restock", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("fail_invalid_body", func(t *testing.T) {
		reqBody := `{"quantity":}`
		req := httptest.NewRequest("POST", "/api/products/prod-123/restock", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestHTTPHandler_UpdateProductPrice(t *testing.T) {
	mockService := &mockInventoryService{
		UpdateProductPriceFunc: func(id string, newPrice float64) error {
			if id == "prod-456" {
				return domain.ErrProductNotFound
			}
			if newPrice <= 0 {
				return domain.ErrProductInvalid
			}
			return nil
		},
	}
	handler := NewHTTPHandler(mockService, nil)
	router := newTestRouter(handler)

	t.Run("success", func(t *testing.T) {
		reqBody := `{"price": 99.99}`
		req := httptest.NewRequest("PUT", "/api/products/prod-123/price", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("fail_not_found", func(t *testing.T) {
		reqBody := `{"price": 99.99}`
		req := httptest.NewRequest("PUT", "/api/products/prod-456/price", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+getTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestHTTPHandler_Logout(t *testing.T) {
	handler := NewHTTPHandler(nil, nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest("POST", "/logout", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "logout successful") {
		t.Errorf("body does not contain logout message")
	}
}
