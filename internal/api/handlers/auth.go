package handlers

import (
    "net/http"
    "secret-manager/internal/services"
    
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Simple demo authentication - in production, verify against database
    if req.Username == "admin" && req.Password == "password" {
        token, err := h.authService.GenerateToken("1", req.Username, "admin")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "token": token,
            "user": gin.H{
                "id":       "1",
                "username": req.Username,
                "role":     "admin",
            },
        })
        return
    }
    
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

