package encrypter

// Encrypter is the interface for encryption
type Encrypter interface {
	HashPassword(string) (string, error)
	DecryptFiles([]string) error
	EncryptFiles([]string) error
}
