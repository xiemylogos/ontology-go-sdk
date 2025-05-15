package ontology_go_sdk

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/big"
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
	timeout := 10 * time.Second
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
	contractAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm bridge
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

// GetAllTokenPair
func TestGetAllTokenBridgeFromWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("0c5ea3898e9a5a32335818e322eed98cfe266b44") //wasm bridge testnet
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

// getTokenPairInfo
func TestGetTokenPairFromWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("0c5ea3898e9a5a32335818e322eed98cfe266b44") //wasm uniswapv2
	assert.Nil(t, err)
	//tokenPairName := "MKT-MKT"
	//tokenPairName := "MKR-MKR"
	tokenPairName := "MYT-MYT"
	//tokenPairName := "ONT-WONT"
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

func TestBridgeRegisterNeoVmMYTTokenPair(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("b09746c446d87e218f0cf04dfe51db0aa977bc72") //bridge contract
	assert.Nil(t, err)
	//https://explorer.ont.io/testnet/contract/other/19ed155fd0adf1d598af054068538d41a3a192ff/10/1
	oep4TokenAddr, err := common.AddressFromHexString("19ed155fd0adf1d598af054068538d41a3a192ff") //neovm oep4 token contract
	assert.Nil(t, err)
	oep4Decimals := 18
	erc20TokenAddr, err := common.AddressFromHexString("a7b848e1c02a927d65dd2746debd1a9c9898a48c") //erc20 token contract:0x8ca498989c1abdde4627dd657d922ac0e148b8a7
	assert.Nil(t, err)
	erc20Decimals := 18
	tokenPairName := "MYT-MYT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair",
		[]interface{}{
			tokenPairName,
			oep4TokenAddr,
			oep4Decimals,
			erc20TokenAddr,
			erc20Decimals,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestBridgeRegisterOntErc20TokenPair(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("b09746c446d87e218f0cf04dfe51db0aa977bc72") //bridge contract
	assert.Nil(t, err)
	//https://explorer.ont.io/testnet/contract/other/0100000000000000000000000000000000000000/10/1
	oep4TokenAddr, err := common.AddressFromHexString("0100000000000000000000000000000000000000") //oep4 token contract
	assert.Nil(t, err)
	oep4Decimals := 9
	//wont
	//https://explorer.ont.io/testnet/contract/orc20/0x4ce5619038209524f14477008d770f0bb1283c19/10/1
	erc20TokenAddr, err := common.AddressFromHexString("193c28b10b0f778d007744f1249520389061e54c") //erc20 wont token contract:0x4ce5619038209524f14477008d770f0bb1283c19
	assert.Nil(t, err)
	erc20Decimals := 9
	tokenPairName := "ONT-WONT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair",
		[]interface{}{
			tokenPairName,
			oep4TokenAddr,
			oep4Decimals,
			erc20TokenAddr,
			erc20Decimals,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetBalanceWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("b09746c446d87e218f0cf04dfe51db0aa977bc72") //wasm bridge contract address
	assert.Nil(t, err)
	oep4ContractAddr, err := common.AddressFromHexString("19ed155fd0adf1d598af054068538d41a3a192ff")
	assert.Nil(t, err)
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "balanceOf", []interface{}{oep4ContractAddr, testDefAcc.Address})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("balance:%v", binary.LittleEndian.Uint64(bs))
}

func TestTransferWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("b09746c446d87e218f0cf04dfe51db0aa977bc72") //wasm bridge contract
	assert.Nil(t, err)
	oep4ContractAddr, err := common.AddressFromHexString("19ed155fd0adf1d598af054068538d41a3a192ff")
	assert.Nil(t, err)
	amount := big.NewInt(1)
	toAddr, err := common.AddressFromBase58("AcvHRs4SVqBh5dNnusakmUTu6uUCpfdZRt")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "transfer",
		[]interface{}{
			testDefAcc.Address,
			toAddr,
			amount,
			oep4ContractAddr,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestBirdgeRegisterWasmOep4Token(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm bridge contract
	assert.Nil(t, err)
	oep4TokenAddr, err := common.AddressFromHexString("f4adeee0867b34cf40f3a23233353307101b6c9b") //wasm oep4 token contract
	assert.Nil(t, err)
	oep4Decimals := 18
	erc20TokenAddr, err := common.AddressFromHexString("a7b848e1c02a927d65dd2746debd1a9c9898a48c") //erc20 token contract:0x8ca498989c1abdde4627dd657d922ac0e148b8a7
	assert.Nil(t, err)
	erc20Decimals := 18
	tokenPairName := "MKT-MKT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair",
		[]interface{}{
			tokenPairName,
			oep4TokenAddr,
			oep4Decimals,
			erc20TokenAddr,
			erc20Decimals,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestBirdgeRegisterWasmOep4TokenTwo(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm bridge contract
	assert.Nil(t, err)
	oep4TokenAddr, err := common.AddressFromHexString("97e32e59b0af10ea8575bab82ed8156f47c91a06") //wasm oep4 token contract
	assert.Nil(t, err)
	oep4Decimals := 18
	erc20TokenAddr, err := common.AddressFromHexString("b0d3508c103cbc084600baa2be430f4c91936fab") //erc20 token contract:0xab6f93914c0f43bea2ba004608bc3c108c50d3b0
	assert.Nil(t, err)
	erc20Decimals := 18
	tokenPairName := "MKR-MKR"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair",
		[]interface{}{
			tokenPairName,
			oep4TokenAddr,
			oep4Decimals,
			erc20TokenAddr,
			erc20Decimals,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetTokenPairForWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm birdge
	assert.Nil(t, err)
	tokenPairName := "MKT-MKT"
	//tokenPairName := "MKR-MKR"
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getTokenPair", []interface{}{tokenPairName})
	assert.Nil(t, err)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	tokenPair := parseBridgeTokenPair(bs, t)
	assert.Nil(t, err)
	t.Logf("TokenPair: erc20Addr:%x,decimal:%d,oep4Addr:%x,decimal:%d", tokenPair.Erc20, tokenPair.Erc20Decimals, tokenPair.Oep4, tokenPair.Oep4Decimals)
}


func TestWasmOep4ToOrc20(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	wasmBridgeAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm birdge contract
	assert.Nil(t, err)
	amountIn := big.NewInt(6e18)
	//tokenPairName := "MKT-MKT"
	tokenPairName := "MKR-MKR"
	ethAddr, err := common.AddressFromHexString("d7214bfe969b98d236e72d4f6ca3c912b1499812")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, wasmBridgeAddr, "wasmOep4ToOrc20",
		[]interface{}{
			testDefAcc.Address,
			ethAddr,
			amountIn,
			tokenPairName,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetWasmBalanceWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm bridge contract address
	assert.Nil(t, err)
	oep4ContractAddr, err := common.AddressFromHexString("f4adeee0867b34cf40f3a23233353307101b6c9b")
	assert.Nil(t, err)
	userAddr, err := common.AddressFromBase58("AcvHRs4SVqBh5dNnusakmUTu6uUCpfdZRt")
	assert.Nil(t, err)
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "wasmBalanceOf", []interface{}{oep4ContractAddr, userAddr})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("balance:%v", binary.LittleEndian.Uint64(bs))
}

func TestTransferWasmOep4UseWasmBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm bridge contract
	assert.Nil(t, err)
	oep4ContractAddr, err := common.AddressFromHexString("f4adeee0867b34cf40f3a23233353307101b6c9b")
	assert.Nil(t, err)
	amount := big.NewInt(1)
	toAddr, err := common.AddressFromBase58("AcvHRs4SVqBh5dNnusakmUTu6uUCpfdZRt")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "wasmTransfer",
		[]interface{}{
			testDefAcc.Address,
			toAddr,
			amount,
			oep4ContractAddr,
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestNewSwapExactOep4TokensForWasmTokensSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("930bb8674dc8651dc9a6b8c686f5da8a8ca2b467") //wasm uniswap bridge(support wasm oep4token)
	assert.Nil(t, err)
	evmUniswapAddr, err := common.AddressFromHexString("fd7333783574ef3a4e5d19289b706123bfb9e8f9") //evm swapv2 router contract 0xf9e8b9bf2361709b28195d4e3aef7435783373fd
	assert.Nil(t, err)
	wasmBridgeAddr, err := common.AddressFromHexString("1febb9e5f7dec0b8026629d02d33070884c0228f") //wasm birdge contract
	assert.Nil(t, err)
	amountOutMin := big.NewInt(0) // 最小输出代币数量（设为 0 表示接受任何数量）
	//swap(oep4->ope4) from MKT->MKR
	amountIn := big.NewInt(7)
	wasmbrigeTokenPariName1 := "MKT-MKT"
	wasmbrigeTokenPariName2 := "MKR-MKR"

	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "swapExactTokensForTokens",
		[]interface{}{
			testDefAcc.Address,
			evmUniswapAddr,
			wasmBridgeAddr,
			amountIn,
			amountOutMin,
			[]interface{}{wasmbrigeTokenPariName2, wasmbrigeTokenPariName1},
			big.NewInt(time.Now().Unix() + 3600),
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
	//https://explorer.ont.io/testnet/tx/207dfe10373bca3dd00cb6c78105b16807be2c6e93ef89869afd3dea3c7f349b
}