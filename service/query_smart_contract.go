package service

import (
	"context"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
)

func init() {
	err := godotenv.Load(".env")
	U.CheckErr(err)
}

func QuerySmartContract(key string) string {
	ethEndpoint := os.Getenv("ETH_ENDPOINT")
	client, _ := ethclient.Dial(ethEndpoint)
	defer client.Close()

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	contractABI, err := os.ReadFile(os.Getenv("ABI_PATHS"))
	U.CheckErr(err)

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
