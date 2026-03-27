package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

// IniConfig loads and decrypts an INI configuration file into a struct
// stage: environment stage (e.g., "development", "production")
// secret: encryption secret key
// Returns: parsed config struct of type T
func IniConfig[T any](stage, secret string) T {
	if stage == "" {
		panic("stage is required")
	}
	if secret == "" {
		panic("secret is required")
	}

	encFile := fmt.Sprintf("%s.ini.enc", stage)

	// 1. decrypt
	decryptBytes := DecryptBytesFromFile(encFile, secret)

	// 2. load ini from memory
	cfgFile, err := ini.LoadSources(
		ini.LoadOptions{
			IgnoreInlineComment: true,
		},
		decryptBytes,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to parse ini: %v", err))
	}

	// 3. map to struct
	var cfg T
	if err := cfgFile.MapTo(&cfg); err != nil {
		panic(fmt.Sprintf("failed to map config to struct: %v", err))
	}

	return cfg
}

// DecryptBytesFromFile reads and decrypts bytes from a file
// path: path to the encrypted file
// key: encryption key for decryption
// Returns: decrypted bytes
func DecryptBytesFromFile(path string, key string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("DecryptBytesFromFile: failed to read %s: %w", path, err))
	}

	encrypted := string(data)
	return decryptBytes(encrypted, key)
}

// decryptBytes decrypts AES-GCM encrypted data
// encodedData: base64 URL-encoded encrypted data (nonce + ciphertext)
// key: encryption secret key
// Returns: decrypted plaintext bytes
func decryptBytes(encodedData string, key string) []byte {
	stretchedKey := keyStretch(key)

	dataBytes, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(stretchedKey)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	if len(dataBytes) < nonceSize {
		panic("ciphertext too short")
	}

	nonce, ciphertext := dataBytes[:nonceSize], dataBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return plaintext
}

// keyStretch derives a 32-byte key from a string secret using PBKDF2
// This is a placeholder - implement based on your key derivation needs
// For now, it uses SHA256 to expand the key to 32 bytes
func keyStretch(key string) []byte {
	hash := sha256.Sum256([]byte(key)) // Use SHA-256 to stretch the key
	return hash[:32]                   // Return a 32-byte key for AES-256
}
