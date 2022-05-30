package ipfs

import (
	"context"
	"fmt"
	"os"
	"testing"

	shell "github.com/ipfs/go-ipfs-api"
)

func TestIpfsAdd(t *testing.T) {
	ipfsCli = shell.NewShell("xxxx:5001")
	file, err := os.Open("./deployment-cheatsheet.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fileHash, err := ipfsCli.Add(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("fileHash:-----------------", fileHash)
}

func TestIpfsGet(t *testing.T) {
	ipfsCli = shell.NewShell("xxxx:5001")
	fileHash := "xxxxxxxxxxxxxxxxxxxxx"
	err := ipfsCli.Get(fileHash, "./test.png")
	if err != nil {
		fmt.Println(err)
	}
}

func TestIpfsAddDir(t *testing.T) {
	ipfsCli = shell.NewShell("xxxx:5001")
	file, err := os.Open("./deployment-cheatsheet.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	err = ipfsCli.FilesWrite(context.Background(), "/test/test.png", file, shell.FilesWrite.Create(true))
	if err != nil {
		fmt.Println(err)
	}
}
