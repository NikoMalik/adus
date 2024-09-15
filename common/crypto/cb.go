package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/curve25519"
)

const (
	keySize       = 32 // Size for X25519 and AES-GCM keys
	nonceSize     = 12 // Size of nonce for AES-GCM
	signatureSize = 64 // Size of Ed25519 signatures
)

var _ Cipher = (*X25519Cipher)(nil)

type X25519Cipher struct {
	key [keySize]byte
}

func (c *X25519Cipher) GenerateKey() (ed25519.PrivateKey, error) {
	privKey := make(ed25519.PrivateKey, ed25519.PrivateKeySize)
	if _, err := io.ReadFull(rand.Reader, privKey); err != nil {
		return nil, err
	}
	copy(c.key[:], privKey[:keySize])
	return privKey, nil
}

// Generate a shared secret using X25519
func (c *X25519Cipher) GenerateSharedSecret(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (*[keySize]byte, error) {
	var sharedSecret [keySize]byte
	xSharedSecret, err := curve25519.X25519(privateKey[:32], publicKey[:32])
	if err != nil {
		return nil, err
	}

	copy(sharedSecret[:], xSharedSecret)
	return &sharedSecret, nil
}

// Encrypt data using AES-GCM
func (c *X25519Cipher) Encrypt(sharedSecret *[keySize]byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(sharedSecret[:32]) // Use 256-bit key for AES create new block cipher
	if err != nil {
		return nil, err
	}
	// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
	// with the standard nonce length.
	//
	// In general, the GHASH operation performed by this implementation of GCM is not constant-time.
	// An exception is when the underlying [Block] was created by aes.NewCipher
	// on systems with hardware support for AES. See the [crypto/aes] package documentation for details.
	aesgcm, err := cipher.NewGCM(block) // create a Galois/Counter Mode
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize()) // generated nonce
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	//
	// To reuse plaintext's storage for the encrypted output, use plaintext[:0]
	// as dst. Otherwise, the remaining capacity of dst must not overlap plaintext.
	ciphertext := aesgcm.Seal(nonce, nonce, data, nil) // Seal encrypt
	return ciphertext, nil
}

// Decrypt data using AES-GCM
func (c *X25519Cipher) Decrypt(sharedSecret *[keySize]byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(sharedSecret[:32]) // Use 256-bit key for AES create new block cipher with existing key
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block) // create cipher for decryption
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize() // get nonce size from aegcm that get the key from block
	if len(ciphertext) < nonceSize {
		return nil, io.ErrUnexpectedEOF
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	// Open decrypts and authenticates ciphertext, authenticates the
	// additional data and, if successful, appends the resulting plaintext
	// to dst, returning the updated slice. The nonce must be NonceSize()
	// bytes long and both it and the additional data must match the
	// value passed to Seal.
	//
	// To reuse ciphertext's storage for the decrypted output, use ciphertext[:0]
	// as dst. Otherwise, the remaining capacity of dst must not overlap plaintext.
	//
	// Even if the function fails, the contents of dst, up to its capacity,
	// may be overwritten.
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil) // decrypt data
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
