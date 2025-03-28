package ontology_go_sdk

import (
	"io/ioutil"
	"math/big"
	"testing"
	"time"
	"fmt"
	"encoding/binary"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

func TestDeploySwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	wasmfile := "./uniswapv2.wasm"
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
		"wasmuniswap",
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

func TestInitSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //uniswapv2
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "init", []interface{}{testDefAcc.Address})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestBalanceOfBridgeWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("a0c0e665a4a87c48dbb0b2e169e5ab7dd9c68f61") //wasm uniswap v2
	assert.Nil(t, err)
	wong, err := common.AddressFromHexString("197239029febc9158bd65e46d6509175832da6a5") //wong contract
	assert.Nil(t, err)
	user, err := common.AddressFromHexString("357be1f9e1c98b83f1ca1e363b79743019936d0f")
	assert.Nil(t, err)
	tx, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "balanceOf", []interface{}{wong, user})
	assert.Nil(t, err)
	t.Logf("txHash:%s", tx.ToHexString())
}

func TestPauseSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasmoep4 bridge
	assert.Nil(t, err)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "unpause", []interface{}{})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

//test succ
//https://explorer.ont.io/testnet/tx/7e5dba27c2da0b7f650704bc6f155881f273b17a39ea6bad60682dbf6552f3f5
func TestBalanceOfSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap v2
	assert.Nil(t, err)
	wong, err := common.AddressFromHexString("a5a62d83759150d6465ed68b15c9eb9f02397219") //0x197239029febc9158bd65e46d6509175832da6a5
	assert.Nil(t, err)
	user, err := common.AddressFromHexString("0f6d93193074793b361ecaf1838bc9e1f9e17b35") //0x357be1f9e1c98b83f1ca1e363b79743019936d0f
	assert.Nil(t, err)
	tx, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "balanceOf", []interface{}{wong, user})
	assert.Nil(t, err)
	t.Logf("txHash:%s", tx.ToHexString())
}

func TestGetTxEvent(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap v2
	assert.Nil(t, err)
	user, err := common.AddressFromHexString("357be1f9e1c98b83f1ca1e363b79743019936d0f")
	assert.Nil(t, err)
	//txHash := "69ee725996c514e1d4d89763a4513bd317a6fd0b2ad524c58494b7dfa952e152" //ont balanceOf
	txHash := "2c480d22ef813fa89bdf77d2a3d5a2c7edd4b7031fc8b82f60826d2e2cfb54ac"
	event, err := testOntSdk.GetSmartContractEvent(txHash)
	assert.Nil(t, err)
	t.Logf("event:%v", event)
	if event != nil && len(event.Notify) > 0 {
		for _, notify := range event.Notify {
			t.Logf("notify contract address:%s", notify.ContractAddress)
			if notify.ContractAddress == contractAddr.ToHexString() {
				bs := notify.States.([]byte)
				t.Logf("Raw result bytes: %x, length: %d", bs, len(bs))
				if len(bs) >= 16 {
					balance := new(big.Int).SetBytes(bs[len(bs)-16:]) // 解析 U128
					t.Logf("Balance of %s: %s", user.ToHexString(), balance.String())
				} else {
					t.Errorf("Invalid return length: %d", len(bs))
				}
			}
		}
	} else {
		t.Logf("No event found for tx: %s", txHash)
	}
}

func TestRegisterTokenPairSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap v2
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

func TestUnRegisterTokenPairSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "unregisterTokenPair", []interface{}{tokenPairName})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

type SwapTokenPair struct {
	Token0       common.Address
	Token1       common.Address
	ContractAddr common.Address
}

func TestGetTokenPairSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswapv2
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getTokenPair", []interface{}{tokenPairName})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("Raw bytes (hex): %x", bs)
	t.Logf("Raw bytes length: %d", len(bs))
	tokenPair, err := parseTokenPair(bs, t)
	if err != nil {
		t.Logf("Parse error: %v", err)
	}
	assert.Nil(t, err)
	// Log the parsed TokenPair
	t.Logf("TokenPair - Token0: %s", tokenPair.Token0.ToHexString())
	t.Logf("TokenPair - Token1: %s", tokenPair.Token1.ToHexString())
	t.Logf("TokenPair - ContractAddr: %s", tokenPair.ContractAddr.ToHexString())
}

func TestGetAllTokenPairSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap
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

func TestSwapExactTokensForTokensSwapWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap v2
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	wong, err := common.AddressFromHexString("a5a62d83759150d6465ed68b15c9eb9f02397219") //0x197239029febc9158bd65e46d6509175832da6a5
	assert.Nil(t, err)
	wont, err := common.AddressFromHexString("193c28b10b0f778d007744f1249520389061e54c") //0x4ce5619038209524f14477008d770f0bb1283c19
	assert.Nil(t, err)
	toAddr, err := common.AddressFromHexString("0f6d93193074793b361ecaf1838bc9e1f9e17b35") //0x357be1f9e1c98b83f1ca1e363b79743019936d0f
	assert.Nil(t, err)
	ethAddr, err := common.AddressFromHexString("25350e95e87d0b385d82e42b43f0ef272b484b5f") //0x5f4b482b27eff0432be4825d380b7de8950e3525
	assert.Nil(t, err)
	evmUniswapAddr, err := common.AddressFromHexString("b1ecc98a3ffedb9bd652981df286b1e971e6cf12") //evm swap contract 0x12cfe671e9b186f21d9852d69bdbfe3f8ac9ecb1
	assert.Nil(t, err)
	amountOutMin := big.NewInt(0) // 最小输出代币数量（设为 0 表示接受任何数量）
	amountIn := big.NewInt(1e16)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "swapExactTokensForTokens",
		[]interface{}{
			testDefAcc.Address,
			ethAddr,
			evmUniswapAddr,
			tokenPairName,
			amountIn,
			amountOutMin,
			[]interface{}{wong, wont},
			toAddr,
			big.NewInt(time.Now().Unix() + 3600),
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestTransferErc20WasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	gasprice := uint64(2500)
	invokegaslimit := uint64(200000)
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap v2
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	wong, err := common.AddressFromHexString("a5a62d83759150d6465ed68b15c9eb9f02397219") //0x197239029febc9158bd65e46d6509175832da6a5
	assert.Nil(t, err)
	wont, err := common.AddressFromHexString("193c28b10b0f778d007744f1249520389061e54c") //0x4ce5619038209524f14477008d770f0bb1283c19
	assert.Nil(t, err)
	toAddr, err := common.AddressFromHexString("0f6d93193074793b361ecaf1838bc9e1f9e17b35") //0x357be1f9e1c98b83f1ca1e363b79743019936d0f
	assert.Nil(t, err)
	ethAddr, err := common.AddressFromHexString("25350e95e87d0b385d82e42b43f0ef272b484b5f") //0x5f4b482b27eff0432be4825d380b7de8950e3525
	assert.Nil(t, err)
	amountOutMin := big.NewInt(0) // 最小输出代币数量（设为 0 表示接受任何数量）
	amountIn := big.NewInt(1e16)
	txHash, err := testOntSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, testDefAcc, contractAddr, "transferErc20",
		[]interface{}{
			testDefAcc.Address,
			ethAddr,
			tokenPairName,
			amountIn,
			amountOutMin,
			[]interface{}{wong, wont},
			toAddr,
			big.NewInt(time.Now().Unix() + 3600),
		})
	assert.Nil(t, err)
	t.Logf("txHash:%s", txHash.ToHexString())
}

func TestGetAmountsOutWasmContract(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	contractAddr, err := common.AddressFromHexString("1a8e70da5a613d6176689be54dbbc9cc22cfbc6b") //wasm uniswap v2
	assert.Nil(t, err)
	tokenPairName := "ONG_TO_WONT"
	amountIn := big.NewInt(1e16)
	res, err := testOntSdk.WasmVM.PreExecInvokeWasmVMContract(
		contractAddr, "getAmountsOut",
		[]interface{}{
			tokenPairName,
			amountIn})
	assert.Nil(t, err)
	t.Logf("PreExec result (full): %+v", res)
	bs, err := res.Result.ToByteArray()
	assert.Nil(t, err)
	t.Logf("Raw bytes (hex): %x", bs)
	t.Logf("Raw bytes length: %d", len(bs))

    // 解析 Vec<U128>
    amounts, err := parseWasmVecU128(bs)
	assert.Nil(t, err)
    // 打印结果
    for i, amount := range amounts {
		t.Logf("Amount[%d]: %s", i, amount.String())
    }
}

func parseWasmVecU128(res []byte) ([]*big.Int, error) {
    if len(res) < 4 {
        return nil, fmt.Errorf("invalid result length: too short")
    }

    // 读取长度 (u32, 4 字节)
    length := binary.BigEndian.Uint32(res[0:4])
    data := res[4:]

    // 验证数据长度
    if len(data) != int(length)*16 {
        return nil, fmt.Errorf("invalid data length: expected %d, got %d", length*16, len(data))
    }

    // 解析每个 U128
    amounts := make([]*big.Int, length)
    for i := uint32(0); i < length; i++ {
        start := i * 16
        end := start + 16
        // 假设 U128 是大端编码
        amount := new(big.Int).SetBytes(data[start:end])
        amounts[i] = amount
    }

    return amounts, nil
}