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

type KeyPair struct {
	priv *PrivateKey
	pub  *PublicKey
}

func (k *KeyPair) PrivateKey() *PrivateKey {
	return k.priv
}

func (k *KeyPair) PublicKey() *PublicKey {
	return k.pub
}

func NewKeyPair(priv *PrivateKey) *KeyPair {
	return &KeyPair{priv: priv, pub: GetPublicKeyFromPrivate(priv)}
}

type ViewOnlyKeyPair struct {
	view  *KeyPair
	spend *PublicKey
}

func (p *ViewOnlyKeyPair) ViewKeyPair() *KeyPair {
	return p.view
}

func (p *ViewOnlyKeyPair) SpendPublicKey() *PublicKey {
	return p.spend
}

func NewViewOnlyKeyPair(view *PrivateKey, spend *PublicKey) *ViewOnlyKeyPair {
	return &ViewOnlyKeyPair{view: NewKeyPair(view), spend: spend}
}

type FullKeyPair struct {
	view  *KeyPair
	spend *KeyPair
}

func (p *FullKeyPair) ViewKeyPair() *KeyPair {
	return p.view
}

func (p *FullKeyPair) SpendKeyPair() *KeyPair {
	return p.spend
}

func (p *FullKeyPair) ViewOnlyKeyPair() *ViewOnlyKeyPair {
	return &ViewOnlyKeyPair{view: p.view, spend: p.spend.pub}
}

func NewFullKeyPair(view *PrivateKey, spend *PrivateKey) *FullKeyPair {
	return &FullKeyPair{view: NewKeyPair(view), spend: NewKeyPair(spend)}
}

func NewFullKeyPairSpendPrivateKey(spend *PrivateKey) (*FullKeyPair, error) {
	view, err := GetPrivateViewKeyFromPrivateSpendKey(spend)
	if err != nil {
		return nil, err
	}

	return &FullKeyPair{view: NewKeyPair(view), spend: NewKeyPair(spend)}, nil
}
