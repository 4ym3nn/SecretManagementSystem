package models

import (
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Secret struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name        string    `json:"name" gorm:"not null;index"`
    Value       string    `json:"value" gorm:"not null"`
    Description string    `json:"description"`
    Tags        []Tag     `json:"tags" gorm:"many2many:secret_tags;"`
    ProjectID   *uuid.UUID `json:"project_id" gorm:"type:uuid;index"`
    Project     *Project  `json:"project,omitempty"`
    Version     int       `json:"version" gorm:"default:1"`
    CreatedBy   string    `json:"created_by" gorm:"not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type SecretVersion struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    SecretID  uuid.UUID `json:"secret_id" gorm:"type:uuid;not null;index"`
    Secret    Secret    `json:"secret"`
    Value     string    `json:"value" gorm:"not null"`
    Version   int       `json:"version" gorm:"not null"`
    CreatedBy string    `json:"created_by" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name      string    `json:"name" gorm:"unique;not null"`
    Color     string    `json:"color" gorm:"default:#007bff"`
    CreatedAt time.Time `json:"created_at"`
}

type Project struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name        string    `json:"name" gorm:"unique;not null"`
    Description string    `json:"description"`
    CreatedBy   string    `json:"created_by" gorm:"not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type AuditLog struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    UserID    string    `json:"user_id" gorm:"not null;index"`
    Action    string    `json:"action" gorm:"not null"`
    Resource  string    `json:"resource" gorm:"not null"`
    ResourceID string   `json:"resource_id" gorm:"index"`
    Details   string    `json:"details"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
    CreatedAt time.Time `json:"created_at"`
}

type User struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Role      string    `json:"role" gorm:"default:user"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (s *Secret) BeforeCreate(tx *gorm.DB) error {
    if s.ID == uuid.Nil {
        s.ID = uuid.New()
    }
    return nil
}
