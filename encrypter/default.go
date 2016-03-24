package encrypter

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/andrewrynhard/go-mask"
	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/logger"
)

var (
	defaultPublicKeyOutput  = os.Getenv("HOME") + "/.alterant/pubring.gpg"
	defaultPrivateKeyOutput = os.Getenv("HOME") + "/.alterant/secring.gpg"
)

// DefaultEncryption is a basic encryption handler and is the default
type DefaultEncryption struct {
	logger   *logWrapper.LogWrapper
	Password string
	Private  string
	Public   string
	Remove   bool
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

func savePublicKey(e *openpgp.Entity) error {
	pubKey, err := os.Create(defaultPublicKeyOutput)
	if err != nil {
		return err
	}

	w, err := armor.Encode(pubKey, openpgp.PublicKeyType, nil)
	if err != nil {
		return err
	}

	err = e.Serialize(w)
	if err != nil {
		return err
	}

	w.Close()

	return nil
}

func savePrivateKey(e *openpgp.Entity, pgpCfg *packet.Config) error {
	privKey, err := os.Create(defaultPrivateKeyOutput)
	if err != nil {
		return err
	}

	w, err := armor.Encode(privKey, openpgp.PrivateKeyType, nil)
	if err != nil {
		return err
	}

	err = e.SerializePrivate(w, pgpCfg)
	if err != nil {
		return err
	}

	w.Close()

	return nil
}

// NewKeyPair Creates a new RSA/RSA  private/public key pair with the provided identity details and signs the
// public key with the private key
func NewKeyPair(name string, comment string, email string) error {
	pgpCfg := newPGPConfig()

	e, err := openpgp.NewEntity(name, comment, email, pgpCfg)
	if err != nil {
		return err
	}

	for _, id := range e.Identities {
		err := id.SelfSignature.SignUserId(id.UserId.Id, e.PrimaryKey, e.PrivateKey, pgpCfg)
		if err != nil {
			return err
		}

		// https://github.com/golang/go/issues/12153
		// https://github.com/inversepath/interlock/blob/master/src/openpgp.go#L330
		// FIXES: openpgp: invalid argument: cannot encrypt because no candidate hash functions are compiled in. (Wanted RIPEMD160 in this case.)
		id.SelfSignature.PreferredHash = []uint8{8}
	}

	// SerializePrivate must be called BEFORE Serialize, so we create the private
	// key first. See https://github.com/golang/go/issues/6483
	savePrivateKey(e, pgpCfg)
	savePublicKey(e)

	return nil
}

func signedEntity(privateKey string) (*openpgp.Entity, error) {
	// open ascii armored private key
	sign, err := os.Open(privateKey)
	defer sign.Close()
	if err != nil {
		return nil, err
	}

	// decode armor and check key type
	signBlock, err := armor.Decode(sign)
	if err != nil {
		return nil, err
	}

	if signBlock.Type != openpgp.PrivateKeyType {
		return nil, fmt.Errorf("sign key type:%s", signBlock.Type)
	}

	// parse and decrypt decoded key
	signReader := packet.NewReader(signBlock.Body)
	signed, err := openpgp.ReadEntity(signReader)
	if err != nil {
		return nil, err
	}

	return signed, nil
}

func encryptFile(file string, to *openpgp.EntityList, signed *openpgp.Entity, pgpCfg *packet.Config) error {
	// read the file intended for encryption into a buffer
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	ciphertext := new(bytes.Buffer)

	// create the encryption writer
	w, err := openpgp.Encrypt(ciphertext, *to, signed, nil, pgpCfg)
	if err != nil {
		return err
	}

	// encrypt the data
	_, err = w.Write(content)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(ciphertext)
	if err != nil {
		return err
	}

	// encode the encypted data as a base64 string
	encoded := base64.StdEncoding.EncodeToString(bytes)

	// write the encoded/encypted data to disk
	writeToFile(encoded, file+".encrypted")

	return nil
}

// EncryptFiles encrypts a file
func (de *DefaultEncryption) EncryptFiles(cfg *config.Config) error {
	pgpCfg := newPGPConfig()

	// open ascii armored public key
	f, err := os.Open(de.Public)
	defer f.Close()
	if err != nil {
		return err
	}

	// retrieve the entities in the keyring
	to, err := openpgp.ReadArmoredKeyRing(f)
	if err != nil {
		return err
	}

	// obtain a private key for signing
	signed, err := signedEntity(de.Private)
	if err != nil {
		return err
	}

	for _, task := range cfg.Tasks {
		for _, link := range task.Links {
			file := string(link.Target)
			encrypted := link.Encrypted

			if encrypted {
				exists, err := isFile(file)
				if !exists {
					return err
				}

				de.logger.Info("Encrypting: %s", file)
				err = encryptFile(file, &to, signed, pgpCfg)
				if err != nil {
					return err
				}

				if de.Remove {
					de.logger.Info("Removing: %s", file)
					err = os.Remove(file)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func decryptFile(file string, to *openpgp.EntityList, pgpCfg *packet.Config) error {
	// read the file inteded for decryption into a buffer
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	// decode the base64 encypted data
	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return err
	}

	// decrypt the data
	md, err := openpgp.ReadMessage(bytes.NewBuffer(decoded), to, nil, pgpCfg)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return err
	}

	plaintext := string(bytes)

	// write the decoded/decypted data to disk
	file = strings.TrimSuffix(file, filepath.Ext(file))

	writeToFile(plaintext, file)

	return nil
}

// DecryptFiles decrypts a file
func (de *DefaultEncryption) DecryptFiles(cfg *config.Config) error {
	pgpCfg := newPGPConfig()

	// open the private key file
	privateKeyring, err := os.Open(de.Private)
	defer privateKeyring.Close()
	if err != nil {
		return err
	}

	// retrieve the entities in the keyring
	to, err := openpgp.ReadArmoredKeyRing(privateKeyring)
	if err != nil {
		return err
	}

	for _, task := range cfg.Tasks {
		for _, link := range task.Links {
			file := string(link.Target)
			encrypted := link.Encrypted

			if encrypted {
				file := file + ".encrypted"
				exists, err := isFile(file)
				if !exists {
					return err
				}

				de.logger.Info("Decrypting: %s", file)
				err = decryptFile(file, &to, pgpCfg)
				if err != nil {
					return err
				}

				if de.Remove {
					de.logger.Info("Removing: %s", file)
					err = os.Remove(file)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func newPGPConfig() *packet.Config {
	pgpCfg := &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       4096,
	}

	return pgpCfg
}

// NewDefaultEncryption returns an instance of `DefaultEncryption`
func NewDefaultEncryption(password string, private string, public string, remove bool, logger *logWrapper.LogWrapper) *DefaultEncryption {
	return &DefaultEncryption{
		logger:   logger,
		Password: password,
		Private:  private,
		Public:   public,
		Remove:   remove,
	}
}
