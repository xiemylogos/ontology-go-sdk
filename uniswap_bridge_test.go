package ontology_go_sdk

import (
	"encoding/binary"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTokenBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm bridge testnet
	//contractAddr, err := common.AddressFromHexString("7340c7b8d611e7fca29c351103b10c21ffdb1e3c") //wasm bridge mainnet
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

func TestGetTokenPairBridgeContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm uniswapv2 testnet
	//contractAddr, err := common.AddressFromHexString("7340c7b8d611e7fca29c351103b10c21ffdb1e3c") //wasm bridge mainnet
	assert.Nil(t, err)
	tokenPairName := "WING-WING"
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getTokenPair", []interface{}{tokenPairName})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("Raw bytes (hex): %x", bs)
	t.Logf("Raw bytes length: %d", len(bs))
	tokenPair := parseBridgeTokenPair(bs, t)
	if err != nil {
		t.Logf("Parse error: %v", err)
	}
	assert.Nil(t, err)
	// Log the parsed TokenPair
	// 输出解析结果
	t.Logf("TokenPair:")
	t.Logf("ERC20 Address:%x", tokenPair.Erc20)
	t.Logf("ERC20 Decimals: %d", tokenPair.Erc20Decimals)
	t.Logf("OEP4 Address: %x", tokenPair.Oep4)
	t.Logf("OEP4 Decimals: %d", tokenPair.Oep4Decimals)
}

// TokenPair 定义与合约中的 TokenPair 结构体对应的 Go 结构
type BridgeTokenPair struct {
	Erc20         [20]byte // 20 字节地址
	Erc20Decimals uint32   // 4 字节无符号整数
	Oep4          [20]byte // 20 字节地址
	Oep4Decimals  uint32   // 4 字节无符号整数
}

func parseBridgeTokenPair(data []byte, t *testing.T) BridgeTokenPair {
	// 默认值
	defaultPair := BridgeTokenPair{}

	// 检查数据长度，至少需要 64 字节
	if len(data) < 64 {
		return defaultPair
	}

	// 从第 20 字节开始解析，跳过前 20 字节
	data = data[20:68] // 取 48 字节 (20 + 48 = 68)

	// 解析字段
	var pair BridgeTokenPair

	// erc20: 前 20 字节
	copy(pair.Erc20[:], data[0:20])

	// erc20_decimals: 接下来的 4 字节 (大端序)
	pair.Erc20Decimals = binary.LittleEndian.Uint32(data[20:24])

	// oep4: 接下来的 20 字节
	copy(pair.Oep4[:], data[24:44])

	// oep4_decimals: 最后 4 字节 (大端序)
	pair.Oep4Decimals = binary.LittleEndian.Uint32(data[44:48])

	return pair
}

func TestDeploySwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	wasmfile := "./uniswapv2bridge.wasm"
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

func TestInitSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //uniswapv2gatway
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "init", []interface{}{testDefAcc.Address})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestOngWrapperSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("852b906fab77aba58414527ec5f729f628ffbde1") //wasm uniswap bridge
	assert.Nil(t, err)
	wong, err := common.AddressFromHexString("a5a62d83759150d6465ed68b15c9eb9f02397219") //0x197239029febc9158bd65e46d6509175832da6a5
	assert.Nil(t, err)
	amount := big.NewInt(7e16)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "ong_wrapper", []interface{}{testDefAcc.Address, wong, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestSwapExactTokensForTokensSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //wasm uniswap bridge
	assert.Nil(t, err)
	wong, err := common.AddressFromHexString("a5a62d83759150d6465ed68b15c9eb9f02397219") //0x197239029febc9158bd65e46d6509175832da6a5
	assert.Nil(t, err)
	wont, err := common.AddressFromHexString("193c28b10b0f778d007744f1249520389061e54c") //0x4ce5619038209524f14477008d770f0bb1283c19
	assert.Nil(t, err)
	evmUniswapAddr, err := common.AddressFromHexString("b1ecc98a3ffedb9bd652981df286b1e971e6cf12") //evm swap contract 0x12cfe671e9b186f21d9852d69bdbfe3f8ac9ecb1
	assert.Nil(t, err)
	wasmBridgeAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm birdge contract
	assert.Nil(t, err)
	amountOutMin := big.NewInt(0) // 最小输出代币数量（设为 0 表示接受任何数量）
	amountIn := big.NewInt(1e16)
	wasmbrigeTokenPariName := "ONT-WONT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "swapExactTokensForTokens",
		[]interface{}{
			testDefAcc.Address,
			wasmbrigeTokenPariName,
			wasmBridgeAddr,
			evmUniswapAddr,
			amountIn,
			amountOutMin,
			[]interface{}{wong, wont},
			big.NewInt(time.Now().Unix() + 3600),
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetAmountsOutWaspBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //wasm uniswap v2
	assert.Nil(t, err)
	evmUniswapAddr, err := common.AddressFromHexString("b1ecc98a3ffedb9bd652981df286b1e971e6cf12") //evm swap contract 0x12cfe671e9b186f21d9852d69bdbfe3f8ac9ecb1
	assert.Nil(t, err)
	wasmBridgeAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm birdge contract
	assert.Nil(t, err)
	amountIn := big.NewInt(1e16)
	wasmbrigeTokenPariName1 := "ONG"
	wasmbrigeTokenPariName2 := "ONT-WONT"
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getAmountsOut",
		[]interface{}{
			evmUniswapAddr,
			wasmBridgeAddr,
			amountIn,
			[]interface{}{wasmbrigeTokenPariName1, wasmbrigeTokenPariName2},
		})
	assert.Nil(t, err)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	result, err := parseEvmUintArrayInfo(bs)
	assert.Nil(t, err)
	for _, u128 := range result {
		t.Logf("value:%s", u128.String())
	}
}

func TestNewSwapExactTokensForTokensSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //wasm uniswap bridge
	assert.Nil(t, err)
	evmUniswapAddr, err := common.AddressFromHexString("b1ecc98a3ffedb9bd652981df286b1e971e6cf12") //evm swap contract 0x12cfe671e9b186f21d9852d69bdbfe3f8ac9ecb1
	assert.Nil(t, err)
	wasmBridgeAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm birdge contract
	assert.Nil(t, err)
	amountOutMin := big.NewInt(0) // 最小输出代币数量（设为 0 表示接受任何数量）
	//swap from ong->ont
	amountIn := big.NewInt(1e16)
	wasmbrigeTokenPariName1 := "ONG"
	wasmbrigeTokenPariName2 := "ONT-WONT"
	//swap from ont->ong
	/*
			amountIn := big.NewInt(1e7)
		wasmbrigeTokenPariName1 := "ONT-WONT"
		wasmbrigeTokenPariName2 := "ONG"
	*/
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "newSwapExactTokensForTokens",
		[]interface{}{
			testDefAcc.Address,
			evmUniswapAddr,
			wasmBridgeAddr,
			amountIn,
			amountOutMin,
			[]interface{}{wasmbrigeTokenPariName1, wasmbrigeTokenPariName2},
			big.NewInt(time.Now().Unix() + 3600),
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetAllTokenPariSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //wasm uniswap bridge
	assert.Nil(t, err)
	wasmBridgeAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm birdge contract
	assert.Nil(t, err)
	ong, err := common.AddressFromHexString("0200000000000000000000000000000000000000") //0200000000000000000000000000000000000000
	assert.Nil(t, err)
	ont, err := common.AddressFromHexString("0100000000000000000000000000000000000000") //0100000000000000000000000000000000000000
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "getBridgeAllTokenPairs", []interface{}{wasmBridgeAddr, []interface{}{ong, ont}})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetBridgeErc20AddrSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //wasm uniswap bridge
	assert.Nil(t, err)
	wasmBridgeAddr, err := common.AddressFromHexString("724191fd4acc22e6b5c813ea5fa4cd161810f75d") //wasm birdge contract
	assert.Nil(t, err)
	wasmbrigeTokenPariName1 := "ONG"
	wasmbrigeTokenPariName2 := "ONT-WONT"
	//wasmbrigeTokenPariName := "ONT-WONT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "getBridgeErc20Addr", []interface{}{wasmBridgeAddr, []interface{}{wasmbrigeTokenPariName1, wasmbrigeTokenPariName2}})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestOngWithDrawSwapBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("6b0712beb0140a3a7a591857d5f99d807fb280f9") //wasm uniswap bridge
	assert.Nil(t, err)
	amount := big.NewInt(7e16)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "withdrawOng", []interface{}{testDefAcc.Address, amount})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}
