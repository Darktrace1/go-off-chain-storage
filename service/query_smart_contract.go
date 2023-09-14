package service

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
)

func QuerySmartContract(key string) string {
	// ethEndpoint := "http://172.15.200.2:8545"
	ethEndpoint := "http://bum0448.iptime.org:8545"
	client, _ := ethclient.Dial(ethEndpoint)
	defer client.Close()

	contractAddress := common.HexToAddress("0x66c59762390A016F17453447a5FE1c14d0A91B5B")

	contractABI := `[ { "inputs": [ { "internalType": "string", "name": "key", "type": "string" }, { "internalType": "string", "name": "value", "type": "string" } ], "name": "setConfigValue", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [], "stateMutability": "nonpayable", "type": "constructor" }, { "inputs": [ { "internalType": "string", "name": "key", "type": "string" } ], "name": "getConfigValue", "outputs": [ { "internalType": "string", "name": "", "type": "string" } ], "stateMutability": "view", "type": "function" } ]`

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	U.CheckErr(err)

	callData, err := parsedABI.Pack("getConfigValue", key)
	U.CheckErr(err)

	res, err := client.CallContract(context.TODO(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	U.CheckErr(err)

	var Cluster_EndPoint string
	err = parsedABI.UnpackIntoInterface(&Cluster_EndPoint, "getConfigValue", res)
	U.CheckErr(err)

	return Cluster_EndPoint
}
