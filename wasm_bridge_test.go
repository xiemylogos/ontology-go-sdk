package ontology_go_sdk

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

// deploy bridge wasm contract
func TestDeployBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	wasmfile := "./wasm_file/bridge.wasm"
	//wasmfile := "./bridge.wasm"
	code, err := ioutil.ReadFile(wasmfile)
	assert.Nil(t, err)
	codeHash := common.ToHexString(code)
	gasprice := uint64(2500)
	//invokegaslimit := uint64(200000)
	t.Logf("signer addr:%s", testDefAcc.Address.ToBase58())
	deploygaslimit := uint64(200000000)
	txHash, err := testOntSdk.WasmVM.DeployWasmVMSmartContract(
		gasprice,
		deploygaslimit,
		testDefAcc,
		codeHash,
		"tokenbridge",
		"1.0",
		"author",
		"email",
		"desc",
	)
	assert.Nil(t, err)
	timeout := 20 * time.Second
	_, err = testOntSdk.WaitForGenerateBlock(timeout)
	assert.Nil(t, err)
	t.Logf("deploy wasm contract txhash is %s", txHash.ToHexString())
	contractAddr, err := utils.GetContractAddress(codeHash)
	assert.Nil(t, err)
	t.Logf("the contractAddr is:%s", contractAddr.ToHexString())
}

// init bridge wasm contract
func TestInitBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("816c02c3d1b30fb10d1bb493d8e4d2af5731449c") //uniswapv2gatway
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "init", []interface{}{testDefAcc.Address})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestBridgeRegisterOep4Erc20TokenPair(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("816c02c3d1b30fb10d1bb493d8e4d2af5731449c") //bridge contract
	assert.Nil(t, err)
	oep4TokenAddr, err := common.AddressFromHexString("077fc2906e1f2c9beab3ca1c2c7de191f81494d6") //oep4 token contract
	assert.Nil(t, err)
	oep4Decimals := 9
	erc20TokenAddr, err := common.AddressFromHexString("a7b848e1c02a927d65dd2746debd1a9c9898a48c") //erc20 token contract:0x8ca498989c1abdde4627dd657d922ac0e148b8a7
	assert.Nil(t, err)
	erc20Decimals := 18
	tokenPariName := "MKT-MKT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair",
		[]interface{}{
			tokenPariName,
			oep4TokenAddr,
			oep4Decimals,
			erc20TokenAddr,
			erc20Decimals,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestBridgeRegisterOtherOep4Erc20TokenPair(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("816c02c3d1b30fb10d1bb493d8e4d2af5731449c") //bridge contract
	assert.Nil(t, err)
	oep4TokenAddr, err := common.AddressFromHexString("22757e9ceb405f2d3b3d8a8a656aa7b1f8c68534") //oep4 token contract
	assert.Nil(t, err)
	oep4Decimals := 9
	erc20TokenAddr, err := common.AddressFromHexString("b0d3508c103cbc084600baa2be430f4c91936fab") //erc20 token contract:0xab6f93914c0f43bea2ba004608bc3c108c50d3b0
	assert.Nil(t, err)
	erc20Decimals := 18
	tokenPariName := "MKR-MKR"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair",
		[]interface{}{
			tokenPariName,
			oep4TokenAddr,
			oep4Decimals,
			erc20TokenAddr,
			erc20Decimals,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

//GetAllTokenPair
func TestGetAllTokenBridgeFromWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("816c02c3d1b30fb10d1bb493d8e4d2af5731449c") //wasm bridge testnet
	assert.Nil(t, err)
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getAllTokenPairName", []interface{}{})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	t.Logf("ResultItem value: %v", res.Result) // Hex string before decoding

	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("Raw bytes (hex): %x", bs)
	t.Logf("Raw bytes length: %d", len(bs))

	tokenPairNames, err := parseWasmVecVecU8(bs, t)
	if err != nil {
		t.Logf("Parse error: %v", err)
	}
	assert.Nil(t, err)

	// Log parsed token pair names
	for i, name := range tokenPairNames {
		t.Logf("name:%s", string(name))
		t.Logf("Token pair #%d: %s (hex: %x, length: %d)", i, string(name), name, len(name))
	}

}

//getTokenPairInfo
func TestGetTokenPairFromWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("816c02c3d1b30fb10d1bb493d8e4d2af5731449c") //wasm uniswapv2
	assert.Nil(t, err)
	//tokenPairName := "MKT-MKT"
	tokenPairName := "MKR-MKR"
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getTokenPair", []interface{}{tokenPairName})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("Raw bytes (hex): %x", bs)
	t.Logf("Raw bytes length: %d", len(bs))
	tokenPairInfo, err := parseTokenPairInfo(bs)
	assert.Nil(t, err)
	t.Logf("%v", tokenPairInfo)
}

// TokenPair represents the structure of the TokenPair in the contract
type TokenPairInfo struct {
	Erc20Address  string
	Erc20Decimals uint32
	Oep4Address   string
	Oep4Decimals  uint32
}

func parseTokenPairInfo(rawResult []byte) (*TokenPairInfo, error) {
	if len(rawResult) < 48 { // 20 (erc20) + 4 (erc20_decimals) + 20 (oep4) + 4 (oep4_decimals)
		return nil, fmt.Errorf("invalid result length: %d", len(rawResult))
	}
	tokenPairInfo := &TokenPairInfo{}
	erc20Address, err := common.AddressParseFromBytes(rawResult[0:20])
	if err != nil {
		return nil, fmt.Errorf("AddressParseFromBytes, err: %s", err)
	}
	tokenPairInfo.Erc20Address = erc20Address.ToHexString()
	tokenPairInfo.Erc20Decimals = binary.LittleEndian.Uint32(rawResult[20:24])
	oep4Address, err := common.AddressParseFromBytes(rawResult[24:44])
	if err != nil {
		return nil, fmt.Errorf("AddressParseFromBytes, err: %s", err)
	}
	tokenPairInfo.Oep4Address = oep4Address.ToHexString()
	tokenPairInfo.Oep4Decimals = binary.LittleEndian.Uint32(rawResult[44:48])
	return tokenPairInfo, nil
}