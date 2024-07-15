package core

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand/v2"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/schollz/progressbar/v3"
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

func ReadLinesFromTxt(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func Sleep(delayRange [2]int) {
	sleepTime := rand.IntN(delayRange[1]-delayRange[0]) + delayRange[0]
	bar := progressbar.Default(int64(sleepTime))

	for i := 0; i < sleepTime; i++ {
		bar.Add(1)
		time.Sleep(1 * time.Second)
	}
}

func PrintGreeting() {
	fmt.Print(`
 ____ ____ ____ ____ ____ ____ ____ ____       
||S |||t |||a |||r |||g |||a |||t |||e ||      
||__|||__|||__|||__|||__|||__|||__|||__||      
|/__\|/__\|/__\|/__\|/__\|/__\|/__\|/__\|      
 _________  ____ ____ ____ ____ ____ ____ ____ 
||       ||||B |||r |||i |||d |||g |||e |||r ||
||_______||||__|||__|||__|||__|||__|||__|||__||
|/_______\||/__\|/__\|/__\|/__\|/__\|/__\|/__\|

Choose module:
1. Create database
2. Bridge batch
`)
}
