package ipfs

import (
	"errors"
	"fmt"

	shell "github.com/ipfs/go-ipfs-api"
)

const endpoint = "xxxxx"

var ipfsCli *shell.Shell

// GetIpfsConnection : 获取ipfs连接
func GetIpfsConnection(ip, port string) (*shell.Shell, error) {
	if ip == "" {
		return nil, errors.New("")
	}
	if port == "" {
		return nil, errors.New("")
	}

	if ipfsCli != nil {
		return ipfsCli, nil
	}
	ipfsCli = shell.NewShell(fmt.Sprintf("%s:%s", ip, port))
	return ipfsCli, nil
}

func init() {
	ipfsCli = shell.NewShell(endpoint)
	ver, comm, err := ipfsCli.Version()
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to ipfs:", endpoint, ver, comm)
}

func Client() *shell.Shell {
	return ipfsCli
}
