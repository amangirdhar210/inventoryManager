###INVENTORY MANAGEMENT SYSTEM

Features:
1. Add Product with stock quantity and price 
2. Update stock when items are sold or restocked
3. Generate low-stock alerts
4. Print total inventory value

Commands to run InventoryManager Server:
git clone git@github.com:amangirdhar210/inventoryManager.git

go mod tidy

go run ./cmd/server/

Command to run Inventory Client

go run tools/inventoryClient/cmd/client/main.go

