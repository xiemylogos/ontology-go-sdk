package ontology_go_sdk

import (
	"testing"

	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

func TestInit_EnsWasmContract(t *testing.T) {
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

func TestRegister_EnsWasmContract(t *testing.T) {
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

func TestResolve_EnsWasmContract(t *testing.T) {
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
