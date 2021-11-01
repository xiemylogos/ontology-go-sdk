package ontology_go_sdk

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
	"github.com/ontio/ontology/common"
)

type Address common.Address

func (this *Address) ToHexString() string {
	return fmt.Sprintf("%x", common.ToArrayReverse(this[:]))
}

// ToBase58 returns base58 encoded address string
func (f *Address) ToBase58() string {
	data := append([]byte{23}, f[:]...)
	temp := sha256.Sum256(data)
	temps := sha256.Sum256(temp[:])
	data = append(data, temps[0:4]...)

	bi := new(big.Int).SetBytes(data).String()
	encoded, _ := base58.BitcoinEncoding.Encode([]byte(bi))
	return string(encoded)
}
