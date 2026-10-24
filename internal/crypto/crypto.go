package crypto

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

// AEAD wraps an AEAD cipher instance.
type AEAD struct {
	aead cipher
}

// Using an interface to make testing easier.
type cipher interface {
	NonceSize() int
	Overhead() int
	Seal(dst, nonce, plaintext, additionalData []byte) []byte
	Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error)
}

// NewAEADKey initializes ChaCha20-Poly1305 AEAD with key.
func NewAEADKey(key []byte) (*AEAD, error) {
	if len(key) < chacha20poly1305.KeySize {
		// For simplicity accept keys of length < 32 by zero-padding (not recommended for prod)
		return nil, errors.New("encryption key must be at least 32 bytes")
	}
	c, err := chacha20poly1305.NewX(key[:chacha20poly1305.KeySize])
	if err != nil {
		return nil, err
	}
	return &AEAD{aead: c}, nil
}

// Encrypt returns ciphertext and nonce
func (a *AEAD) Encrypt(plaintext []byte) (ciphertext []byte, nonce []byte, err error) {
	n := a.aead.NonceSize()
	nonce = make([]byte, n)
	if _, err = rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	ciphertext = a.aead.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// Decrypt returns plaintext or an error
func (a *AEAD) Decrypt(ciphertext, nonce []byte) (plaintext []byte, err error) {
	return a.aead.Open(nil, nonce, ciphertext, nil)
}
