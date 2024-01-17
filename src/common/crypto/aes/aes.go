package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type GcmProvider struct {
	cipher cipher.AEAD
}

func NewAesGcmProvider(key string) (*GcmProvider, error) {
	hexKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes gcm provider %w", err)
	}
	block, err := aes.NewCipher(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to new block %w", err)
	}

	c, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to new cipher %w", err)
	}
	return &GcmProvider{cipher: c}, nil
}

func (a *GcmProvider) Seal(plain []byte) ([]byte, error) {
	nonceSize := a.cipher.NonceSize()
	nonce := make([]byte, nonceSize)

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to create random nonce %w", err)
	}

	return a.cipher.Seal(nonce, nonce, plain, nil), nil

}

func (a *GcmProvider) Open(enc []byte) ([]byte, error) {
	nonceSize := a.cipher.NonceSize()

	if len(enc) <= nonceSize {
		return nil, fmt.Errorf("invalid enc, current length %v", len(enc))
	}
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plain, err := a.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open gcm %w", err)
	}

	return plain, nil
}

func (a *GcmProvider) Compare(plain, enc []byte) (bool, error) {
	nonceSize := a.cipher.NonceSize()

	if len(enc) <= nonceSize {
		return false, fmt.Errorf("invalid enc, current length %v", len(enc))
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plainOfEnc, err := a.cipher.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return false, fmt.Errorf("failed to open gcm %w", err)
	}

	return bytes.Equal(plain, plainOfEnc), nil
}
