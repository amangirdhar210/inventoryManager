package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/handler"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/notifier"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/repository"
	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/amangirdhar210/inventory-manager/utils/auth"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := SetupDatabase("./inventory.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	sqliteRepo := repository.NewSQLiteRepository(db)
	logNotifier := notifier.NewLogNotifier()
	tokenGenerator := auth.NewJWTGenerator(config.JWTSecretKey)

	inventoryService := service.NewInventoryService(sqliteRepo, logNotifier)
	authService := service.NewAuthService(sqliteRepo, tokenGenerator)

	inventoryHandler := handler.NewHTTPHandler(inventoryService, authService)

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
	fmt.Println("Inventory Management Server starting on port 8080....")
	log.Fatal(server.ListenAndServe())
}
