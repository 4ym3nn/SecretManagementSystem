package handlers

import (
    "net/http"
    "secret-manager/internal/services"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type AuditHandler struct {
    auditService *services.AuditService
}

func NewAuditHandler(auditService *services.AuditService) *AuditHandler {
    return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) GetLogs(c *gin.Context) {
    userID := c.Query("user_id")
    resource := c.Query("resource")
    
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    logs, err := h.auditService.GetLogs(userID, resource, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "logs":  logs,
        "limit": limit,
        "offset": offset,
    })
}
