package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	a, err := NewAEADKey(key)
	if err != nil {
		t.Fatalf("failed to create aead: %v", err)
	}
	pt := []byte("hello privacy")
	ct, nonce, err := a.Encrypt(pt)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	out, err := a.Decrypt(ct, nonce)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(out) != string(pt) {
		t.Fatalf("plaintext mismatch: %s != %s", out, pt)
	}
}
