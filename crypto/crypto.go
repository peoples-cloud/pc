package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"

	b58 "github.com/jbenet/go-base58"
)

func fetchTextAndBlock(source string, key []byte) ([]byte, cipher.Block) {
	// read content from your file
	text, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return text, block
}

func saveFile(text []byte, destination string) {
	// create a new file for saving the decrypted data.
	f, err := os.Create(destination)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f, bytes.NewReader(text))
	if err != nil {
		panic(err.Error())
	}
}

// generating functions from: https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// func generateRandomString(s int) (string, error) {
// 	b, err := GenerateRandomBytes(s)
// 	fmt.Println(len(b), " -> ", b)
// }
func GenerateRandomString(s int) string {
	b, err := generateRandomBytes(s)
	if err != nil {
		panic(err)
	}
	return b58.Encode(b)
}

func Encrypt(source string) (string, string) {
	destination := source + ".aes"
	aeskey, err := generateRandomBytes(32)
	if err != nil {
		panic(err)
	}
	plaintext, block := fetchTextAndBlock(source, aeskey)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	saveFile(ciphertext, destination)

	keystring := b58.Encode(aeskey)
	return keystring, destination
}

func Decrypt(source, key, destination string) {
	// aeskey, err := base64.URLEncoding.DecodeString(key)
	aeskey := b58.Decode(key)
	ciphertext, block := fetchTextAndBlock(source, aeskey)
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	iv := ciphertext[:aes.BlockSize]
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	saveFile(plaintext, destination)
}

// func main() {
// 	source := os.Args[1]
// 	dest := os.Args[2]
// 	// generateRandomString(32)
// 	base64key := encrypt(source, dest)
// 	fmt.Println(base64key)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// decrypt(dest, base64key, "wow-decrypted-crypto.tar.gz")
// }
