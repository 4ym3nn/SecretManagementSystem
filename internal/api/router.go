package api

import (
    "secret-manager/internal/api/handlers"
    "secret-manager/internal/api/middleware"
    "secret-manager/internal/services"
    
    "github.com/gin-gonic/gin"
)

func SetupRouter(secretService *services.SecretService, auditService *services.AuditService, authService *services.AuthService) *gin.Engine {
    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()
    
    r.Use(middleware.CORSMiddleware())
    r.Use(middleware.AuditMiddleware(auditService))
    
    secretHandler := handlers.NewSecretHandler(secretService, auditService)
    authHandler := handlers.NewAuthHandler(authService)
    auditHandler := handlers.NewAuditHandler(auditService)
    
    api := r.Group("/api/v1")
    {
        api.POST("/auth/login", authHandler.Login)
        
        protected := api.Group("/")
        protected.Use(middleware.AuthMiddleware(authService))
        {
            secrets := protected.Group("/secrets")
            {
                secrets.POST("/", secretHandler.CreateSecret)
                secrets.GET("/", secretHandler.ListSecrets)
                secrets.GET("/:id", secretHandler.GetSecret)
                secrets.PUT("/:id", secretHandler.UpdateSecret)
                secrets.DELETE("/:id", secretHandler.DeleteSecret)
                secrets.GET("/:id/versions", secretHandler.GetSecretVersions)
                secrets.GET("/:id/versions/:version", secretHandler.GetSecretVersion)
            }
            
            audit := protected.Group("/audit")
            {
                audit.GET("/logs", auditHandler.GetLogs)
            }
        }
    }
    
    return r
}
