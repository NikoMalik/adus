package crypto

type Cipher interface {
	Encrypt(*[keySize]byte, []byte) ([]byte, error)
	Decrypt(*[keySize]byte, []byte) ([]byte, error)
}
