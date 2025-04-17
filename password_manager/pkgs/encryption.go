package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func EncryptPass(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptPass(cipherText string, key []byte) (string, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if(err != nil){
		return "", err
	}
	block, err := aes.NewCipher(key)

	if(err != nil){
		return "", err
	}

	gcm, err := cipher.NewGCM(block)

	if(err != nil){
		return "", err
	}

	nonceSize := gcm.NonceSize()
	
	if len(cipherTextBytes) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}