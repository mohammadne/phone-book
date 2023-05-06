package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// procedure is as follow:
//
// 1. plainText + salt -> binary format
//
// 2. encrypt aes in CTR mode
//
// 3. base64 encode
func Encrypt(plainText, secret string) (string, error) {
	binaryText, binarySecret := []byte(plainText), []byte(secret)

	// Create new AES cipher block
	block, err := aes.NewCipher(binarySecret)
	if err != nil {
		return "", err
	}

	// The IV (Initialization Vector) need to be unique, but not secure.
	// Therefore, it's common to include it at the beginning of the cipher text.
	cipherText := make([]byte, aes.BlockSize+len(binaryText))

	// Create IV
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], binaryText)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}
