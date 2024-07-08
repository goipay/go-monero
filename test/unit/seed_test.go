package test

import (
	"strings"
	"testing"

	"github.com/chekist32/go-monero/utils"
	"github.com/stretchr/testify/assert"
)

func TestSeedMnemonic(t *testing.T) {
	spend, err := utils.NewPrivateKey("0cca07dc4e90fc738fffdb2561dddd7a94d0dc8977d0229303d7509a10c9d705")
	if err != nil {
		t.Fatal(err)
	}

	exFullKeys, err := utils.NewFullKeyPairSpendPrivateKey(spend)
	if err != nil {
		t.Fatal(err)
	}

	seed, err := utils.NewSeedMnemonic("wiggle drowning auburn aquarium attire meant impel phase soothe heron android mechanic inroads energy smog niece enforce syllabus exquisite lush bluntly rage siblings soda syllabus", utils.English)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exFullKeys.ViewKeyPair().PrivateKey().Bytes(), seed.FullKeyPair().ViewKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exFullKeys.ViewKeyPair().PublicKey().Bytes(), seed.FullKeyPair().ViewKeyPair().PublicKey().Bytes())
	assert.Equal(t, exFullKeys.SpendKeyPair().PrivateKey().Bytes(), seed.FullKeyPair().SpendKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exFullKeys.SpendKeyPair().PublicKey().Bytes(), seed.FullKeyPair().SpendKeyPair().PublicKey().Bytes())
	assert.Equal(t, utils.English, seed.MnemonicLanguage())
}

func TestSeed(t *testing.T) {
	seed, err := utils.NewSeed(utils.English)
	if err != nil {
		t.Fatal(err)
	}

	exSeed, err := utils.NewSeedMnemonic(strings.Join(seed.Mnemonic(), " "), utils.English)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exSeed.FullKeyPair().ViewKeyPair().PrivateKey().Bytes(), seed.FullKeyPair().ViewKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exSeed.FullKeyPair().ViewKeyPair().PublicKey().Bytes(), seed.FullKeyPair().ViewKeyPair().PublicKey().Bytes())
	assert.Equal(t, exSeed.FullKeyPair().SpendKeyPair().PrivateKey().Bytes(), seed.FullKeyPair().SpendKeyPair().PrivateKey().Bytes())
	assert.Equal(t, exSeed.FullKeyPair().SpendKeyPair().PublicKey().Bytes(), seed.FullKeyPair().SpendKeyPair().PublicKey().Bytes())
	assert.Equal(t, utils.English, seed.MnemonicLanguage())
}
