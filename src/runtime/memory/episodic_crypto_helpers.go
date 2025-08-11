package runtimememory

import (
    "errors"

    runtimesecurity "github.com/contexis-cmp/contexis/src/runtime/security"
)

// securityKeyProvider returns a closure that yields the AES key bytes for episodic encryption.
func securityKeyProvider() func() ([]byte, error) {
    provider := runtimesecurity.EnvKeyProvider{}
    return func() ([]byte, error) {
        // tenant-aware keying could be added by passing through context/tenant id here
        key, err := provider.GetKey("")
        if err != nil {
            return nil, err
        }
        if len(key) != 32 {
            return nil, errors.New("episodic key must be 32 bytes")
        }
        return key, nil
    }
}

func encryptBytes(key []byte, plaintext []byte) ([]byte, error) {
    return runtimesecurity.EncryptGCM(key, plaintext)
}

func decryptBytes(key []byte, data []byte) ([]byte, error) {
    return runtimesecurity.DecryptGCM(key, data)
}


