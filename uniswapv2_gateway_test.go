package ontology_go_sdk

import (
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)
func TestDeploySwapGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	wasmfile := "./wasm_file/uniswapv2gateway.wasm"
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
		"wasmuniswapgateway",
		"1.0",
		"author",
		"email",
		"desc",
	)
	assert.Nil(t, err)
	timeout := 90 * time.Second
	_, err = testOntSdk.WaitForGenerateBlock(timeout)
	assert.Nil(t, err)
	t.Logf("deploy wasm contract txhash is %s", txHash.ToHexString())
	contractAddr, err := utils.GetContractAddress(codeHash)
	assert.Nil(t, err)
	t.Logf("the contractAddr is:%s", contractAddr.ToHexString())
}

func TestInitSwapGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("68271157db1557af35f6dd8493204041cdb07f32") //uniswapv2gatway
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "init", []interface{}{testDefAcc.Address})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestRegisterTokenPairSwapGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("68271157db1557af35f6dd8493204041cdb07f32") //wasm uniswap v2gateway
	assert.Nil(t, err)
	evmUniswapAddr, err := common.AddressFromHexString("b1ecc98a3ffedb9bd652981df286b1e971e6cf12") //evm swap contract 0x12cfe671e9b186f21d9852d69bdbfe3f8ac9ecb1
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	wong, err := common.AddressFromHexString("a5a62d83759150d6465ed68b15c9eb9f02397219") //0x197239029febc9158bd65e46d6509175832da6a5
	assert.Nil(t, err)
	wont, err := common.AddressFromHexString("193c28b10b0f778d007744f1249520389061e54c") //0x4ce5619038209524f14477008d770f0bb1283c19
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "registerTokenPair", []interface{}{tokenPairName, wong, wont, evmUniswapAddr})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestSwapExactTokensForTokensSwapGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("68271157db1557af35f6dd8493204041cdb07f32") //wasm uniswap v2
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	toAddr, err := common.AddressFromHexString("0f6d93193074793b361ecaf1838bc9e1f9e17b35") //0x357be1f9e1c98b83f1ca1e363b79743019936d0f
	assert.Nil(t, err)
	ethAddr, err := common.AddressFromHexString("25350e95e87d0b385d82e42b43f0ef272b484b5f") //0x5f4b482b27eff0432be4825d380b7de8950e3525
	assert.Nil(t, err)
	amountOutMin := big.NewInt(0)
	amountIn := big.NewInt(1e16)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "swapExactTokensForTokens",
		[]interface{}{
			testDefAcc.Address,
			ethAddr,
			tokenPairName,
			amountIn,
			amountOutMin,
			toAddr,
			big.NewInt(time.Now().Unix() + 3600),
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestPauseUniswapGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("68271157db1557af35f6dd8493204041cdb07f32") //wasmoep4 bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "pause", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestSetAdminUniswapGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("68271157db1557af35f6dd8493204041cdb07f32") //wasmoep4 bridge
	assert.Nil(t, err)
	admin, err := common.AddressFromBase58("AZqW2S7FRp8n45V48YWkuLe7BnkqBaXbyP")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "setAdmin", []interface{}{admin})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetAmountsOutWaspGatewayWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	contractAddr, err := common.AddressFromHexString("68271157db1557af35f6dd8493204041cdb07f32") //wasm uniswap v2
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	amountIn := big.NewInt(1e16)
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getAmountsOut",
		[]interface{}{
			tokenPairName,
			amountIn})
	assert.Nil(t, err)
	bs,err := res.Result.ToByteArray()
	assert.Nil(t, err)
	result, err := parseEvmUintArrayInfo(bs)
	assert.Nil(t, err)
	for _, u128 := range result {
		t.Logf("value:%s", u128.String())
	}
}

