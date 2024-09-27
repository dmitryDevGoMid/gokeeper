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
	"path/filepath"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/checkdir"
	//"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
)

const pathGokeeperKey = "gokeeperspace/gokeeperkey"

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
	AllSet()
	ReadServerPublicKey() (*rsa.PublicKey, error)
	EncryptByServerKey(msg string) ([]byte, error)
	SetPublicServerKey(key string) error
	EncryptByServerKeyParts(msg string) ([]byte, error)
	DecryptOAEP(msg []byte) ([]byte, error)
	DecryptOAEPServer(msg []byte) ([]byte, error)

	ReadAndGetClientPublicKey(key string) (*rsa.PublicKey, error)
	EncryptByClientKeyParts(msg string, key string) ([]byte, error)
}

type AsimEncryptStruct struct {
	//cfg            *config.Config
	PathEncryptKey string
	PublicKey      *rsa.PublicKey
	PrivateKey     *rsa.PrivateKey

	PublicKeyByte  []byte
	PrivateKeyByte []byte

	PublicServerKey  []byte
	PrivateServerKey *rsa.PrivateKey
}

func NewAsimEncrypt() *AsimEncryptStruct {
	//cfg *config.Config) AsimEncryptStruct {
	return &AsimEncryptStruct{}
}

func (asme *AsimEncryptStruct) GetBytePrivate() []byte {
	return asme.PrivateKeyByte
}

func (asme *AsimEncryptStruct) GetBytePublic() []byte {
	return asme.PublicKeyByte
}

func (asme *AsimEncryptStruct) AllSet() {
	asme.SetPrivateKey()
	asme.SetPublicKey()
	asme.ReadPrivateKeyGetByte()
	asme.ReadPublicKeyGetByte()
}

// GenerateRsaKeyPair generates an RSA key pair and returns the private and public keys
func (asme *AsimEncryptStruct) GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func (asme *AsimEncryptStruct) SetPublicServerKey(key string) error {
	asme.PublicServerKey = []byte(key)

	return nil
}

// Get path to keys
func (asme *AsimEncryptStruct) GetPathToKey() string {
	PathEncryptKey := "/" //asme.cfg.PathEncrypt.PathEncryptKey

	if PathEncryptKey != "" {

		d, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		return d + "/" + PathEncryptKey
	}
	return ""

}

// Get path to keys
func (asme *AsimEncryptStruct) GetPathOnly(name string) string {

	d, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return d + "/" + name

}

func (asme *AsimEncryptStruct) GetPathExecutable() (string, error) {
	// Получаем путь к исполняемому файлу
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return "", nil
	}

	// Получаем путь к директории, содержащей исполняемый файл
	dirPath := filepath.Dir(executablePath)
	return dirPath, nil
}

func (asme *AsimEncryptStruct) CheckCreateFileToDirectory() {

	//dirPath, err := asme.GetPathExecutable()
	dirPath, _, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, "")
	if err != nil {
		return
	}
	//fmt.Println("dirPath CheckCreateFileToDirectory", dirPath)
	//time.Sleep(10 * time.Second)
	// Проверяем права доступа на запись в текущую директорию
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	if fileInfo.IsDir() {
		if fileInfo.Mode().Perm()&os.ModePerm&0222 == 0 {
			fmt.Println("No write permission in current directory")
		} else {
			fmt.Println("Write permission in current directory")
		}
	} else {
		fmt.Println("Error: current path is not a directory")
	}
}

func (asme *AsimEncryptStruct) CheckFile(name string) error {
	// Проверяем, существует ли файл
	//currentPath, err := asme.GetPathExecutable()
	currentPath, _, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, name)
	if err != nil {
		return err
	}

	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}

	file := fmt.Sprintf("%s/%s", currentPath, name)

	_, err = os.Stat(file)
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
func (asme *AsimEncryptStruct) SetPublicKey() error {
	//Check config by private key for decode body
	//PathEncryptKey := asme.GetPathOnly("keeper_public.pem")
	_, PathEncryptKey, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, "keeper_public.pem")
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}

	if PathEncryptKey != "" {
		PublicKey, err := asme.ReadPublicKey(PathEncryptKey)
		if err != nil {
			return err
		}
		asme.PublicKey = PublicKey
		asme.PathEncryptKey = PathEncryptKey

		//asme.cfg.PathEncrypt.KeyEncryptEnbled = true
		//fmt.Println(asme.PublicKey)
	}

	return nil
}

// Set Private key
func (asme *AsimEncryptStruct) SetPrivateKey() error {
	//Check config by private key for decode body
	//PathEncryptKey := asme.GetPathOnly("keeper_private.pem")
	_, PathEncryptKey, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, "keeper_private.pem")
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}

	if PathEncryptKey != "" {
		PrivateKey, err := asme.ReadPrivateKey(PathEncryptKey)
		if err != nil {
			return err
		}
		asme.PrivateKey = PrivateKey
		asme.PathEncryptKey = PathEncryptKey

		//asme.cfg.PathEncrypt.KeyEncryptEnbled = true
		//fmt.Println(asme.PrivateKey)
	}

	return nil
}

// GenerateKey generates a new RSA key pair and saves the private key to a file
func (asme *AsimEncryptStruct) GenerateKeyLink() (*rsa.PrivateKey, *rsa.PublicKey) {
	// Generate a new RSA key pair
	priv, public := asme.GenerateRsaKeyPair()

	return priv, public
}

// GenerateKey generates a new RSA key pair and saves the private key to a file
func (asme *AsimEncryptStruct) GenerateKeyFile(prefix string) error {

	asme.CheckCreateFileToDirectory()

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

	currentPath, _, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, "")
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}
	/*currentPath, err := asme.GetPathExecutable()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}*/

	fileCreatePrivate := fmt.Sprintf("%s/%s", currentPath, prefix+"_private.pem")
	filePrivate, err := os.Create(fileCreatePrivate)
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
	fileCreatePublic := fmt.Sprintf("%s/%s", currentPath, prefix+"_public.pem")
	filePublic, err := os.Create(fileCreatePublic)
	if err != nil {
		return err
	}
	defer filePublic.Close()
	err = pem.Encode(filePublic, blockPublic)
	if err != nil {
		return err
	}

	//fmt.Println("Private Key : ", priv)
	//fmt.Println("Public key ", &priv.PublicKey)

	return nil
}

func (asme *AsimEncryptStruct) ReadPublicKeyGetByte() error {
	// Read the file
	//file := asme.GetPathOnly("keeper_public.pem")
	_, file, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, "keeper_public.pem")
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}
	//fmt.Println(file)
	data, err := os.ReadFile(file) // changed
	if err != nil {
		return err
	}

	asme.PublicKeyByte = data

	return nil
}

// ReadPublicKey reads a public key from a file
func (asme *AsimEncryptStruct) ReadServerPublicKey() (*rsa.PublicKey, error) {

	// Decode the file from PEM format
	block, _ := pem.Decode(asme.PublicServerKey)
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
func (asme *AsimEncryptStruct) ReadPublicKey(filename string) (*rsa.PublicKey, error) {
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

func (asme *AsimEncryptStruct) ReadPrivateKeyGetByte() error {
	// Read the file
	//file := asme.GetPathOnly("keeper_private.pem")
	_, file, _, err := checkdir.EnsureDirectoryExists(pathGokeeperKey, "keeper_private.pem")
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}
	//fmt.Println(file)
	data, err := os.ReadFile(file) // changed
	if err != nil {
		return err
	}

	asme.PublicKeyByte = data

	return nil
}

// ReadPrivateKey reads a private key from a file
func (asme *AsimEncryptStruct) ReadPrivateKey(filename string) (*rsa.PrivateKey, error) {
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

// ReadPublicKey reads a public key from a file
func (asme *AsimEncryptStruct) ReadAndGetClientPublicKey(key string) (*rsa.PublicKey, error) {

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

func (asme *AsimEncryptStruct) EncryptByClientKeyParts(msg string, key string) ([]byte, error) {
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

func (asme *AsimEncryptStruct) EncryptByServerKeyParts(msg string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function
	pub, _ := asme.ReadServerPublicKey()
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
// func (asme *AsimEncryptStruct) Encrypt(pub *rsa.PublicKey, msg string) ([]byte, error) {
func (asme *AsimEncryptStruct) EncryptByServerKey(msg string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function

	pub, _ := asme.ReadServerPublicKey()

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

// Encrypt encrypts a message using the public key and the OAEP method
// func (asme *AsimEncryptStruct) Encrypt(pub *rsa.PublicKey, msg string) ([]byte, error) {
func (asme *AsimEncryptStruct) Encrypt(msg string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function

	pub := asme.PublicKey

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

// Decrypt decrypts an encrypted message using the private key and the OAEP method
// func (asme *AsimEncryptStruct) Decrypt(priv *rsa.PrivateKey, ciphertext []byte) (string, error) {
func (asme *AsimEncryptStruct) Decrypt(ciphertext []byte) (string, error) {
	// Use an empty label (label) and the SHA-256 hash function
	priv := asme.PrivateKey

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

func (asme *AsimEncryptStruct) DecryptOAEP(msg []byte) ([]byte, error) {

	private := asme.PrivateKey

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

func (asme *AsimEncryptStruct) DecryptOAEPServer(msg []byte) ([]byte, error) {

	private := asme.PrivateServerKey

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
