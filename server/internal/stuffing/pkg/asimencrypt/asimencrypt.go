package asimencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
)

type AsimEncrypt interface {
	SetPrivateKey() error
	SetPublicKey() error
	GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey)
	Decrypt(ciphertext []byte) (string, error)
	GenerateKeyLink() (*rsa.PrivateKey, *rsa.PublicKey)
	GenerateKeyFile(prefix string) error
	ReadPublicKey(filename string) (*rsa.PublicKey, error)
	ReadPrivateKey(filename string) (*rsa.PrivateKey, error)
	Encrypt(msg string) ([]byte, error)
	CheckFile(name string) error
	ReadPrivateKeyGetByte() error
	ReadPublicKeyGetByte() error
	GetBytePrivate() []byte
	GetBytePublic() []byte
	AllSet() error
	ReadClientPublicKey() (*rsa.PublicKey, error)
	EncryptByClientKey(msg string, key string) ([]byte, error)
	DecryptOAEP(msg []byte) ([]byte, error)
	EncryptByClientKeyParts(msg string, key string) ([]byte, error)
	DecryptOAEPClient(msg []byte) ([]byte, error)
}

type asimEncrypt struct {
	pathEncryptKey string
	publicKey      *rsa.PublicKey
	privateKey     *rsa.PrivateKey

	publicKeyByte  []byte
	privateKeyByte []byte

	PublicClientKey  []byte
	PrivateClientKey *rsa.PrivateKey
}

func NewAsimEncrypt() *asimEncrypt {
	return &asimEncrypt{}
}

func (asme *asimEncrypt) GetBytePrivate() []byte {
	return asme.privateKeyByte
}

func (asme *asimEncrypt) GetBytePublic() []byte {
	return asme.publicKeyByte
}

func (asme *asimEncrypt) AllSet() error {
	err := asme.SetPrivateKey()
	if err != nil {
		return err
	}

	err = asme.SetPublicKey()
	if err != nil {
		return err
	}

	err = asme.ReadPrivateKeyGetByte()
	if err != nil {
		return err
	}

	err = asme.ReadPublicKeyGetByte()

	if err != nil {
		return err
	}

	return nil
}

// GenerateRsaKeyPair generates an RSA key pair and returns the private and public keys
func (asme *asimEncrypt) GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func (asme *asimEncrypt) SetPublicServerKey(key string) error {
	asme.PublicClientKey = []byte(key)

	return nil
}

// Get path to keys
func (asme *asimEncrypt) GetPathToKey() string {
	pathEncryptKey := "/"

	if pathEncryptKey != "" {

		d, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		return d + "/" + pathEncryptKey
	}
	return ""

}

// Get path to keys
func (asme *asimEncrypt) GetPathOnly(name string) string {

	d, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return d + "/" + name

}

func (asme *asimEncrypt) CheckFile(name string) error {
	// Проверяем, существует ли файл
	_, err := os.Stat(fmt.Sprintf("./%s", name))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(err)
			log.Println("Файл не существует")
			return nil
		} else {
			fmt.Println(err)
			log.Printf("Произошла ошибка при проверке существования файла: %v", err)
			return errors.New("An error occurred while checking the existence of a file")
		}
	} else {
		fmt.Println(err)
		log.Println("Файл существует")
		return errors.New("The file exists! Overwriting will result in loss of information!")
	}
}

// Set Public key
func (asme *asimEncrypt) SetPublicKey() error {
	//Check config by private key for decode body
	pathEncryptKey := asme.GetPathOnly("keeper_public.pem")

	if pathEncryptKey != "" {
		publicKey, err := asme.ReadPublicKey(pathEncryptKey)
		if err != nil {
			return err
		}
		asme.publicKey = publicKey
		asme.pathEncryptKey = pathEncryptKey
	}

	return nil
}

// Set Private key
func (asme *asimEncrypt) SetPrivateKey() error {
	//Check config by private key for decode body
	pathEncryptKey := asme.GetPathOnly("keeper_private.pem")

	if pathEncryptKey != "" {
		privateKey, err := asme.ReadPrivateKey(pathEncryptKey)
		if err != nil {
			return err
		}
		asme.privateKey = privateKey
		asme.pathEncryptKey = pathEncryptKey
	}

	return nil
}

// GenerateKey generates a new RSA key pair and saves the private key to a file
func (asme *asimEncrypt) GenerateKeyLink() (*rsa.PrivateKey, *rsa.PublicKey) {
	// Generate a new RSA key pair
	priv, public := asme.GenerateRsaKeyPair()

	return priv, public
}

// GenerateKey generates a new RSA key pair and saves the private key to a file
func (asme *asimEncrypt) GenerateKeyFile(prefix string) error {

	pri, pub := false, false
	err := asme.CheckFile(prefix + "_private.pem")

	if err != nil {
		pri = true
	}

	err = asme.CheckFile(prefix + "_public.pem")

	if err != nil {
		pub = true
	}

	if (!pri && pub) || (pri && !pub) {
		return errors.New("Oдин из ключей отсутствует!")
	}

	if pri && pub {
		return nil
	}

	// Generate a new RSA key pair
	priv, public := asme.GenerateRsaKeyPair()

	// Encode the private key in PKCS#1 format
	derStream := x509.MarshalPKCS1PrivateKey(priv)

	// Write the private key to a file in PEM format
	blockPrivate := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	filePrivate, err := os.Create(prefix + "_private.pem")
	if err != nil {
		return err
	}
	defer filePrivate.Close()
	err = pem.Encode(filePrivate, blockPrivate)
	if err != nil {
		return err
	}

	// Encode the public key in PKIX format
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return err
	}

	// Write the public key to a file in PEM format
	blockPublic := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	filePublic, err := os.Create(prefix + "_public.pem")
	if err != nil {
		return err
	}
	defer filePublic.Close()
	err = pem.Encode(filePublic, blockPublic)
	if err != nil {
		return err
	}

	return nil
}

func (asme *asimEncrypt) ReadPublicKeyGetByte() error {
	// Read the file
	file := asme.GetPathOnly("keeper_public.pem")
	data, err := os.ReadFile(file) // changed
	if err != nil {
		return err
	}

	asme.publicKeyByte = data

	return nil
}

// ReadPublicKey reads a public key from a file
func (asme *asimEncrypt) ReadAndGetClientPublicKey(key string) (*rsa.PublicKey, error) {

	// Decode the file from PEM format
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse deocde server key")
	}

	// Parse the public key in PKIX format
	publick, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publick.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

// ReadPublicKey reads a public key from a file
func (asme *asimEncrypt) ReadClientPublicKey() (*rsa.PublicKey, error) {

	// Decode the file from PEM format
	block, _ := pem.Decode(asme.PublicClientKey)
	if block == nil {
		return nil, errors.New("failed to parse deocde server key")
	}

	// Parse the public key in PKIX format
	publick, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publick.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

// ReadPublicKey reads a public key from a file
func (asme *asimEncrypt) ReadPublicKey(filename string) (*rsa.PublicKey, error) {
	// Read the file
	data, err := os.ReadFile(filename) // changed
	if err != nil {
		return nil, err
	}

	// Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	// Parse the public key in PKIX format
	publick, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publick.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

func (asme *asimEncrypt) ReadPrivateKeyGetByte() error {
	// Read the file
	file := asme.GetPathOnly("keeper_private.pem")
	data, err := os.ReadFile(file) // changed
	if err != nil {
		return err
	}

	asme.publicKeyByte = data

	return nil
}

// ReadPrivateKey reads a private key from a file
func (asme *asimEncrypt) ReadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	// Read the file
	data, err := os.ReadFile(filename) // changed
	if err != nil {
		return nil, err
	}

	// Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	// Parse the private key in PKCS#1 format
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// Encrypt encrypts a message using the public key and the OAEP method
// func (asme *asimEncrypt) Encrypt(pub *rsa.PublicKey, msg string) ([]byte, error) {
func (asme *asimEncrypt) EncryptByClientKey(msg string, key string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function
	pub, err := asme.ReadAndGetClientPublicKey(key)
	if err != nil {
		log.Println("Error reading client public key: ", err)
		return nil, err
	}

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	// Encrypt the message
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pub,
		[]byte(msg),
		hashedLabel)
	if err != nil {
		log.Println("Error encrypting message", err)
		return nil, err
	}

	return ciphertext, nil
}

func (asme *asimEncrypt) EncryptByClientKeyParts(msg string, key string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function
	pub, err := asme.ReadAndGetClientPublicKey(key)
	if err != nil {
		log.Println("Error reading client public key: ", err)
		return nil, err
	}

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	// Calculate the maximum message size that can be encrypted
	maxMessageSize := pub.Size() - 2*hash.Size() - 2

	// Split the message into parts if it's too large to be encrypted at once
	var ciphertext []byte
	for len(msg) > maxMessageSize {
		part, err := rsa.EncryptOAEP(
			hash,
			rand.Reader,
			pub,
			[]byte(msg[:maxMessageSize]),
			hashedLabel)
		if err != nil {
			return nil, err
		}
		ciphertext = append(ciphertext, part...)
		msg = msg[maxMessageSize:]
	}

	// Encrypt the remaining message part
	part, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pub,
		[]byte(msg),
		hashedLabel)
	if err != nil {
		return nil, err
	}

	ciphertext = append(ciphertext, part...)

	return ciphertext, nil
}

// Encrypt encrypts a message using the public key and the OAEP method
// func (asme *asimEncrypt) Encrypt(pub *rsa.PublicKey, msg string) ([]byte, error) {
func (asme *asimEncrypt) Encrypt(msg string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function

	pub := asme.publicKey

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	// Encrypt the message
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pub,
		[]byte(msg),
		hashedLabel)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (asme *asimEncrypt) DecryptOAEPClient(msg []byte) ([]byte, error) {

	private := asme.PrivateClientKey

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, rand.Reader, private, msg[start:finish], hashedLabel)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func (asme *asimEncrypt) DecryptOAEP(msg []byte) ([]byte, error) {

	private := asme.privateKey

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, rand.Reader, private, msg[start:finish], hashedLabel)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

// Decrypt decrypts an encrypted message using the private key and the OAEP method
// func (asme *asimEncrypt) Decrypt(priv *rsa.PrivateKey, ciphertext []byte) (string, error) {
func (asme *asimEncrypt) Decrypt(ciphertext []byte) (string, error) {
	// Use an empty label (label) and the SHA-256 hash function
	priv := asme.privateKey

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	// Decrypt the message
	plaintext, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		priv,
		ciphertext,
		hashedLabel)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
