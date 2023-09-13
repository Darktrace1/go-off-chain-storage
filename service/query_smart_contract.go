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
	ethEndpoint := "http://bum0448.iptime.org:8545"
	client, _ := ethclient.Dial(ethEndpoint)
	defer client.Close()

	contractAddress := common.HexToAddress("0x9F5844648746c6ae351F0E586fE435bA1E193199")
	contractABI := U.ReadFile("../ABI.json")

	parsedABI, err := abi.JSON(strings.NewReader(string(contractABI)))
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
