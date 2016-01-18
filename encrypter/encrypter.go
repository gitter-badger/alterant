package encrypter

import "github.com/autonomy/alterant/config"

// Encrypter is the interface for encryption
type Encrypter interface {
	HashPassword(string) (string, error)
	DecryptFiles(*config.Config) error
	EncryptFiles(*config.Config) error
}
