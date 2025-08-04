package database

import (
    "secret-manager/internal/models"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

func Migrate(db *gorm.DB) error {
    db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
    
    return db.AutoMigrate(
        &models.User{},
        &models.Project{},
        &models.Tag{},
        &models.Secret{},
        &models.SecretVersion{},
        &models.AuditLog{},
    )
}
