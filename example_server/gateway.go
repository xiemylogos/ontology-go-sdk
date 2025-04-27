package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ontio/ontology-go-sdk"
	com "github.com/ontio/ontology/common"
)

type Request struct {
	CallData string `json:"callData"`
}

func main() {
	ontSdk := ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress("http://testnet1.ont.io:40336")

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		dataHandler(w, r, ontSdk)
	})

	fmt.Println("Gateway running on :6008")
	log.Fatal(http.ListenAndServe(":6008", nil))
}

func dataHandler(w http.ResponseWriter, r *http.Request, ontSdk *ontology_go_sdk.OntologySdk) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	callDataBytes, err := hex.DecodeString(strings.TrimPrefix(req.CallData, "0x"))
	if err != nil {
		http.Error(w, "Invalid callData", http.StatusBadRequest)
		return
	}
	name := extractNameFromCallData(callDataBytes)
	if !strings.HasSuffix(name, ".ont.im") {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	ensName := strings.TrimSuffix(name, ".ont.im")
	fmt.Println("ensName:", ensName)
	addr, err := getAddressFromOntology(ontSdk, ensName)
	if err != nil {
		log.Printf("Failed to query %s: %v", ensName, err)
		addr = "0x00"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(addr))
}

func extractNameFromCallData(callData []byte) string {
	if len(callData) > 68 {
		length := int(callData[35])
		return string(callData[36 : 36+length])
	}
	return ""
}

func getAddressFromOntology(ontSdk *ontology_go_sdk.OntologySdk, ensName string) (string, error) {
	ensContractAddr, err := com.AddressFromHexString("2042c3b54c681e4ba728a022d01ab950612445ac") //ont ens contract
	if err != nil {
		return "",err
	}
	res, err := ontSdk.WasmVM.PreExecInvokeWasmVMContract(ensContractAddr, "resolve", []interface{}{ensName})
	if err != nil {
		return "",err
	}
	bs, err := res.Result.ToByteArray()
	if err != nil {
		return "",err
	}
	addr, err := com.AddressParseFromBytes(bs)
	if err != nil {
		return "",err
	}
	return addr.ToBase58(),nil
}
