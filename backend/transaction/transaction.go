package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	v2 "github.com/2022AA/bytes-linked/backend/pkg/logging/v2"
	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	"github.com/FISCO-BCOS/go-sdk/examples"
	"github.com/FISCO-BCOS/go-sdk/precompiled/cns"
)

const KVTableTestABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"string\"}],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"},{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"item_price\",\"type\":\"int256\"},{\"name\":\"item_name\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"count\",\"type\":\"int256\"}],\"name\":\"SetResult\",\"type\":\"event\"}]"
const KTTableName = "kt_test"
const KTVersion = "1"
const KTTABLE_ADDRESS = "xxxxxxxxxxx"

func GenConfig(privateKey string) (*conf.Config, error) {

	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("decode hex failed of %v", err)
	}
	config := &conf.Config{IsHTTP: true, ChainID: 1, IsSMCrypto: false, GroupID: 1,
		PrivateKey: privateKeyBytes, NodeURL: "http://127.0.0.1:8545"}
	return config, nil
}

func Transaction(fileId string, data string, privateKey string) (string, error) {

	config, err := GenConfig(privateKey)
	if err != nil {
		return "", err
	}
	cli, err := client.Dial(config)
	if err != nil {
		return "", err
	}
	v2.Infoln("-------------------starting deploy contract-----------------------")

	newService, err := cns.NewCnsService(cli)
	if err != nil {
		return "", err
	}

	addr, err := newService.GetAddressByContractNameAndVersion(KTTableName, KTVersion)
	if err != nil {
		return "", err
	}

	instance, err := examples.NewKVTableTest(addr, cli)
	if err != nil {
		return "", err
	}

	v2.Infoln("-------------------starting invoke Set to insert info-----------------------")

	kvtabletestSession := &examples.KVTableTestSession{
		Contract: instance, CallOpts: *cli.GetCallOpts(),
		TransactOpts: *cli.GetTransactOpts(),
	}
	tx, receipt, err := kvtabletestSession.Set(fileId, big.NewInt(1), data) // call set API
	if err != nil {
		return "", err
	}

	transactionId := tx.Hash().Hex()[2:34] + receipt.BlockHash[2:34] + receipt.BlockNumber[2:]

	v2.Infof("tx sent: %s, receipt: %s, block_hash: %s, block_number: %s", tx.Hash().Hex(), receipt.Output, receipt.BlockHash, receipt.BlockNumber)
	v2.Infoln(transactionId)
	setedLines, err := parseOutput(examples.KVTableTestABI, "set", receipt)
	if err != nil {
		v2.Errorf("error when transfer string to int: %v\n", err)
	}
	v2.Infof("seted lines: %v", setedLines.Int64())
	return transactionId, nil
}

func parseOutput(abiStr, name string, receipt *types.Receipt) (*big.Int, error) {
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		fmt.Printf("parse ABI failed, err: %v", err)
	}
	var ret *big.Int
	b, err := hex.DecodeString(receipt.Output[2:])
	if err != nil {
		return nil, fmt.Errorf("decode receipt.Output[2:] failed, err: %v", err)
	}
	err = parsed.Unpack(&ret, name, b)
	if err != nil {
		return nil, fmt.Errorf("unpack %v failed, err: %v", name, err)
	}
	return ret, nil
}
