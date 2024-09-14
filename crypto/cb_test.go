package crypto

import (
	"crypto/ed25519"
	"testing"
)

func TestCrypto(t *testing.T) {
	c := &X25519Cipher{}

	data := make([]byte, 1024)

	edata, err := c.Encrypt(&c.key, data)

	if err != nil {
		t.Error(err)
	}
	ddata, err := c.Decrypt(&c.key, edata)
	if err != nil {
		t.Error(err)
	}

	if !equal(data, ddata) {
		t.Error("data not equal")
	}
}

func TestX25519Cipher(t *testing.T) {

	c := &X25519Cipher{}
	privKey, _ := c.GenerateKey()
	pubKey := privKey.Public().(ed25519.PublicKey)

	sharedSecret, err := c.GenerateSharedSecret(privKey, pubKey)
	if err != nil {
		t.Fatalf("GenerateSharedSecret() error: %v", err)
	}

	data := []byte("Test data for encryption")
	encryptedData, err := c.Encrypt(sharedSecret, data)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	decryptedData, err := c.Decrypt(sharedSecret, encryptedData)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if !equal(data, decryptedData) {
		t.Errorf("Decrypted data does not match original data. Got %s, want %s", decryptedData, data)
	}
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
