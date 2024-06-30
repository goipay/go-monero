package test

import (
	"chekist32/go-monero/utils"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	addr  string
	dec   string
	view  string
	spend string
	payId string
	nt    utils.NetworkType
	at    utils.AddressType
}{
	// Mainnet
	{
		addr:  "48ukkZtBSBRL8iva7k3p2sBVMLWTfNwsTbW1aVh5M84g21muDCssvCHTpoZCaSc6rq8M9QLZ3sQMrMn1bq2RD2anGnyHhtq",
		dec:   "12c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b1060ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8bfa84423e",
		view:  "0ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b",
		spend: "c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b106",
		payId: "",
		nt:    utils.Mainnet,
		at:    utils.Primary,
	},
	{
		addr:  "84nvgV2eTnG1vAKbg87MnbfjWrSY3eH3s2eykmggk549C8zdNk4PPD7iv7BPfPsnoH9NjXaRhjC19FY6PBmXZUtoG5SEiY7",
		dec:   "2a3dba53246e6981057ad2a9eff6d164e791cdef4578caee09e4b6ea03af774e4296b37bede10496fa98e0a7c58c49403211b225643621e456e7d2d4d8c13b6885a192d0a8",
		view:  "96b37bede10496fa98e0a7c58c49403211b225643621e456e7d2d4d8c13b6885",
		spend: "3dba53246e6981057ad2a9eff6d164e791cdef4578caee09e4b6ea03af774e42",
		payId: "",
		nt:    utils.Mainnet,
		at:    utils.Sub,
	},
	{
		addr:  "4JcRmNhg3SwL8iva7k3p2sBVMLWTfNwsTbW1aVh5M84g21muDCssvCHTpoZCaSc6rq8M9QLZ3sQMrMn1bq2RD2anQMXGeH3St98D3GPzmn",
		dec:   "13c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b1060ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b9f9739432368cb6ab57ad0c9",
		view:  "0ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b",
		spend: "c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b106",
		payId: "9f9739432368cb6a",
		nt:    utils.Mainnet,
		at:    utils.Integrated,
	},

	// Stagenet
	{
		addr:  "53zEYzu2hi3e97tdMTqTvSRAfFYXwxA7LBJEHLWvFnm699WgcsE8CJujENwNAQotKyY2u94vpbGEZTiwahuMcMfX3x6NFwY",
		dec:   "1838e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a25777ccb",
		view:  "b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a",
		spend: "38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130",
		payId: "",
		nt:    utils.Stagenet,
		at:    utils.Primary,
	},
	{
		addr:  "74xhb5sXRsnDZv8RKFEv7LAMfUq5AmGEEB77SVvsUJf8bLvFMSEfc8YYyJHF6xNNnjAZQmgqZp76AjT8bD6qKkLZLeR42oi",
		dec:   "2447a69d7aa0d0b14b22e2ff185b2e8b37effd771fee0c8b3c6a7ff9d53910ffcd536bab26fc7101bf23da6b6a39daa83925f5393df58a1bfdcb8edcc388726aae2014d691",
		view:  "536bab26fc7101bf23da6b6a39daa83925f5393df58a1bfdcb8edcc388726aae",
		spend: "47a69d7aa0d0b14b22e2ff185b2e8b37effd771fee0c8b3c6a7ff9d53910ffcd",
		payId: "",
		nt:    utils.Stagenet,
		at:    utils.Sub,
	},
	{
		addr:  "5DguZoiXJyZe97tdMTqTvSRAfFYXwxA7LBJEHLWvFnm699WgcsE8CJujENwNAQotKyY2u94vpbGEZTiwahuMcMfX5MsmWgk84zrS4MPMnW",
		dec:   "1938e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a10f5ebd54675efde19e372ef",
		view:  "b4cdbf52851002fc7b098b99536df8b9885aa6cb8db24e9fc46103674dc9421a",
		spend: "38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130",
		payId: "10f5ebd54675efde",
		nt:    utils.Stagenet,
		at:    utils.Integrated,
	},

	// Testnet
	{
		addr:  "9zvkxwHbuHxX8B82zA8G9yBh6oKzbXS8viKexKeBCVBwNeP246aVAKSiC1DyVoETYZ11qDdmibSShX88HWGevRbp3G6hKyK",
		dec:   "35ccc9377cde8377b4190b13b1384c2c3feb697bc98a783ff70bbdf789d623e2816763ef0f8d3b41f641db860acbe5360015f00f6d0f05b6b417c0ad4708277b1408cfe1ba",
		view:  "6763ef0f8d3b41f641db860acbe5360015f00f6d0f05b6b417c0ad4708277b14",
		spend: "ccc9377cde8377b4190b13b1384c2c3feb697bc98a783ff70bbdf789d623e281",
		payId: "",
		nt:    utils.Testnet,
		at:    utils.Primary,
	},
	{
		addr:  "Be4mtTzNR3gGe7S9foMPdiLnv7jP2cfR3FB4YiNjNQFh8Q6WGortUdtXgwumP6xRu8MxdozhRjXMf4gCwwwE7NtRQ5jMkJd",
		dec:   "3f9b5365b83a95d35d8127e5af7b9fa976539907f255ad9254bcd54705a7e6582c3b1aba22c7292fb779d907deb9fbd77d4e96286f22053615fa273513938b9acc95b06f0a",
		view:  "3b1aba22c7292fb779d907deb9fbd77d4e96286f22053615fa273513938b9acc",
		spend: "9b5365b83a95d35d8127e5af7b9fa976539907f255ad9254bcd54705a7e6582c",
		payId: "",
		nt:    utils.Testnet,
		at:    utils.Sub,
	},
	{
		addr:  "AAdRyk76WZUX8B82zA8G9yBh6oKzbXS8viKexKeBCVBwNeP246aVAKSiC1DyVoETYZ11qDdmibSShX88HWGevRbp4MEeHy9Xttg2tozpVS",
		dec:   "36ccc9377cde8377b4190b13b1384c2c3feb697bc98a783ff70bbdf789d623e2816763ef0f8d3b41f641db860acbe5360015f00f6d0f05b6b417c0ad4708277b14058bca5e06c79110c7fd10f5",
		view:  "6763ef0f8d3b41f641db860acbe5360015f00f6d0f05b6b417c0ad4708277b14",
		spend: "ccc9377cde8377b4190b13b1384c2c3feb697bc98a783ff70bbdf789d623e281",
		payId: "058bca5e06c79110",
		nt:    utils.Testnet,
		at:    utils.Integrated,
	},
}

func TestDecodeMoneroAddress(t *testing.T) {
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			expected, err := hex.DecodeString(c.dec)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := utils.DecodeMoneroAddress(c.addr)
			assert.NoError(t, err)

			assert.Equal(t, expected, actual)
		})
	}
}

func TestEncodeMoneroAddress(t *testing.T) {
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			addr, err := hex.DecodeString(c.dec)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := utils.EncodeMoneroAddress(addr)
			assert.NoError(t, err)

			assert.Equal(t, []byte(c.addr), actual)
		})
	}
}

func TestAddress(t *testing.T) {
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			expectedView, err := hex.DecodeString(c.view)
			if err != nil {
				t.Fatal(err)
			}
			expectedSpend, err := hex.DecodeString(c.spend)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := utils.NewAddress(c.addr)
			assert.NoError(t, err)

			assert.Equal(t, c.nt, actual.NetworkType())
			assert.Equal(t, c.at, actual.AddressType())
			assert.Equal(t, expectedSpend, actual.PublicSpendKey().Bytes())
			assert.Equal(t, expectedView, actual.PublicViewKey().Bytes())

			if actual.AddressType() == utils.Integrated {
				act, ok := actual.(*utils.IntegratedAddress)
				assert.True(t, ok)

				expectedpayId, err := hex.DecodeString(c.payId)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, expectedpayId, act.PaymentId())
			}
		})
	}

	viewKey, err := utils.NewPublicKey("0ef3c9e1146ed2a05f0eb4b25e41662bed41fa246251257c363a8ba95750cb8b")
	if err != nil {
		t.Fatal(err)
	}
	spendKey, err := utils.NewPublicKey("c04ac8adc844e07263bf9a4dd337883eb55db89743c9aece4357381ae6c0b106")
	if err != nil {
		t.Fatal(err)
	}

	addr, err := utils.NewAddress("48ukkZtBSBRL8iva7k3p2sBVMLWTfNwsTbW1aVh5M84g21muDCssvCHTpoZCaSc6rq8M9QLZ3sQMrMn1bq2RD2anGnyHhtq")
	assert.NoError(t, err)

	assert.Equal(t, viewKey.Bytes(), addr.PublicViewKey().Bytes())
	assert.Equal(t, spendKey.Bytes(), addr.PublicSpendKey().Bytes())

}
