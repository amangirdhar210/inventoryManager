package http

import (
	"net/http"
	"time"

	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/gorilla/mux"
)

func NewHTTPServer(inventoryService service.InventoryService, authService service.AuthService) *http.Server {
	inventoryHandler := NewHTTPHandler(inventoryService, authService)

	router := mux.NewRouter()

	router.HandleFunc("/login", inventoryHandler.Login).Methods("POST")
	router.HandleFunc("/logout", inventoryHandler.Logout).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(inventoryHandler.AuthMiddleware)

	apiRouter.HandleFunc("/products", inventoryHandler.AddProduct).Methods("POST")
	apiRouter.HandleFunc("/products/{id}", inventoryHandler.GetProduct).Methods("GET")
	apiRouter.HandleFunc("/products/{id}/sell", inventoryHandler.SellProductUnits).Methods("POST")
	apiRouter.HandleFunc("/products/{id}/restock", inventoryHandler.RestockProduct).Methods("POST")
	apiRouter.HandleFunc("/products/{id}/price", inventoryHandler.UpdateProductPrice).Methods("PUT")
	apiRouter.HandleFunc("/products/{id}", inventoryHandler.DeleteProduct).Methods("DELETE")
	apiRouter.HandleFunc("/products", inventoryHandler.GetAllProducts).Methods("GET")
	apiRouter.HandleFunc("/inventory/value", inventoryHandler.GetInventoryValue).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return server

}
