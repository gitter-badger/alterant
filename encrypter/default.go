package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/andrewrynhard/go-mask"
	"github.com/autonomy/alterant/logger"
)

// DefaultEncryption is a basic encryption method and is the default
type DefaultEncryption struct {
	logger   *logWrapper.LogWrapper
	Password string
	Remove   bool
}

// TODO: verify that the password works
func decrypt(cipherstring string, keystring string) string {
	ciphertext := []byte(cipherstring)

	key := []byte(keystring)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func encrypt(plainstring, keystring string) string {
	plaintext := []byte(plainstring)

	key := []byte(keystring)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return string(ciphertext)
}

func isFile(file string) (bool, error) {
	s, err := os.Stat(file)

	if os.IsNotExist(err) {
		return false, err
	}

	if s.IsDir() {
		return false, err
	}

	return true, nil
}

func writeToFile(data, file string) {
	ioutil.WriteFile(file, []byte(data), 0666)
}

func readFromFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	return data, err
}

func decryptFile(file string, key string) error {
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	decrypted := decrypt(string(content), string(key))
	writeToFile(decrypted, file[:len(file)-4])

	return nil
}

func encryptFile(file string, key string) error {
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	encrypted := encrypt(string(content), string(key))
	writeToFile(encrypted, file+".aes")

	return nil
}

func getMaskedInput() ([]byte, error) {
	maskedReader := mask.NewMaskedReader()

	key, err := maskedReader.GetInputConfirmMasked()
	if err != nil {
		return nil, err
	}

	return key, nil
}

func key32BitFromPassword(password []byte) string {
	hasher := md5.New()

	hasher.Write([]byte(password))

	return hex.EncodeToString(hasher.Sum(nil))
}

// HashPassword creates an md5 sum from a string, ensuring a 32 byte key
func (de *DefaultEncryption) HashPassword(password string) (string, error) {
	var key string

	if password != "" {
		key = key32BitFromPassword([]byte(password))
	} else {
		fmt.Println("Enter your password:")

		k, err := getMaskedInput()
		if err != nil {
			return "", err
		}

		key = key32BitFromPassword(k)
	}

	return key, nil
}

// DecryptFiles decrypts `files` based on the hash of `password`
func (de *DefaultEncryption) DecryptFiles(files []string) error {
	key, err := de.HashPassword(de.Password)
	if err != nil {
		return err
	}

	for _, file := range files {
		file = file + ".aes"
		exists, err := isFile(file)
		if !exists {
			return err
		}

		de.logger.Info("Decrypting: %s", file)
		err = decryptFile(file, key)
		if err != nil {
			return err
		}

		if de.Remove {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// EncryptFiles encrypts `files` based on the hash of `password`
func (de *DefaultEncryption) EncryptFiles(files []string) error {
	key, err := de.HashPassword(de.Password)
	if err != nil {
		return err
	}

	for _, file := range files {
		exists, err := isFile(file)
		if !exists {
			return err
		}

		de.logger.Info("Encrypting: %s", file)
		err = encryptFile(file, key)
		if err != nil {
			return err
		}

		if de.Remove {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewDefaultEncryption returns an instance of `DefaultEncryption`
func NewDefaultEncryption(password string, remove bool, logger *logWrapper.LogWrapper) *DefaultEncryption {
	return &DefaultEncryption{
		logger:   logger,
		Password: password,
		Remove:   remove,
	}
}
