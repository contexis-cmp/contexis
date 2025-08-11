package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "io"
    "os"
)

// KeyProvider abstracts encryption key retrieval (per-tenant)
type KeyProvider interface {
    // GetKey returns a 32-byte key for AES-256
    GetKey(tenantID string) ([]byte, error)
}

// EnvKeyProvider loads key material from environment variables
// CMP_EPISODIC_KEY (hex/base64 not required here: expect 32-byte raw via env or we derive from sha256)
type EnvKeyProvider struct{}

func (EnvKeyProvider) GetKey(tenantID string) ([]byte, error) {
    // For demo/dev: single key for all tenants; production should use KMS per tenant
    k := os.Getenv("CMP_EPISODIC_KEY")
    if len(k) == 32 {
        return []byte(k), nil
    }
    if k == "" {
        return nil, errors.New("CMP_EPISODIC_KEY not set")
    }
    // If key length is not 32, derive a 32-byte key by hashing
    h := sha256Sum(k)
    // take first 32 bytes from hex-decoded hash
    // h is hex, length 64; convert to bytes: naive approach here for brevity
    // we can just use the hex string bytes but ensure 32 bytes
    if len(h) < 32 {
        return nil, errors.New("derived key too short")
    }
    return []byte(h[:32]), nil
}

// EncryptGCM encrypts plaintext with AES-GCM using the provided key
func EncryptGCM(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    // result: nonce || ciphertext
    out := make([]byte, 0, len(nonce)+len(plaintext)+gcm.Overhead())
    out = append(out, nonce...)
    out = append(out, gcm.Seal(nil, nonce, plaintext, nil)...)
    return out, nil
}

// DecryptGCM decrypts data produced by EncryptGCM
func DecryptGCM(key, data []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    if len(data) < gcm.NonceSize() {
        return nil, errors.New("ciphertext too short")
    }
    nonce := data[:gcm.NonceSize()]
    ciphertext := data[gcm.NonceSize():]
    return gcm.Open(nil, nonce, ciphertext, nil)
}


