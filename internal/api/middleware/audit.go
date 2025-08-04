package middleware

import (
    "secret-manager/internal/services"
    
    "github.com/gin-gonic/gin"
)

func AuditMiddleware(auditService *services.AuditService) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if userID, exists := c.Get("user_id"); exists {
            auditService.Log(
                userID.(string),
                c.Request.Method,
                c.FullPath(),
                "",
                c.ClientIP(),
                c.GetHeader("User-Agent"),
                gin.H{
                    "status_code": c.Writer.Status(),
                    "path":        c.Request.URL.Path,
                },
            )
        }
    }
}
