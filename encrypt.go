package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/andrewrynhard/go-mask"
)

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
		fmt.Println(err)
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

func keyFromPassword(password []byte) string {
	hasher := md5.New()

	hasher.Write([]byte(password))

	return hex.EncodeToString(hasher.Sum(nil))
}

func hashPassword(password string) (string, error) {
	var key string

	if password != "" {
		key = keyFromPassword([]byte(password))
	} else {
		fmt.Println("Enter your password:")

		k, err := getMaskedInput()
		if err != nil {
			return "", err
		}

		key = keyFromPassword(k)
	}

	return key, nil
}

func decryptFiles(files []string, password string) error {
	key, err := hashPassword(password)
	if err != nil {
		return err
	}

	for _, file := range files {
		file = file + ".aes"
		exists, err := isFile(file)
		if !exists {
			return err
		}

		log.Printf("Decrypting: %s", file)
		err = decryptFile(file, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func encryptFiles(files []string, password string) error {
	key, err := hashPassword(password)
	if err != nil {
		return err
	}

	for _, file := range files {
		exists, err := isFile(file)
		if !exists {
			return err
		}

		log.Printf("Encrypting: %s", file)
		err = encryptFile(file, key)
		if err != nil {
			return err
		}

		err = os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}
