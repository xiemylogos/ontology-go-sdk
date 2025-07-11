/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package ontology_go_sdk

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology/core/types"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/constants"
	"github.com/stretchr/testify/assert"
)

func TestOnt_BalanceV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	balance, err := testOntSdk.Native.Ont.BalanceOfV2(testDefAcc.Address)
	assert.Nil(t, err)
	t.Logf("balance:%d", balance)
}

func TestOng_BalanceV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	balance, err := testOntSdk.Native.Ong.BalanceOfV2(testDefAcc.Address)
	assert.Nil(t, err)
	t.Logf("balance:%v", balance)
}

func TestOnt_DecimalsV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	res, err := testOntSdk.Native.Ont.DecimalsV2()
	assert.Nil(t, err)
	assert.Equal(t, 9, int(res))
}

func TestOng_DecimalsV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	res, err := testOntSdk.Native.Ong.DecimalsV2()
	assert.Nil(t, err)
	assert.Equal(t, 18, int(res))
}

func TestOnt_TotalSupplyV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	res, err := testOntSdk.Native.Ont.TotalSupplyV2()
	assert.Nil(t, err)
	assert.Equal(t, new(big.Int).SetInt64(constants.ONT_TOTAL_SUPPLY_V2), res)
}

func TestOng_TotalSupplyV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	res, err := testOntSdk.Native.Ong.TotalSupplyV2()
	assert.Nil(t, err)
	t.Logf("rest:%v", res)
	assert.Equal(t, common.BigIntToNeoBytes(constants.ONG_TOTAL_SUPPLY_V2.BigInt()), bytesReverse(res.Bytes()))
}

func bytesReverse(u []byte) []byte {
	for i, j := 0, len(u)-1; i < j; i, j = i+1, j-1 {
		u[i], u[j] = u[j], u[i]
	}
	return u
}

func TestOng_UnboundONGV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	res, err := testOntSdk.Native.Ong.UnboundONGV2(testDefAcc.Address)
	assert.Nil(t, err)
	t.Logf("res:%v", res)
}

func TestOnt_TransferV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	addr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	txHash, err := testOntSdk.Native.Ont.TransferV2(testGasPrice, testGasLimit, nil, testDefAcc, addr, new(big.Int).SetInt64(98))
	assert.Nil(t, err)
	t.Logf("hash:%v", txHash.ToHexString())

}

func TestOng_TransferV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	addr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	txHash, err := testOntSdk.Native.Ong.TransferV2(testGasPrice, testGasLimit, nil, testDefAcc, addr, new(big.Int).SetInt64(10000000000887776))
	assert.Nil(t, err)
	t.Logf("hash:%v", txHash.ToHexString())
	_, err = testOntSdk.WaitForGenerateBlock(30 * time.Second)
	assert.Nil(t, err)
	contractEvent, err := testOntSdk.GetSmartContractEvent(txHash.ToHexString())
	assert.Nil(t, err)
	for _, notify := range contractEvent.Notify {
		transfer, err := testOntSdk.ParseNativeTransferEventV2(notify)
		assert.Nil(t, err)
		t.Logf("transfer:%v", transfer)
	}
}

// send ong tx test
func TestOng_Transfer(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./testconsens.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	addr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	txHash, err := testOntSdk.Native.Ong.TransferV2(testGasPrice, testGasLimit, nil, testDefAcc, addr, new(big.Int).SetInt64(17))
	assert.Nil(t, err)
	t.Logf("hash:%v", txHash.ToHexString())
}

func TestOnt_NewTransferTransactionV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	toAddr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	mutableTransaction, err := testOntSdk.Native.Ont.NewTransferTransactionV2(testGasPrice, testGasLimit, testDefAcc.Address, toAddr, new(big.Int).SetInt64(1111111111125))
	assert.Nil(t, err)
	ontTx, err := mutableTransaction.IntoImmutable()
	assert.Nil(t, err)
	res, err := ParseNativeTxPayloadV2(ontTx.ToArray())
	assert.Nil(t, err)
	t.Logf("res:%v", res)
}

func TestOnt_NewTransferTransaction(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	toAddr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	mutableTransaction, err := testOntSdk.Native.Ont.NewTransferTransaction(testGasPrice, testGasLimit, testDefAcc.Address, toAddr, 111111125)
	assert.Nil(t, err)
	ontTx, err := mutableTransaction.IntoImmutable()
	assert.Nil(t, err)
	res, err := ParseNativeTxPayload(ontTx.ToArray())
	assert.Nil(t, err)
	t.Logf("res:%v", res)
}

func TestOng_NewTransferTransactionV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	toAddr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	mutableTransaction, err := testOntSdk.Native.Ong.NewTransferTransactionV2(testGasPrice, testGasLimit, testDefAcc.Address, toAddr, new(big.Int).SetInt64(889100112999999025))
	assert.Nil(t, err)
	ongTx, err := mutableTransaction.IntoImmutable()
	assert.Nil(t, err)
	res, err := ParseNativeTxPayloadV2(ongTx.ToArray())
	assert.Nil(t, err)
	t.Logf("res:%v", res)
}

func TestOnt_NewTransferFromTransactionV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	toAddr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	mutableTransaction, err := testOntSdk.Native.Ont.NewTransferFromTransactionV2(testGasPrice, testGasLimit, testDefAcc.Address, testDefAcc.Address, toAddr, new(big.Int).SetInt64(1111111111125))
	assert.Nil(t, err)
	ontTx, err := mutableTransaction.IntoImmutable()
	assert.Nil(t, err)
	res, err := ParseNativeTxPayloadV2(ontTx.ToArray())
	assert.Nil(t, err)
	t.Logf("res:%v", res)
}

func TestOng_NewTransferFromTransactionV2(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testOntSdk.OpenWallet("./wallet.dat")
	testDefAcc, err := testWallet.GetDefaultAccount(testPasswd)
	assert.Nil(t, err)
	toAddr, err := common.AddressFromBase58("AWRBh9yYVzYHAfAb3tuWtdKjwGxNubimPo")
	assert.Nil(t, err)
	mutableTransaction, err := testOntSdk.Native.Ong.NewTransferFromTransactionV2(testGasPrice, testGasLimit, testDefAcc.Address, testDefAcc.Address, toAddr, new(big.Int).SetInt64(100000000000000008))
	assert.Nil(t, err)
	ontTx, err := mutableTransaction.IntoImmutable()
	assert.Nil(t, err)
	t.Logf("rawHex:%s", hex.EncodeToString(ontTx.Raw))
	res, err := ParseNativeTxPayloadV2(ontTx.ToArray())
	assert.Nil(t, err)
	t.Logf("res:%v", res)
}

// send ong tx test
func TestOngLegdgerV2_Transfer(t *testing.T) {
	testOntSdk = NewOntologySdk()
	testOntSdk.NewRpcClient().SetAddress(testNetUrl)
	txData := "00d1526aebfac409000000000000204e00000000000035eac618f06ed688813467cdfa1b3ed1aedf76397b00c66b1435eac618f06ed688813467cdfa1b3ed1aedf76396a7cc8141243502e7243c889181add5bc16a32bbc1e068356a7cc8080080e03779c311006a7cc86c51c10a7472616e7366657256321400000000000000000000000000000000000000020068164f6e746f6c6f67792e4e61746976652e496e766f6b650000"
	mutableTx, err := testOntSdk.GetMutableTx(txData)
	/*
		res,err := json.MarshalIndent(mutableTx,""," ")
		assert.Nil(t, err)
		t.Logf("%v",string(res))
	*/
	t.Logf("tx type:%d", mutableTx.TxType)
	t.Logf("tx payer:%s", mutableTx.Payer.ToBase58())
	assert.Nil(t, err)
	derSigData := "3045022100943f0841e849a0ca12e46b3d73469e45caa2c16517f41758bf3563c9f085412302200282de53f8384bc0db6926c523a3471a19fbfe08ac5256ad9adb5fab837cfe1c"
	r, s, err := ParseDerSig(derSigData)
	assert.Nil(t, err)
	sig := &signature.Signature{
		Scheme: signature.SHA256withECDSA,
		Value: &signature.DSASignature{
			R:     r,
			S:     s,
			Curve: elliptic.P256(),
		},
	}
	sigData, err := signature.Serialize(sig)
	assert.Nil(t, err)
	pub := "04c7f99c0f1cd9d37bec54c612d14de2474bf65e41e9724e77f3809a68d92c75f9513033e2813f5c1f5ea31d3fac7446a50ec8939ec9ea717cf3080fa29177e59a"
	pubKey, err := hex.DecodeString(pub)
	assert.Nil(t, err)
	pk, err := keypair.DeserializePublicKey(pubKey)
	assert.Nil(t, err)
	pkBytes := keypair.SerializePublicKey(pk)
	t.Logf("pub:%s", pub)
	t.Logf("pkBytes data:%s", hex.EncodeToString(pkBytes))
	//sigBytes,err := hex.DecodeString("")
	mutableTx.Sigs = append(mutableTx.Sigs, types.Sig{
		SigData: [][]byte{sigData},
		PubKeys: []keypair.PublicKey{pk},
		M:       1,
	})
	sendTx, err := testOntSdk.GetTxData(mutableTx)
	assert.Nil(t, err)
	t.Logf("sendTx:%s", sendTx)
	txHash, err := testOntSdk.SendTransaction(mutableTx)
	assert.Nil(t, err)
	t.Logf("hash:%v", txHash.ToHexString())
}

func ParseDerSig(derSigData string) (*big.Int, *big.Int, error) {
	derSig, _ := hex.DecodeString(derSigData)
	r := new(big.Int)
	s := new(big.Int)
	offset := 2
	if derSig[offset] != 0x02 {
		return nil, nil, fmt.Errorf("Invalid R integer marker")
	}
	offset++
	rLen := int(derSig[offset])
	offset++
	r.SetBytes(derSig[offset : offset+rLen])
	offset += rLen
	if derSig[offset] != 0x02 {
		return nil, nil, fmt.Errorf("Invalid S integer marker")
	}
	offset++
	sLen := int(derSig[offset])
	offset++
	s.SetBytes(derSig[offset : offset+sLen])
	return r, s, nil
}

func TestPubkeyToAddr(t *testing.T) {
	pub := "0253719ac66d7cafa1fe49a64f73bd864a346da92d908c19577a003a8a4160b7fa"
	pubKey, err := hex.DecodeString(pub)
	assert.Nil(t, err)
	pk, err := keypair.DeserializePublicKey(pubKey)
	assert.Nil(t, err)
	address := types.AddressFromPubKey(pk)
	t.Logf("address:%s", address.ToBase58())
}
