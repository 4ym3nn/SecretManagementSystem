package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "secret-manager/internal/api"
    "secret-manager/internal/config"
    "secret-manager/internal/database"
    "secret-manager/internal/models"
    "secret-manager/internal/services"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
    gin.SetMode(gin.TestMode)
    
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    database.Migrate(db)
    
    cfg := &config.Config{
        JWTSecret:     "test-secret",
        EncryptionKey: "test-encryption-key-32-bytes-long",
        APIKeys:       []string{"test-api-key"},
    }
    
    encryptionService := services.NewEncryptionService(cfg.EncryptionKey)
    secretService := services.NewSecretService(db, encryptionService)
    auditService := services.NewAuditService(db)
    authService := services.NewAuthService(cfg.JWTSecret, cfg.APIKeys)
    
    return api.SetupRouter(secretService, auditService, authService), db
}

func TestCreateSecret(t *testing.T) {
    router, _ := setupTestRouter()
    
    secret := models.Secret{
        Name:        "test-secret",
        Value:       "secret-value",
        Description: "Test secret",
    }
    
    jsonValue, _ := json.Marshal(secret)
    req, _ := http.NewRequest("POST", "/api/v1/secrets", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", "test-api-key")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response models.Secret
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "test-secret", response.Name)
    assert.Equal(t, "***", response.Value)
}

func TestGetSecret(t *testing.T) {
    router, db := setupTestRouter()
    
    encService := services.NewEncryptionService("test-encryption-key-32-bytes-long")
    encryptedValue, _ := encService.Encrypt("secret-value")
    
    secret := models.Secret{
        ID:          uuid.New(),
        Name:        "test-secret",
        Value:       encryptedValue,
        Description: "Test secret",
        CreatedBy:   "test-user",
    }
    
    db.Create(&secret)
    
    req, _ := http.NewRequest("GET", "/api/v1/secrets/"+secret.ID.String(), nil)
    req.Header.Set("X-API-Key", "test-api-key")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response models.Secret
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "test-secret", response.Name)
    assert.Equal(t, "secret-value", response.Value)
}

func TestAuthentication(t *testing.T) {
    router, _ := setupTestRouter()
    
    req, _ := http.NewRequest("GET", "/api/v1/secrets", nil)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin(t *testing.T) {
    router, _ := setupTestRouter()
    
    loginReq := map[string]string{
        "username": "admin",
        "password": "password",
    }
    
    jsonValue, _ := json.Marshal(loginReq)
    req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Contains(t, response, "token")
    assert.Contains(t, response, "user")
}
