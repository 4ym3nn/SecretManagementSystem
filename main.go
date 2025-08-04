
// main.go
package main

import (
    "log"
    "secret-manager/internal/api"
    "secret-manager/internal/config"
    "secret-manager/internal/database"
    "secret-manager/internal/services"
)

func main() {
    cfg := config.Load()
    
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
    
    encryptionService := services.NewEncryptionService(cfg.EncryptionKey)
    secretService := services.NewSecretService(db, encryptionService)
    auditService := services.NewAuditService(db)
    authService := services.NewAuthService(cfg.JWTSecret, cfg.APIKeys)
    
    router := api.SetupRouter(secretService, auditService, authService)
    
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(router.Run(":" + cfg.Port))
}



