package services

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
    jwtSecret string
    apiKeys   []string
}

type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func NewAuthService(jwtSecret string, apiKeys []string) *AuthService {
    return &AuthService{
        jwtSecret: jwtSecret,
        apiKeys:   apiKeys,
    }
}

func (a *AuthService) GenerateToken(userID, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(a.jwtSecret))
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(a.jwtSecret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}

func (a *AuthService) ValidateAPIKey(apiKey string) bool {
    for _, key := range a.apiKeys {
        if key == apiKey {
            return true
        }
    }
    return false
}

