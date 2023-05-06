package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// procedure is as follow:
//
// 1. base64 decode
//
// 2. decrypt aes in CTR mode
//
// 3. remove salt
func Decrypt(cipherText, secret string) (string, error) {
	binaryCipherText, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	binarySecret := []byte(secret)

	// Create new AES cipher block
	block, err := aes.NewCipher(binarySecret)
	if err != nil {
		return "", err
	}

	// Decrpt
	decryptedText := make([]byte, len(binaryCipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, binaryCipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, binaryCipherText[aes.BlockSize:])

	return string(decryptedText), nil
}
