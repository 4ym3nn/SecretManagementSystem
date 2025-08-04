package handlers

import (
    "net/http"
    "secret-manager/internal/models"
    "secret-manager/internal/services"
    "strconv"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type SecretHandler struct {
    secretService *services.SecretService
    auditService  *services.AuditService
}

func NewSecretHandler(secretService *services.SecretService, auditService *services.AuditService) *SecretHandler {
    return &SecretHandler{
        secretService: secretService,
        auditService:  auditService,
    }
}

func (h *SecretHandler) CreateSecret(c *gin.Context) {
    var secret models.Secret
    if err := c.ShouldBindJSON(&secret); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetString("user_id")
    secret.CreatedBy = userID
    
    if err := h.secretService.CreateSecret(&secret); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    h.auditService.Log(userID, "CREATE", "secret", secret.ID.String(), c.ClientIP(), c.GetHeader("User-Agent"), gin.H{"name": secret.Name})
    
    secret.Value = "***"
    c.JSON(http.StatusCreated, secret)
}

func (h *SecretHandler) GetSecret(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    
    userID := c.GetString("user_id")
    secret, err := h.secretService.GetSecret(id, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
        return
    }
    
    h.auditService.Log(userID, "READ", "secret", id.String(), c.ClientIP(), c.GetHeader("User-Agent"), gin.H{"name": secret.Name})
    
    c.JSON(http.StatusOK, secret)
}

func (h *SecretHandler) ListSecrets(c *gin.Context) {
    userID := c.GetString("user_id")
    
    var projectID *uuid.UUID
    if projectIDStr := c.Query("project_id"); projectIDStr != "" {
        if id, err := uuid.Parse(projectIDStr); err == nil {
            projectID = &id
        }
    }
    
    var tags []string
    if tagsStr := c.Query("tags"); tagsStr != "" {
        tags = strings.Split(tagsStr, ",")
    }
    
    secrets, err := h.secretService.ListSecrets(userID, projectID, tags)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    h.auditService.Log(userID, "LIST", "secret", "", c.ClientIP(), c.GetHeader("User-Agent"), gin.H{"count": len(secrets)})
    
    c.JSON(http.StatusOK, secrets)
}

func (h *SecretHandler) UpdateSecret(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetString("user_id")
    if err := h.secretService.UpdateSecret(id, updates, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    h.auditService.Log(userID, "UPDATE", "secret", id.String(), c.ClientIP(), c.GetHeader("User-Agent"), updates)
    
    c.JSON(http.StatusOK, gin.H{"message": "Secret updated successfully"})
}

func (h *SecretHandler) DeleteSecret(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    
    userID := c.GetString("user_id")
    if err := h.secretService.DeleteSecret(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    h.auditService.Log(userID, "DELETE", "secret", id.String(), c.ClientIP(), c.GetHeader("User-Agent"), nil)
    
    c.JSON(http.StatusOK, gin.H{"message": "Secret deleted successfully"})
}

func (h *SecretHandler) GetSecretVersions(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    
    userID := c.GetString("user_id")
    versions, err := h.secretService.GetSecretVersions(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    h.auditService.Log(userID, "LIST_VERSIONS", "secret", id.String(), c.ClientIP(), c.GetHeader("User-Agent"), gin.H{"count": len(versions)})
    
    c.JSON(http.StatusOK, versions)
}

func (h *SecretHandler) GetSecretVersion(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    
    version, err := strconv.Atoi(c.Param("version"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version number"})
        return
    }
    
    userID := c.GetString("user_id")
    secretVersion, err := h.secretService.GetSecretVersion(id, version, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
        return
    }
    
    h.auditService.Log(userID, "READ_VERSION", "secret", id.String(), c.ClientIP(), c.GetHeader("User-Agent"), gin.H{"version": version})
    
    c.JSON(http.StatusOK, secretVersion)
}
