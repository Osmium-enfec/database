package util

import (
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

// decryptBytes is a placeholder for the actual decryption logic
// Replace this with your actual decryption implementation
func decryptBytes(encrypted, key string) []byte {
	// TODO: Implement actual decryption logic
	// This is a placeholder - implement based on your encryption scheme
	return []byte(encrypted)
}
