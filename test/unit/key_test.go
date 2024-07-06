package test

import (
	"encoding/hex"
	"testing"

	"github.com/chekist32/go-monero/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicKeyFromPrivate(t *testing.T) {
	privKey, err := utils.NewPrivateKey("ac413c16b815899b69393d72086fa86d31e8e352895606180c4c8fadd707450a")
	if err != nil {
		t.Fatal(err)
	}

	pub := utils.GetPublicKeyFromPrivate(privKey)

	expected, err := utils.NewPublicKey("0ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Bytes(), pub.Bytes())
}

func TestPublicKey(t *testing.T) {
	expected, err := hex.DecodeString("0ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := utils.NewPublicKey("0ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b")
	assert.NoError(t, err)

	assert.Equal(t, expected, actual.Bytes())
}

func TestPrivateKey(t *testing.T) {
	expected, err := hex.DecodeString("ac413c16b815899b69393d72086fa86d31e8e352895606180c4c8fadd707450a")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := utils.NewPrivateKey("ac413c16b815899b69393d72086fa86d31e8e352895606180c4c8fadd707450a")
	assert.NoError(t, err)

	assert.Equal(t, expected, actual.Bytes())
}
