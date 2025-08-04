package services

import (

    "secret-manager/internal/models"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type SecretService struct {
    db         *gorm.DB
    encryption *EncryptionService
}

func NewSecretService(db *gorm.DB, encryption *EncryptionService) *SecretService {
    return &SecretService{
        db:         db,
        encryption: encryption,
    }
}

func (s *SecretService) CreateSecret(secret *models.Secret) error {
    encryptedValue, err := s.encryption.Encrypt(secret.Value)
    if err != nil {
        return err
    }
    
    secret.Value = encryptedValue
    secret.Version = 1
    
    if err := s.db.Create(secret).Error; err != nil {
        return err
    }
    
    version := &models.SecretVersion{
        SecretID:  secret.ID,
        Value:     encryptedValue,
        Version:   1,
        CreatedBy: secret.CreatedBy,
    }
    
    return s.db.Create(version).Error
}

func (s *SecretService) GetSecret(id uuid.UUID, userID string) (*models.Secret, error) {
    var secret models.Secret
    if err := s.db.Preload("Tags").Preload("Project").First(&secret, "id = ?", id).Error; err != nil {
        return nil, err
    }
    
    decryptedValue, err := s.encryption.Decrypt(secret.Value)
    if err != nil {
        return nil, err
    }
    
    secret.Value = decryptedValue
    return &secret, nil
}

func (s *SecretService) ListSecrets(userID string, projectID *uuid.UUID, tags []string) ([]models.Secret, error) {
    var secrets []models.Secret
    query := s.db.Preload("Tags").Preload("Project")
    
    if projectID != nil {
        query = query.Where("project_id = ?", *projectID)
    }
    
    if len(tags) > 0 {
        query = query.Joins("JOIN secret_tags ON secrets.id = secret_tags.secret_id").
            Joins("JOIN tags ON secret_tags.tag_id = tags.id").
            Where("tags.name IN ?", tags).
            Group("secrets.id")
    }
    
    if err := query.Find(&secrets).Error; err != nil {
        return nil, err
    }
    
    for i := range secrets {
        secrets[i].Value = "***"
    }
    
    return secrets, nil
}

func (s *SecretService) UpdateSecret(id uuid.UUID, updates map[string]interface{}, userID string) error {
    var secret models.Secret
    if err := s.db.First(&secret, "id = ?", id).Error; err != nil {
        return err
    }
    
    if value, exists := updates["value"]; exists {
        encryptedValue, err := s.encryption.Encrypt(value.(string))
        if err != nil {
            return err
        }
        updates["value"] = encryptedValue
        updates["version"] = secret.Version + 1
        
        version := &models.SecretVersion{
            SecretID:  secret.ID,
            Value:     encryptedValue,
            Version:   secret.Version + 1,
            CreatedBy: userID,
        }
        s.db.Create(version)
    }
    
    return s.db.Model(&secret).Updates(updates).Error
}

func (s *SecretService) DeleteSecret(id uuid.UUID) error {
    return s.db.Delete(&models.Secret{}, "id = ?", id).Error
}

func (s *SecretService) GetSecretVersions(secretID uuid.UUID) ([]models.SecretVersion, error) {
    var versions []models.SecretVersion
    err := s.db.Where("secret_id = ?", secretID).Order("version DESC").Find(&versions).Error
    
    for i := range versions {
        versions[i].Value = "***"
    }
    
    return versions, err
}

func (s *SecretService) GetSecretVersion(secretID uuid.UUID, version int, userID string) (*models.SecretVersion, error) {
    var secretVersion models.SecretVersion
    if err := s.db.Where("secret_id = ? AND version = ?", secretID, version).First(&secretVersion).Error; err != nil {
        return nil, err
    }
    
    decryptedValue, err := s.encryption.Decrypt(secretVersion.Value)
    if err != nil {
        return nil, err
    }
    
    secretVersion.Value = decryptedValue
    return &secretVersion, nil
}
