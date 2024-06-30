package test

import (
	"chekist32/go-monero/daemon"
	"chekist32/go-monero/utils"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func addresses(addrs []string) ([]utils.MoneroAddress, error) {
	res := make([]utils.MoneroAddress, 0, len(addrs))

	for _, v := range addrs {
		addr, err := utils.NewAddress(v)
		if err != nil {
			return nil, err
		}
		res = append(res, addr)
	}

	return res, nil
}

// 1619000 - 1620200 Height
// 91 Txs
func BenchmarkBlockchainScanningPerformance(b *testing.B) {
	viewKey, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		b.Fatal(err)
	}

	addrs, err := addresses([]string{
		"74qBoZ9TLyiFodY4XE9ij6Ab6pWzASN6eaE9Rx6UUK4Mhi8zxqLaXpFfoyAfEDvGkFKhirUVPzqVkjS6ccqL8aCuT3iiPba",
		"72c2F4L6XMu28Wf4e5yiVfKJcb4uDzvM9DxSAydF9o766RUiVqXawkhUcz7y59EBRrDafZB8DezLbLSrtb5xPL7s6PZ2zoj",
		"76ZmdKxcxqJD7xbXejEQGjMtKeA8Uv2kNB6B2t7tptdSQ8yCma4xH3WVjvaf7QVFGMDC9scnFBLaJ2LWWU7SWAN4M6Fmiei",
		"76hyxvqnTZqGMSKiuTkwLdLPgh3x5u1paZVSg4i2HvngA7CebpkEWFq3eRpC4GbuqWCVgLyamaJoALdmeK8NT2zGDgJibEy",
		"79E1kMDMzqjZSFykNxVdGNX4TwJB1wwCV2dYacGp8pSeHEfCCTctJjiR7ufGwjoYHDdMm2htzJjSaEJsSrejHD7g9xqz7e6",
		"74xhb5sXRsnDZv8RKFEv7LAMfUq5AmGEEB77SVvsUJf8bLvFMSEfc8YYyJHF6xNNnjAZQmgqZp76AjT8bD6qKkLZLeR42oi",
		"75yDVY9uQiv1uTL3ccbtBLjTQh1H9Z49V3Jcds9qo83GWrceHPY9SfDXKF4EtKENmmMUxvqBG2EfFZco536rB3Tg7zSHVVs",
		"765sVRqnUyULtxGnWamHG8U4ghXm74U2sjYHnjSV1WohTEn8XXowHRiUvs7kCF1sk8ChUAVH36K1yYJeqbEReH2uPJej7GK",
		"53zEYzu2hi3e97tdMTqTvSRAfFYXwxA7LBJEHLWvFnm699WgcsE8CJujENwNAQotKyY2u94vpbGEZTiwahuMcMfX3x6NFwY",
		"78hRedVbk2N3Mg2DpMMUoCbynA1uZJzAr7R7rnCtBo4Q1FtnDePx7NPAcCGPXVEBTp96AjRnR9uchhan49fbBAnuLTU11cw",
	})
	if err != nil {
		b.Fatal(err)
	}

	data, err := os.ReadFile("txs.json")
	if err != nil {
		log.Fatal(err)
	}

	var txs []daemon.MoneroTx1
	if err := json.Unmarshal(data, &txs); err != nil {
		log.Fatal(err)
	}

	b.ResetTimer()
	for _, v := range txs {
		for i := range v.TxInfo.Vout {
			out := v.TxInfo.Vout[i]
			am := v.TxInfo.RctSignatures.EcdhInfo[i].Amount
			txPub, err := utils.GetTxPublicKeyFromExtra(v.TxInfo.Extra)
			if err != nil {
				b.Fatal(err)
			}

			if len(out.Target.TaggedKey.ViewTag) == 2 {
				res, err := utils.OutputBelongsViewTag(out.Target.TaggedKey.ViewTag, uint32(i), txPub, viewKey)
				if err != nil {
					b.Fatal(err)
				}
				if !res {
					continue
				}
			}

			if out.Target.TaggedKey.Key == "" {
				continue
			}

			outKey, err := utils.NewPublicKey(out.Target.TaggedKey.Key)
			if err != nil {
				b.Fatal(err)
			}

			for _, v := range addrs {
				_, _, err := utils.DecryptOutputPublicSpendKey(v.PublicSpendKey(), uint32(i), outKey, am, txPub, viewKey)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	}
}
