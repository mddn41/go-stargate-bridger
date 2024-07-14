package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func AddressToBytes32(address common.Address) [32]byte {
	var bytes32 [32]byte
	copy(bytes32[12:], address[:])
	return bytes32
}

func WeiToEther(wei *big.Int) *big.Float {
	decimals := new(big.Float).SetFloat64(1e18)
	etherValue := new(big.Float).Quo(new(big.Float).SetInt(wei), decimals)
	return etherValue
}

func ApplySlippage(value *big.Int, slippage float64) *big.Int {
	valueFloat := new(big.Float).SetInt(value)
	multiplierFloat := new(big.Float).SetFloat64(1 - slippage)

	resultFloat := new(big.Float).Mul(valueFloat, multiplierFloat)

	resultInt := new(big.Int)
	resultFloat.Int(resultInt)
	return resultInt
}
