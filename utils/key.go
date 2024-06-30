package utils

import (
	"encoding/hex"

	"filippo.io/edwards25519"
)

type PrivateKey struct {
	key *edwards25519.Scalar
}

func (k *PrivateKey) Bytes() []byte {
	return k.key.Bytes()
}

func newPrivateKeyHelper(keyBytes []byte) (*PrivateKey, error) {
	key, err := new(edwards25519.Scalar).SetCanonicalBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{key: key}, nil
}

// Creates Private Key from a hex string representation
func NewPrivateKey(keyStr string) (*PrivateKey, error) {
	keyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}

	return newPrivateKeyHelper(keyBytes)
}

type PublicKey struct {
	key *edwards25519.Point
}

func (k *PublicKey) Bytes() []byte {
	return k.key.Bytes()
}

func newPublicKeyHelper(keyBytes []byte) (*PublicKey, error) {
	key, err := new(edwards25519.Point).SetBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return &PublicKey{key: key}, nil
}

// Creates Public Key from a hex string representation
func NewPublicKey(keyStr string) (*PublicKey, error) {
	keyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}

	return newPublicKeyHelper(keyBytes)
}
