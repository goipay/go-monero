package test

import (
	"encoding/hex"
	"testing"

	"github.com/chekist32/go-monero/utils"

	"github.com/stretchr/testify/assert"
)

func TestXMRToDecimalTest(t *testing.T) {
	assert.Equal(t, "0.034000200000", utils.XMRToDecimal(34000200000))
	assert.Equal(t, "15.000000000000", utils.XMRToDecimal(15e12))
}

func TestXMRToFloat64Test(t *testing.T) {
	assert.Equal(t, float64(0.02), utils.XMRToFloat64(20000000000))
	assert.Equal(t, float64(3.14), utils.XMRToFloat64(314e10))
}

func TestGetTxPublicKeyFromExtra(t *testing.T) {
	extra, err := hex.DecodeString("0166488b56658159e08f88f90afca5cb9ea3e3fbe460b0328b83a6ef7ca7a81835020901c26ecfa7aabbb41b")
	if err != nil {
		t.Fatal(err)
	}

	txPub, err := utils.GetTxPublicKeyFromExtra(extra)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := utils.NewPublicKey("66488b56658159e08f88f90afca5cb9ea3e3fbe460b0328b83a6ef7ca7a81835")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, txPub)
}

func TestOutputBelongsViewTag(t *testing.T) {
	txPub, err := utils.NewPublicKey("7302dd77bf4095baf868de43b7a32f4a36fe9d8b48ccfff537157a4a786fa364")
	if err != nil {
		t.Fatal(err)
	}

	viewKey, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}

	res, err := utils.OutputBelongsViewTag("1a", 1, txPub, viewKey)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, res)
}

func TestDecryptOutputViewTag(t *testing.T) {
	txPub, err := utils.NewPublicKey("7302dd77bf4095baf868de43b7a32f4a36fe9d8b48ccfff537157a4a786fa364")
	if err != nil {
		t.Fatal(err)
	}

	viewKey, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}

	res, am, err := utils.DecryptOutputViewTag("1a", 1, "5db33f80fd4990bc", txPub, viewKey)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, res)
	assert.Equal(t, utils.Float64ToXMR(0.55), am)
}

func TestOutputBelongsPublicKey(t *testing.T) {
	txPub, err := utils.NewPublicKey("7302dd77bf4095baf868de43b7a32f4a36fe9d8b48ccfff537157a4a786fa364")
	if err != nil {
		t.Fatal(err)
	}

	viewKey, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}

	outKey, err := utils.NewPublicKey("7e4f4427539b206740bed78b81b0dc10acb89aa1545880863f73264492ee0c16")
	if err != nil {
		t.Fatal(err)
	}

	spendKey, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	res, err := utils.OutputBelongsPublicSpendKey(spendKey, 1, outKey, txPub, viewKey)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, res)
}

func TestDecryptOutputPublicSpendKey(t *testing.T) {
	txPub, err := utils.NewPublicKey("7302dd77bf4095baf868de43b7a32f4a36fe9d8b48ccfff537157a4a786fa364")
	if err != nil {
		t.Fatal(err)
	}

	viewKey, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}

	outKey, err := utils.NewPublicKey("7e4f4427539b206740bed78b81b0dc10acb89aa1545880863f73264492ee0c16")
	if err != nil {
		t.Fatal(err)
	}

	spendKey, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	res, am, err := utils.DecryptOutputPublicSpendKey(spendKey, 1, outKey, "5db33f80fd4990bc", txPub, viewKey)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, res)
	assert.Equal(t, utils.Float64ToXMR(0.55), am)
}

func TestParseExtra(t *testing.T) {
	extra, err := hex.DecodeString("01c4865e47b9392d52e6f4957d0a6f6a9feda0ef0ac4807e1127bea56cb3ba583e020901a88dcbe6025e76ce")
	if err != nil {
		t.Fatal(err)
	}

	expectedTxPub, err := hex.DecodeString("c4865e47b9392d52e6f4957d0a6f6a9feda0ef0ac4807e1127bea56cb3ba583e")
	if err != nil {
		t.Fatal(err)
	}

	expectedPayId, err := hex.DecodeString("a88dcbe6025e76ce")
	if err != nil {
		t.Fatal(err)
	}

	txPub, payId := utils.ParseExtra(extra)

	assert.Equal(t, expectedTxPub, txPub.Bytes())
	assert.Equal(t, expectedPayId, payId)
}

func TestGenerateSubaddress(t *testing.T) {
	cases := []struct {
		vk           string
		sk           string
		nt           utils.NetworkType
		major, minor uint32
		expected     string
	}{
		{
			vk:       "ac413c16b815899b69393d72086fa86d31e8e352895606180c4c8fadd707450a",
			sk:       "c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b106",
			nt:       utils.Mainnet,
			major:    1,
			minor:    0,
			expected: "87BTvS4grAXSrwgzonu3N8Tm7N6W29UGAcd3GLumriVYiCJrUbsyPGWQoA92FZ6MgKWStiZjhS6o9Eeh6yinHH5NAgE9CUe",
		},
		{
			vk:       "ac413c16b815899b69393d72086fa86d31e8e352895606180c4c8fadd707450a",
			sk:       "c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b106",
			nt:       utils.Mainnet,
			major:    2,
			minor:    3,
			expected: "87aJx3x1cd56PS9JZdY4rzFSWcjvh274ERV1LYmFzzTwYFbfzWLRfgxTm8zBCPxPmoCpEnmHgDAn3dNAi1zRchNv8zdeQ1i",
		},
		{
			vk:       "8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009",
			sk:       "38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130",
			nt:       utils.Stagenet,
			major:    1,
			minor:    0,
			expected: "72c2F4L6XMu28Wf4e5yiVfKJcb4uDzvM9DxSAydF9o766RUiVqXawkhUcz7y59EBRrDafZB8DezLbLSrtb5xPL7s6PZ2zoj",
		},
		{
			vk:       "8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009",
			sk:       "38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130",
			nt:       utils.Stagenet,
			major:    3,
			minor:    5,
			expected: "74wdCFDsraBfreEwnfyyexK5d5ZkU48bK6Xd1UGjFTvNYes7gQJY47WUdA23hny1ynC2REEM9Rf1DGNuuwbDrsuAEHrwVmv",
		},
	}

	for _, v := range cases {
		t.Run("", func(t *testing.T) {
			viewKey, err := utils.NewPrivateKey(v.vk)
			if err != nil {
				t.Fatal(err)
			}

			spendKey, err := utils.NewPublicKey(v.sk)
			if err != nil {
				t.Fatal(err)
			}

			res, err := utils.GenerateSubaddress(viewKey, spendKey, v.major, v.minor, v.nt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, v.expected, res.Address())
		})
	}

}

func TestPaymentID(t *testing.T) {
	PaymentID256 := utils.NewPaymentID256()
	PaymentID64 := utils.NewPaymentID64()

	assert.Equal(t, len(PaymentID256), 32)
	assert.Equal(t, len(PaymentID64), 8)
}
