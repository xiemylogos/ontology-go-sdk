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

func TestDeployWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	//wasmfile := "./oep4token.wasm"  //ope4token
	//wasmfile := "./oep4bridge.wasm"
	//wasmfile := "./crosscontract.wasm"
	wasmfile := "./nativebridge.wasm"
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
		"wasmoep4bridge",
		"1.0",
		"author",
		"email",
		"desc",
	)
	assert.Nil(t, err)
	timeout := 60 * time.Second
	_, err = testOntSdk.WaitForGenerateBlock(timeout)
	assert.Nil(t, err)
	t.Logf("deploy wasm contract txhash is %s", txHash.ToHexString())
	contractAddr, err := utils.GetContractAddress(codeHash)
	assert.Nil(t, err)
	t.Logf("the contractAddr is:%s", contractAddr.ToHexString())
}

func TestWasmContractAddr(t *testing.T) {
	wasmfile := "./oep4bridge.wasm"
	code, err := ioutil.ReadFile(wasmfile)
	assert.Nil(t, err)
	codeHash := common.ToHexString(code)
	contractAddr, err := utils.GetContractAddress(codeHash)
	assert.Nil(t, err)
	t.Logf("the contractAddr is:%s", contractAddr.ToHexString())
	t.Logf("the base58 contractAddr is:%s", contractAddr.ToBase58())
}
func TestInitWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	//	contractAddr, err := common.AddressFromHexString("b8c3d8d19edf354b6e3fd83bfc6e043099120ef3") //ope4token
	//contractAddr, err := common.AddressFromHexString("f1fb556f9eb49bbc379aaef3aac59ccc53420832") //ope4token
	contractAddr, err := common.AddressFromHexString("e89c93487fd5faf8cac30c90eddcdad73a5ff3ce") //oep4bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "init", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetInfoWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("f1fb556f9eb49bbc379aaef3aac59ccc53420832")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "getInfo", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestApproveOep4TokenWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("9bcd5c0b86b703d09254b8047af71407c4698930")
	assert.Nil(t, err)
	amount := big.NewInt(800000000000)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "approve", []interface{}{contractAddr, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestSetParamsWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("f1fb556f9eb49bbc379aaef3aac59ccc53420832")
	assert.Nil(t, err)
	/*
		txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
			gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "name", []interface{}{})
		assert.Nil(t, err)
		t.Logf("txHash:%s",txHash.ToHexString())

		time.Sleep(20*time.Second);
		txHashSymbol, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
			gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "symbol", []interface{}{})
		assert.Nil(t, err)
		t.Logf("txHash:%s",txHashSymbol.ToHexString())

		time.Sleep(20*time.Second);
	*/
	txHashTotoal, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "totalSupply", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHashTotoal.ToHexString())

}

func TestTransferWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("b8c3d8d19edf354b6e3fd83bfc6e043099120ef3")
	assert.Nil(t, err)
	fromAddr := testDefAcc.Address
	toAddr, err := common.AddressFromBase58("AVrvzC1Uax1QRTvNxkUDFZyC3AdCm7U9Un")
	assert.Nil(t, err)
	amount := big.NewInt(800000000000)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "transfer", []interface{}{fromAddr, toAddr, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestLockOep4TokenBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("f3e2b3d06cebcf962898c3c59523dc50248e3587") //wasmoep4 bridge
	assert.Nil(t, err)
	fromAddr := testDefAcc.Address
	assert.Nil(t, err)
	amount := big.NewInt(8000000)
	validChain := "Ethereum"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "lock", []interface{}{fromAddr, validChain, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestPauseBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("f3e2b3d06cebcf962898c3c59523dc50248e3587") //wasmoep4 bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "unpause", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestSetAdminBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("32acc7f4f99744930fbfa6f55f5fcb83db912ce7") //wasmoep4 bridge
	assert.Nil(t, err)
	admin, err := common.AddressFromBase58("AVrvzC1Uax1QRTvNxkUDFZyC3AdCm7U9Un")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "setAdmin", []interface{}{admin})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestCrossContractTransferWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("d6dfe924b8abd59b0cb4ffcfbc8fe83025bdb867")
	assert.Nil(t, err)
	fromAddr := testDefAcc.Address
	toAddr, err := common.AddressFromBase58("AVrvzC1Uax1QRTvNxkUDFZyC3AdCm7U9Un")
	assert.Nil(t, err)
	amount := big.NewInt(80000)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "callOep4Transfer", []interface{}{fromAddr, toAddr, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestPauseCrossContractWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("e89c93487fd5faf8cac30c90eddcdad73a5ff3ce") //wasmoep4 bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "unpause", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestLockNativeBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("e89c93487fd5faf8cac30c90eddcdad73a5ff3ce") //native bridge
	assert.Nil(t, err)
	nativeContract, err := common.AddressFromBase58("AFmseVrdL9f9oyCzZefL9tG6UbvhfRZMHJ") //define ONG contract addr
	assert.Nil(t, err)
	fromAddr := testDefAcc.Address
	assert.Nil(t, err)
	amount := big.NewInt(660)
	validChain := "Ethereum"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "lock", []interface{}{nativeContract, fromAddr, validChain, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}