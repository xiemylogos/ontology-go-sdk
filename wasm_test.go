package ontology_go_sdk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethcomm "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ontio/ontology/common/constants"
	txtypes "github.com/ontio/ontology/core/types"

	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/core/types"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDid(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	did := "did:ont:" + testDefAcc.Address.ToBase58()
	txhash, err := testOntSdk.Native.OntId.RegIDWithPublicKey(testGasPrice, testGasLimit, testDefAcc, did, testDefAcc)
	assert.Nil(t, err)
	t.Logf("txHash:%s", txhash.ToHexString())
}

func TestGetPubKeyByDid(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	did := "did:ont:ALefsBqE3JCdDjatuaJuxWxDh8fUnHyDWC"
	publicKeys, err := testOntSdk.Native.OntId.GetPublicKeysJson(did)
	assert.Nil(t, err)
	var publicKeyList PublicKeyList
	err = json.Unmarshal(publicKeys, &publicKeyList)
	assert.Nil(t, err)
	for _, pkInfo := range publicKeyList {
		data, err := hex.DecodeString(pkInfo.PublicKeyHex)
		assert.Nil(t, err)
		pk, err := keypair.DeserializePublicKey(data)
		assert.Nil(t, err)
		pkBytes := keypair.SerializePublicKey(pk)
		t.Logf("pkBytes data:%s", hex.EncodeToString(pkBytes))
		address := types.AddressFromPubKey(pk)
		t.Logf("address:%s", address.ToBase58())
	}
}

func TestDeployWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	//wasmfile := "./oep4token.wasm"  //ope4token
	//wasmfile := "./oep4bridge.wasm"
	//wasmfile := "./crosscontract.wasm"
	//wasmfile := "./nativebridge.wasm"
	wasmfile := "./proxybridge.wasm"
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
		"wasmproxybridge",
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

//setLogicContract
func TestSetProxyBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	proxyContractAddr, err := common.AddressFromHexString("6d0e0deb199edbae1106145295fa67cfe0b2bc4b") //proxy bridge
	assert.Nil(t, err)
	implContractAddr, err := common.AddressFromHexString("e89c93487fd5faf8cac30c90eddcdad73a5ff3ce") //impl bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, proxyContractAddr, "setLogicContract", []interface{}{implContractAddr})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

//getLogicContract
func TestGetProxyBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	proxyContractAddr, err := common.AddressFromHexString("6d0e0deb199edbae1106145295fa67cfe0b2bc4b") //proxy bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, proxyContractAddr, "getLogicContract", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestFromProxyBridgeCallLockNativeBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	proxyContractAddr, err := common.AddressFromHexString("6d0e0deb199edbae1106145295fa67cfe0b2bc4b") //native bridge
	assert.Nil(t, err)
	nativeContract, err := common.AddressFromBase58("AFmseVrdL9f9oyCzZefL9tG6UbvhfRZMHJ") //define ONG contract addr
	assert.Nil(t, err)
	fromAddr := testDefAcc.Address
	assert.Nil(t, err)
	amount := big.NewInt(660)
	validChain := "Ethereum"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, proxyContractAddr, "call", []interface{}{"lock", nativeContract, fromAddr, validChain, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestDeployEnsWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	//wasmfile := "./oep4token.wasm"  //ope4token
	//wasmfile := "./oep4bridge.wasm"
	//wasmfile := "./crosscontract.wasm"
	//wasmfile := "./nativebridge.wasm"
	wasmfile := "./ensont.wasm"
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
		"wasmontens",
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

func TestInitEnsWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("2042c3b54c681e4ba728a022d01ab950612445ac") //ont ens contract
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "init", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestRegisterEnsWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	ensContractAddr, err := common.AddressFromHexString("2042c3b54c681e4ba728a022d01ab950612445ac") //ont contract addr
	assert.Nil(t, err)
	ensName := "alice.ont.io"
	owner, err := common.AddressFromBase58("ALefsBqE3JCdDjatuaJuxWxDh8fUnHyDWC")
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, ensContractAddr, "register", []interface{}{ensName, owner})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestResolveEnsWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	ensContractAddr, err := common.AddressFromHexString("2042c3b54c681e4ba728a022d01ab950612445ac") //ont ens contract
	assert.Nil(t, err)
	ensName := "alice.ont.io"
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(ensContractAddr, "resolve", []interface{}{ensName})
	assert.Nil(t, err)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	addr, err := common.AddressParseFromBytes(bs)
	assert.Nil(t, err)
	t.Logf("resolve addr:%s", addr.ToBase58())
}

func TestBuildEipTx(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	value := big.NewInt(1000000000)
	gaslimit := uint64(21000)
	gasPrice := big.NewInt(2500 * constants.GWei)
	toAddress := ethcomm.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	nonce := uint64(0)
	tx := ethtypes.NewTransaction(nonce, toAddress, value, gaslimit, gasPrice, data)
	chainId := big.NewInt(int64(5851))
	priKeyBytes := keypair.SerializePrivateKey(testDefAcc.PrivateKey)
	t.Logf("len:%d",len(priKeyBytes))
	if len(priKeyBytes) != 67 {
		t.Fatalf("Unexpected private key length: %d", len(priKeyBytes))
		return
	}
	privKeyD := priKeyBytes[1:33] // 跳过 1 字节类型标记，取 32 字节 D
	privKeyECDSA, err := crypto.ToECDSA(privKeyD)
	assert.Nil(t, err)
	signer := ethtypes.NewEIP155Signer(chainId)
	signedTx, err := ethtypes.SignTx(tx, signer, privKeyECDSA)
	assert.Nil(t, err)
	bts, _ := rlp.EncodeToBytes(signedTx)
	fmt.Printf("rlp:0x%s\n", hex.EncodeToString(bts))
	otx, err := txtypes.TransactionFromEIP155(signedTx)
	assert.Nil(t, err)
	t.Logf("tx hash:%s",otx.Hash().ToHexString())

	sender, err := signer.Sender(signedTx)
	assert.Nil(t, err)
	t.Logf("Recovered sender address: %s", sender.Hex())

	expectedAddr := crypto.PubkeyToAddress(privKeyECDSA.PublicKey)
	t.Logf("Expected address: %s", expectedAddr.Hex())

	if sender != expectedAddr {
		t.Fatalf("Signature verification failed: sender %s != expected %s", sender.Hex(), expectedAddr.Hex())
	}
	t.Logf("Signature verification passed: sender matches expected address")

	txHash := signer.Hash(signedTx)
	t.Logf("Transaction hash: %x", txHash)

	var parsedTx struct {
		Nonce    uint64
		GasPrice *big.Int
		GasLimit uint64
		To       common.Address
		Value    *big.Int
		Data     []byte
		V        *big.Int
		R, S     *big.Int
	}
	err = rlp.DecodeBytes(bts, &parsedTx)
	assert.Nil(t, err)

	t.Logf("Parsed signature values: v=%d, r=%x, s=%x", parsedTx.V, parsedTx.R, parsedTx.S)

	sig := append(ethcomm.LeftPadBytes(parsedTx.R.Bytes(), 32), ethcomm.LeftPadBytes(parsedTx.S.Bytes(), 32)...)

	pubKeyBytes := crypto.FromECDSAPub(&privKeyECDSA.PublicKey)
	if !crypto.VerifySignature(pubKeyBytes, txHash[:], sig) {
		t.Fatalf("Manual signature verification failed")
	}
	t.Logf("Manual signature verification passed")

}