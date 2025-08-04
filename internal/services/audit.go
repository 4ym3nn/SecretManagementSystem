package services

import (
    "encoding/json"
    "secret-manager/internal/models"
    
    "gorm.io/gorm"
)

type AuditService struct {
    db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
    return &AuditService{db: db}
}

func (a *AuditService) Log(userID, action, resource, resourceID, ipAddress, userAgent string, details interface{}) error {
    detailsJSON, _ := json.Marshal(details)
    
    audit := &models.AuditLog{
        UserID:     userID,
        Action:     action,
        Resource:   resource,
        ResourceID: resourceID,
        Details:    string(detailsJSON),
        IPAddress:  ipAddress,
        UserAgent:  userAgent,
    }
    
    return a.db.Create(audit).Error
}

func (a *AuditService) GetLogs(userID string, resource string, limit, offset int) ([]models.AuditLog, error) {
    var logs []models.AuditLog
    query := a.db.Order("created_at DESC")
    
    if userID != "" {
        query = query.Where("user_id = ?", userID)
    }
    
    if resource != "" {
        query = query.Where("resource = ?", resource)
    }
    
    return logs, query.Limit(limit).Offset(offset).Find(&logs).Error
}
