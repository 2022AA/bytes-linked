package transaction

import (
	"fmt"
	"math/big"
	"testing"

	v2 "github.com/2022AA/bytes-linked/backend/pkg/logging/v2"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/examples"
	"github.com/FISCO-BCOS/go-sdk/precompiled/cns"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	privateKey := ""
	cfg, err := GenConfig(privateKey)
	require.NoError(t, err)
	t.Logf("%+v", cfg)
}

func TestInsert(t *testing.T) {
	v2.SetLevel(v2.DebugLevel)
	privateKey := ""
	config, err := GenConfig(privateKey)
	require.NoError(t, err)
	t.Logf("%+v", config)
	client, err := client.Dial(config)
	require.NoError(t, err)

	// deploy contract
	fmt.Println("-------------------starting deploy contract-----------------------")
	address, tx, instance, err := examples.DeployKVTableTest(
		client.GetTransactOpts(), client)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println("contract address: ", address.Hex()) // the address should be saved
	fmt.Println("transaction hash: ", tx.Hash().Hex())
	_ = instance

	// invoke Set to insert info
	fmt.Println("\n-------------------starting invoke Set to insert info-----------------------")
	kvtabletestSession := &examples.KVTableTestSession{Contract: instance, CallOpts: *client.GetCallOpts(), TransactOpts: *client.GetTransactOpts()}
	id := "100010001001"
	item_name := "2000"
	item_price := big.NewInt(6000)
	tx, receipt, err := kvtabletestSession.Set(id, item_price, item_name) // call set API
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
	setedLines, err := parseOutput(examples.KVTableTestABI, "set", receipt)
	if err != nil {
		logrus.Fatalf("error when transfer string to int: %v\n", err)
	}
	fmt.Printf("seted lines: %v\n", setedLines.Int64())

	// invoke Get to query info
	fmt.Println("\n-------------------starting invoke Get to query info-----------------------")
	bool, item_price, item_name, err := kvtabletestSession.Get(id) // call get API
	if err != nil {
		t.Fatal(err)
	}
	if !bool {
		t.Fatalf("id：%v is not found \n", id)
	}
	fmt.Printf("id: %v, item_price: %v, item_name: %v \n", id, item_price, item_name)
}

func TestRegisterAddress(t *testing.T) {
	address := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	privateKey := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	config, err := GenConfig(string(privateKey))
	require.NoError(t, err)
	t.Logf("%+v", config)
	client, err := client.Dial(config)
	require.NoError(t, err)
	newService, err := cns.NewCnsService(client)
	require.NoError(t, err)

	result, err := newService.RegisterCns(KTTableName, KTVersion, common.HexToAddress(address), KVTableTestABI)
	require.NoError(t, err)

	if result != 1 {
		t.Fatalf("TestRegisterCns failed, the result %v is inconsistent with \"%v\"", result, 1)
	}
	t.Logf("TestRegisterCns result: %v", result)

}

func TestAddress(t *testing.T) {
	ok := common.IsHexAddress(KTTABLE_ADDRESS)
	t.Log(ok)

	privateKey := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	config, err := GenConfig(string(privateKey))
	require.NoError(t, err)
	t.Logf("%+v", config)
	client, err := client.Dial(config)
	require.NoError(t, err)
	newService, err := cns.NewCnsService(client)
	require.NoError(t, err)

	addr, err := newService.GetAddressByContractNameAndVersion(KTTableName, KTVersion)

	require.NoError(t, err)
	if addr.Hex() != KTTABLE_ADDRESS {
		t.Fatalf("111111111111")
	}
}
