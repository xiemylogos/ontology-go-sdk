package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum"

)

func main() {
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s3.bnbchain.org:8545")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	resolverAddr := common.HexToAddress("0xe6CD76AfF83054fA10222A9BcD4D38cee884E7e6") // 替换为部署地址
	abiStr := `[
        {"name":"resolve","type":"function","inputs":[{"name":"name","type":"bytes"},{"name":"data","type":"bytes"}],"outputs":[{"name":"","type":"bytes"}],"stateMutability":"view"},
        {"name":"addrCallback","type":"function","inputs":[{"name":"node","type":"bytes32"},{"name":"response","type":"bytes"}],"outputs":[{"name":"","type":"address"}],"stateMutability":"pure"}
    ]`
	contractAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	name := "alice.ont.im"
	node := namehash(name)
	data := append(contractAbi.Methods["addr"].ID, node.Bytes()...)
	callData, err := contractAbi.Pack("resolve", []byte(name), data)
	if err != nil {
		log.Fatalf("Failed to pack callData: %v", err)
	}
	 // Use ethereum.CallMsg instead of types.Message
	 msg := ethereum.CallMsg{
        To:   &resolverAddr,
        Data: callData,
        Gas:  500000,
    }

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		urls := []string{"http://localhost:3000/data"}
		resp, err := sendOffchainRequest(urls[0], callData)
		if err != nil {
			log.Fatalf("Offchain request failed: %v", err)
		}

		callbackData, _ := contractAbi.Pack("addrCallback", node, resp)
		msg.Data = callbackData
		result, err = client.CallContract(context.Background(), msg, nil)
		if err != nil {
			log.Fatalf("addrCallback failed: %v", err)
		}

		addr, err := contractAbi.Unpack("addrCallback", result)
		if err != nil {
			log.Fatalf("Failed to unpack result: %v", err)
		}
		fmt.Printf("Resolved address for %s: %s\n", name, addr[0].(common.Address).Hex())
	} else {
		fmt.Printf("Unexpected direct result: %x\n", result)
	}
}

func namehash(name string) common.Hash {
	labels := strings.Split(name, ".")
	hash := common.Hash{}
	for i := len(labels) - 1; i >= 0; i-- {
		hash = crypto.Keccak256Hash(append(hash[:], crypto.Keccak256([]byte(labels[i]))...))
	}
	return hash
}

func sendOffchainRequest(url string, callData []byte) ([]byte, error) {
	reqBody, err := json.Marshal(map[string]string{"callData": "0x" + hex.EncodeToString(callData)})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimPrefix(string(body), "0x"))
}
