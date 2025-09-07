package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amangirdhar210/inventory-manager/internal/adapters/handler"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/notifier"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/repository"
	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := SetupDatabase("./inventory.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()
	inventoryRepo := repository.NewSQLiteRepository(db)
	logNotifier := notifier.NewLogNotifier()

	inventoryService := service.NewInventoryService(inventoryRepo, logNotifier)

	inventoryHandler := handler.NewHTTPHandler(inventoryService)

	router := mux.NewRouter()

	router.HandleFunc("/products", inventoryHandler.AddProduct).Methods("POST")
	router.HandleFunc("/products/{id}", inventoryHandler.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}/sell", inventoryHandler.SellProductUnits).Methods("POST")
	router.HandleFunc("/products/{id}/restock", inventoryHandler.RestockProduct).Methods("POST")
	router.HandleFunc("/products/{id}", inventoryHandler.UpdateProductPrice).Methods("PUT")
	router.HandleFunc("/products/{id}", inventoryHandler.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/products", inventoryHandler.GetAllProducts).Methods("GET")
	router.HandleFunc("/inventory/value", inventoryHandler.GetInventoryValue).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	fmt.Println("Inventory Management Server starting on port 8080....")
	log.Fatal(server.ListenAndServe())
}
