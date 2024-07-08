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

func TestKeyPair(t *testing.T) {
	exPriv, err := utils.NewPrivateKey("372fcc2abc6bc5015103aae4763822e45c4cfe775d163f97a9ebdd77b0d12c0c")
	if err != nil {
		t.Fatal(err)
	}
	exPub, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	keys := utils.NewKeyPair(exPriv)
	assert.Equal(t, exPriv.Bytes(), keys.PrivateKey().Bytes())
	assert.Equal(t, exPub.Bytes(), keys.PublicKey().Bytes())
}

func TestViewOnlyKeyPair(t *testing.T) {
	exPrivView, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}
	exPubView, err := utils.NewPublicKey("b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a")
	if err != nil {
		t.Fatal(err)
	}
	exPubSpend, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	keys := utils.NewViewOnlyKeyPair(exPrivView, exPubSpend)
	assert.Equal(t, exPrivView.Bytes(), keys.ViewKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exPubView.Bytes(), keys.ViewKeyPair().PublicKey().Bytes())
	assert.Equal(t, exPubSpend.Bytes(), keys.SpendPublicKey().Bytes())
}

func TestGetPrivateViewKeyFromPrivateSpendKey(t *testing.T) {
	exPrivView, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}
	exPrivSpend, err := utils.NewPrivateKey("372fcc2abc6bc5015103aae4763822e45c4cfe775d163f97a9ebdd77b0d12c0c")
	if err != nil {
		t.Fatal(err)
	}

	privView, err := utils.GetPrivateViewKeyFromPrivateSpendKey(exPrivSpend)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exPrivView.Bytes(), privView.Bytes())
}

func TestFullKeyPair(t *testing.T) {
	exPrivView, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}
	exPubView, err := utils.NewPublicKey("b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a")
	if err != nil {
		t.Fatal(err)
	}
	exPrivSpend, err := utils.NewPrivateKey("372fcc2abc6bc5015103aae4763822e45c4cfe775d163f97a9ebdd77b0d12c0c")
	if err != nil {
		t.Fatal(err)
	}
	exPubSpend, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	keys := utils.NewFullKeyPair(exPrivView, exPrivSpend)
	assert.Equal(t, exPrivView.Bytes(), keys.ViewKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exPubView.Bytes(), keys.ViewKeyPair().PublicKey().Bytes())
	assert.Equal(t, exPrivSpend.Bytes(), keys.SpendKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exPubSpend.Bytes(), keys.ViewOnlyKeyPair().SpendPublicKey().Bytes())
}

func TestFullKeyPairSpendPrivateKey(t *testing.T) {
	exPrivView, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}
	exPubView, err := utils.NewPublicKey("b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a")
	if err != nil {
		t.Fatal(err)
	}
	exPrivSpend, err := utils.NewPrivateKey("372fcc2abc6bc5015103aae4763822e45c4cfe775d163f97a9ebdd77b0d12c0c")
	if err != nil {
		t.Fatal(err)
	}
	exPubSpend, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	keys, err := utils.NewFullKeyPairSpendPrivateKey(exPrivSpend)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exPrivView.Bytes(), keys.ViewKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exPubView.Bytes(), keys.ViewKeyPair().PublicKey().Bytes())
	assert.Equal(t, exPrivSpend.Bytes(), keys.SpendKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exPubSpend.Bytes(), keys.ViewOnlyKeyPair().SpendPublicKey().Bytes())
}
