package tests

import (
    "secret-manager/internal/services"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestEncryptionDecryption(t *testing.T) {
    encService := services.NewEncryptionService("test-key-32-bytes-long-for-aes256")
    
    plaintext := "This is a secret message"
    
    encrypted, err := encService.Encrypt(plaintext)
    assert.NoError(t, err)
    assert.NotEqual(t, plaintext, encrypted)
    
    decrypted, err := encService.Decrypt(encrypted)
    assert.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}

func TestEncryptionWithDifferentKeys(t *testing.T) {
    encService1 := services.NewEncryptionService("key1-32-bytes-long-for-aes256-enc")
    encService2 := services.NewEncryptionService("key2-32-bytes-long-for-aes256-enc")
    
    plaintext := "Secret message"
    
    encrypted, _ := encService1.Encrypt(plaintext)
    
    _, err := encService2.Decrypt(encrypted)
    assert.Error(t, err)
}
