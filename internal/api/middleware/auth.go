package middleware

import (
    "net/http"
    "secret-manager/internal/services"
    "strings"
    
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        apiKey := c.GetHeader("X-API-Key")
        
        if apiKey != "" {
            if authService.ValidateAPIKey(apiKey) {
                c.Set("user_id", "api-user")
                c.Set("username", "api-user")
                c.Set("role", "api")
                c.Next()
                return
            }
        }
        
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := authService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}
