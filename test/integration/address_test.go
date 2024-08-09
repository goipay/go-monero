package test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/chekist32/go-monero/utils"
	"github.com/chekist32/go-monero/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGenerateMoneroSubaddress(t *testing.T) {
	t.Parallel()

	privView, err := utils.NewPrivateKey("8aa763d1c8d9da4ca75cb6ca22a021b5cca376c1367be8d62bcc9cdf4b926009")
	if err != nil {
		t.Fatal(err)
	}
	pubSpend, err := utils.NewPublicKey("38e9908d33d034de0ba1281aa7afe3907b795cea14852b3d8fe276e8931cb130")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	containerReq := testcontainers.ContainerRequest{
		Image:        "chekist32/monero-wallet-rpc:0.18.3.3",
		ExposedPorts: []string{"38083/tcp"},
		Mounts:       testcontainers.Mounts(testcontainers.BindMount(fmt.Sprintf("%v/resources/wallet", os.Getenv("PWD")), testcontainers.ContainerMountTarget("/monero/wallet"))),
		Cmd:          []string{"--stagenet", "--daemon-address=http://node.monerodevs.org:38089", "--trusted-daemon", "--rpc-bind-port=38083", "--disable-rpc-login", "--wallet-dir=/monero/wallet"},

		WaitingFor: wait.ForLog(`Starting wallet RPC server`),
	}

	walletRpc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer walletRpc.Terminate(ctx)

	u, err := walletRpc.PortEndpoint(ctx, "38083/tcp", "http")
	if err != nil {
		t.Fatal(err)
	}

	w := wallet.New(wallet.Config{Address: u})
	if err := w.OpenWallet(&wallet.RequestOpenWallet{Filename: "test"}); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 20; i++ {
		t.Run(fmt.Sprintf("Test Generate Monero Subaddress #%v", i), func(t *testing.T) {
			major := rand.Uint32() % 100
			minor := rand.Uint32() % 100

			res, err := w.GetAccounts(&wallet.RequestGetAccounts{})
			if err != nil {
				t.Fatal(err)
			}
			for i := uint32(len(res.SubaddressAccounts)); i <= uint32(major); i++ {
				if _, err := w.CreateAccount(&wallet.RequestCreateAccount{}); err != nil {
					t.Fatal(err)
				}
			}

			exAddr, err := w.GetAddress(&wallet.RequestGetAddress{AccountIndex: uint64(major), AddressIndex: []uint64{uint64(minor)}})
			if err != nil {
				if err.Error() != "address index is out of bound" {
					t.Fatal(err)
				}
			}

			for err != nil {
				if _, err := w.CreateAddress(&wallet.RequestCreateAddress{AccountIndex: uint64(major)}); err != nil {
					t.Fatal(err)
				}

				exAddr, err = w.GetAddress(&wallet.RequestGetAddress{AccountIndex: uint64(major), AddressIndex: []uint64{uint64(minor)}})
				if err != nil {
					if err.Error() != "address index is out of bound" {
						t.Fatal(err)
					}
				}
			}

			addr, err := utils.GenerateSubaddress(privView, pubSpend, uint32(major), uint32(minor), utils.Stagenet)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, exAddr.Addresses[0].Address, addr.Address())
		})
	}

}
