package main

import (
	"fmt"
	"log"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/http"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/notifier"
	"github.com/amangirdhar210/inventory-manager/internal/adapters/repository"
	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/amangirdhar210/inventory-manager/utils/auth"
	"github.com/amangirdhar210/inventory-manager/utils/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlite.SetupDatabase("./inventory.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	managerRepo := repository.NewManagerRepository(db)
	productRepo := repository.NewProductRepository(db)
	logNotifier := notifier.NewLogNotifier()
	tokenGenerator := auth.NewJWTGenerator(config.JWTSecretKey)

	inventoryService := service.NewInventoryService(productRepo, logNotifier)
	authService := service.NewAuthService(managerRepo, tokenGenerator)

	HTTPServer := http.NewHTTPServer(inventoryService, authService)

	fmt.Println("Inventory Management Server starting on port 8080....")
	log.Fatal(HTTPServer.ListenAndServe())
}
