package test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/chekist32/go-monero/daemon"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func getUrlFromEnv(name string) (*url.URL, error) {
	urlStr := strings.TrimSpace(os.Getenv(name))
	if urlStr == "" {
		return nil, errors.New(name + " env can't be empty")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func initSimpleDaemonRpcClient(t *testing.T) daemon.IDaemonRpcClient {
	return initDaemonRpcClientWithCreds(t, "", "")
}

func initDaemonRpcClientWithCreds(t *testing.T, username, password string) daemon.IDaemonRpcClient {
	u, err := getUrlFromEnv("MONERO_DAEMON_RPC_ADDRESS")
	if err != nil {
		t.Fatal(err)
	}

	return daemon.NewDaemonRpcClient(daemon.NewRpcConnection(u, username, password))
}

func initDaemonRpcClientWithCredsAndCustomUrl(t *testing.T, urlStr, username, password string) daemon.IDaemonRpcClient {
	if urlStr == "" {
		t.Fatal("Url can't be empty")
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	return daemon.NewDaemonRpcClient(daemon.NewRpcConnection(u, username, password))
}

func TestGetTransactionPool(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetTransactionPool()
	if err != nil {
		t.Error(err)
	}
}
func TestGetBlockByHash(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockByHash(false, "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428")
	if err != nil {
		t.Error(err)
	}
}

func TestGetBlockByheight(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockByHeight(false, 2751506)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBlockCount(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockCount()
	if err != nil {
		t.Error(err)
	}
}

func TestGetBlockHeaderByHash(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockHeaderByHash(false, "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428")
	if err != nil {
		t.Error(err)
	}
}

func TestGetBlockHeaderByHe(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockHeaderByHeight(false, 2751506)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBlockHeadersRange(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockHeadersRange(false, 2751506, 2751507)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBlockTemplate(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockTemplate("48ukkZtBSBRL8iva7k3p2sBVMLWTfNwsTbW1aVh5M84g21muDCssvCHTpoZCaSc6rq8M9QLZ3sQMrMn1bq2RD2anGnyHhtq", 123)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitBlock(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.SubmitBlock([]string{"0707e6bdfedc053771512f1bc27c62731ae9e8f2443db64ce742f4e57f5cf8d393de28551e441a0000000002fb830a01ffbf830a018cfe88bee283060274c0aae2ef5730e680308d9c00b6da59187ad0352efe3c71d36eeeb28782f29f2501bd56b952c3ddc3e350c2631d3a5086cac172c56893831228b17de296ff4669de020200000000"})
	rpcErr, ok := err.(*daemon.MoneroRpcError)
	assert.True(t, ok)

	assert.Equal(t, rpcErr, &daemon.MoneroRpcError{Code: -7, Message: "Block not accepted"})
}

func TestGetCurrentHeight(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetCurrentHeight()
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastBlockHeader(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetLastBlockHeader(false)
	if err != nil {
		t.Error(err)
	}
}

func TestOnGetBlockHash(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.OnGetBlockHash(2751506)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTransactions(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	cases := []struct {
		p1 []string
		p2 bool
		p3 bool
		p4 bool
	}{
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, true, false, false},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, true, true, false},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, false, false, true},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, true, true, true},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, false, false, false},
	}

	wait := &sync.WaitGroup{}
	wait.Add(len(cases))
	errorChan := make(chan error, len(cases))

	for _, v := range cases {
		go func(params struct {
			p1 []string
			p2 bool
			p3 bool
			p4 bool
		}) {
			_, err := c.GetTransactions(params.p1, params.p2, params.p3, params.p4)
			if err != nil {
				errorChan <- err
			}
		}(v)
	}
	wait.Done()
	close(errorChan)

	for v := range errorChan {
		t.Error(v)
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)
	_, err := c.GetVersion()
	if err != nil {
		t.Error(err)
	}
}

func TestGetFeeEstimate(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)
	_, err := c.GetFeeEstimate()
	if err != nil {
		t.Error(err)
	}
}

func TestGetInfo(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)
	_, err := c.GetInfo()
	if err != nil {
		t.Error(err)
	}
}

func TestDigestAuth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	containerReq := testcontainers.ContainerRequest{
		Image:        "sethsimmons/simple-monerod:v0.18.3.3",
		ExposedPorts: []string{"18081/tcp"},
		Cmd:          []string{"--rpc-restricted-bind-ip=0.0.0.0", "--rpc-bind-ip=0.0.0.0", "--confirm-external-bind", "--rpc-login=user:pass", "--offline"},
		WaitingFor:   wait.ForLog(`Use "help <command>" to see a command's documentation.`),
	}

	monerodC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	defer monerodC.Terminate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	u, err := monerodC.PortEndpoint(ctx, "18081/tcp", "http")
	if err != nil {
		t.Fatal(err)
	}

	client := initDaemonRpcClientWithCredsAndCustomUrl(t, u, "user", "pass")

	_, err = client.GetBlockCount()
	if err != nil {
		t.Error(err)
	}

}
