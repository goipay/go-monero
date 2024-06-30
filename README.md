Monero Go Library
====================

<p align="center">
<img src="./media/img/monero_gopher.png" alt="Monero Gopher" width="200" />
</p>

A client implementation for the Monero wallet and daemon RPC written in go.
This package is inspired by https://github.com/omani/go-monero-rpc-client.

## Wallet RPC Client

[![GoDoc](https://godoc.org/github.com/chekist32/go-monero/wallet?status.svg)](https://godoc.org/github.com/chekist32/go-monero/wallet)

### Monero RPC Version
The ```go-monero/wallet``` package is the RPC client for version `v1.3` of the [Monero Wallet RPC](https://www.getmonero.org/resources/developer-guides/wallet-rpc.html).

### Installation

```sh
go get -u github.com/chekist32/go-monero
```

#### Spawn the monero-wallet-rpc daemon (without rpc login):

```sh
./monero-wallet-rpc --wallet-file /home/$user/stagenetwallet/stagenetwallet --daemon-address pool.cloudissh.com:38081 --stagenet --rpc-bind-port 6061 --password 'mystagenetwalletpassword' --disable-rpc-login
```
You can use our remote node for the stagenet running at pool.cloudissh.com port `38081`.

#### Go code:

```Go
package main

import (
  "encoding/json"
  "fmt"
  "log"

  "github.com/chekist32/go-monero/wallet"
)

func checkerr(err error) {
  if err != nil {
    log.Panic(err)
  }
}

func main() {
  // Start a wallet client instance
  client := wallet.New(wallet.Config{
    Address: "http://127.0.0.1:6061/json_rpc",
  })

  // check wallet balance
  resp, err := client.GetBalance(&wallet.RequestGetBalance{AccountIndex: 0})
  checkerr(err)
  res, _ := json.MarshalIndent(resp, "", "\t")
  fmt.Print(string(res))

  // get incoming transfers
  resp1, err := client.GetTransfers(&wallet.RequestGetTransfers{
    AccountIndex: 0,
    In:           true,
  })
  checkerr(err)
  for _, in := range resp1.In {
    res, _ := json.MarshalIndent(in, "", "\t")
    fmt.Print(string(res))
  }
}
```

### Spawn the monero-wallet-rpc daemon (with rpc login):

```sh
./monero-wallet-rpc --wallet-file /home/$user/stagenetwallet/stagenetwallet --daemon-address pool.cloudissh.com:38081 --stagenet --rpc-bind-port 6061 --password 'mystagenetwalletpassword' --rpc-login test:testpass
```

#### Go code:

```Go
package main

import (
  "encoding/json"
  "fmt"
  "log"

  "github.com/chekist32/go-monero/wallet"
)

func checkerr(err error) {
  if err != nil {
    log.Panic(err)
  }
}

func main() {
  t := httpdigest.New("test", "testpass")

  // Start a wallet client instance
  client := wallet.New(wallet.Config{
    Address: "http://127.0.0.1:6061/json_rpc",
    Transport: t,
  })

  // check wallet balance
  resp, err := client.GetBalance(&wallet.RequestGetBalance{AccountIndex: 0})
  checkerr(err)
  res, _ := json.MarshalIndent(resp, "", "\t")
  fmt.Print(string(res))

  // get incoming transfers
  resp1, err := client.GetTransfers(&wallet.RequestGetTransfers{
    AccountIndex: 0,
    In:           true,
  })
  checkerr(err)
  for _, in := range resp1.In {
    res, _ := json.MarshalIndent(in, "", "\t")
    fmt.Print(string(res))
  }
}
```


## Daemon RPC Client

[![GoDoc](https://godoc.org/github.com/chekist32/go-monero/wallet?status.svg)](https://godoc.org/github.com/chekist32/go-monero/daemon)

Here is a [List of implemented methods.](https://github.com/chekist32/go-monero/issues/5)

**Go code:**
```Go
package main

import (
	"github.com/chekist32/go-monero/daemon"
	"fmt"
	"log"
	"net/url"
)

func main() {
	u, err := url.Parse("http://xmr-node.cakewallet.com:18081")
	if err != nil {
		log.Fatal(err)
	}

	d := daemon.NewDaemonRpcClient(daemon.NewRpcConnection(u, "", ""))

	res, err := d.GetCurrentHeight()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current height: %v\n", res.Height)
}
```

## Monero Utils

[![GoDoc](https://godoc.org/github.com/chekist32/go-monero/wallet?status.svg)](https://godoc.org/github.com/chekist32/go-monero/utils)

This package contains helper methods that can be used for different cases, such as subaddress generation and tx output decryption.

**Go code:**
```Go
package main

import (
	"github.com/chekist32/go-monero/utils"
	"fmt"
	"log"
)

func main() {
	txPub, err := utils.NewPublicKey("7302dd77bf4095baf868de43b7a32f4a36fe9d8b48ccfff537157a4a786fa364")
	if err != nil {
		log.Fatal(err)
	}

	viewKey, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		log.Fatal(err)
	}

	outKey, err := utils.NewPublicKey("7e4f4427539b206740bed78b81b0dc10acb89aa1545880863f73264492ee0c16")
	if err != nil {
		log.Fatal(err)
	}

	spendKey, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		log.Fatal(err)
	}

	res, am, err := utils.DecryptOutputPublicSpendKey(spendKey, 1, outKey, "5db33f80fd4990bc", txPub, viewKey)
	if err != nil {
		log.Fatal(err)
	}

	if res {
		fmt.Printf("Received: %v\n", utils.XMRToFloat64(am))
	} else {
		fmt.Println("The output doesn't belong to the public spend key")
	}
}
```

# Contributing
- Before the actual PR, please create an issue where you can describe the improvements you want to add.


# LICENSE
MIT License
