package dapps

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gleich/logoru"

	"github.com/mddn41/go-stargate-bridger/core"
)

var (
	StargateBusBridgeMode  = []byte{0x01}
	StargateTaxiBridgeMode = []byte{}
)

type StargateSendParam struct {
	DstEid       uint32
	To           [32]byte
	AmountLD     *big.Int
	MinAmountLD  *big.Int
	ExtraOptions []byte
	ComposeMsg   []byte
	OftCmd       []byte
}

type StargateMessagingFee struct {
	NativeFee  *big.Int "json:\"nativeFee\""
	LzTokenFee *big.Int "json:\"lzTokenFee\""
}

var StargateQuoteResult struct {
	Fee StargateMessagingFee
}

var stargateContractAddresses = map[*big.Int]common.Address{
	core.ArbitrumChain.ChainId: common.HexToAddress("0xA45B5130f36CDcA45667738e2a258AB09f4A5f7F"),
	core.OptimismChain.ChainId: common.HexToAddress("0xe8CDF27AcD73a434D661C84887215F7598e7d0d3"),
	core.BaseChain.ChainId:     common.HexToAddress("0xdc181Bd607330aeeBEF6ea62e03e5e1Fb4B6F7C7"),
	core.LineaChain.ChainId:    common.HexToAddress("0x81F6138153d473E8c5EcebD3DC8Cd4903506B075"),
	core.ScrollChain.ChainId:   common.HexToAddress("0xC2b638Cb5042c1B3c5d5C969361fB50569840583"),
}

type StargateBridge struct {
	client          *core.EvmClient
	contractAddress common.Address
	abi             abi.ABI
}

func Stargate(client *core.EvmClient) *StargateBridge {
	reader, _ := os.Open("dapps/abi/StargateNativePool.json")
	abi, _ := abi.JSON(reader)

	return &StargateBridge{
		client:          client,
		abi:             abi,
		contractAddress: stargateContractAddresses[client.Chain.ChainId],
	}
}

func (dapp *StargateBridge) getMessageFee(dstChainLzId int, amount *big.Int, mode []byte) (*StargateMessagingFee, error) {
	data, _ := dapp.abi.Pack("quoteSend", &StargateSendParam{
		uint32(dstChainLzId),
		core.AddressToBytes32(dapp.client.Address),
		amount,
		core.ApplySlippage(amount, 0.005),
		[]byte{},
		[]byte{},
		mode,
	}, false)
	res, err := dapp.client.Provider.CallContract(context.Background(), ethereum.CallMsg{
		To:   &dapp.contractAddress,
		Data: data,
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get message fee: %s", err)
	}

	dapp.abi.UnpackIntoInterface(&StargateQuoteResult, "quoteSend", res)
	return &StargateQuoteResult.Fee, nil
}

func (dapp *StargateBridge) EstimateAmountBeforeFees(dstChainLzId int, amount *big.Int, mode []byte) (*big.Int, error) {
	var simulationValue *big.Int

	if amount == nil {
		balance, err := dapp.client.GetNativeBalance()
		if err != nil {
			return nil, fmt.Errorf("failed to estimate amount before fees: %s", err)
		}
		amount = balance
		simulationValue = big.NewInt(200000000000000)
	} else {
		simulationValue = amount
	}

	messageFee, err := dapp.getMessageFee(dstChainLzId, amount, mode)

	if err != nil {
		return nil, fmt.Errorf("failed to estimate amount before fees: %s", err)
	}

	sendParam := &StargateSendParam{
		uint32(dstChainLzId),
		core.AddressToBytes32(dapp.client.Address),
		simulationValue,
		core.ApplySlippage(simulationValue, 0.005),
		[]byte{},
		[]byte{},
		mode,
	}

	txData, _ := dapp.abi.Pack("send", sendParam, messageFee, dapp.client.Address)

	txValue := big.NewInt(0).Add(simulationValue, messageFee.NativeFee)
	gas, err := dapp.client.EstimateTransaction(&core.TxParams{To: &dapp.contractAddress, Data: txData, Value: txValue})

	if err != nil {
		return nil, fmt.Errorf("failed to estimate amount before fees: %s", err)
	}

	gasPrice, _ := dapp.client.Provider.SuggestGasPrice(context.Background())

	txGasFees := big.NewInt(0).Mul(gasPrice, big.NewInt(int64(gas)))

	new(big.Float).Mul(
		new(big.Float).SetInt(txGasFees),
		core.FullBridgeGasMultiplier,
	).Int(txGasFees)

	finalAmount := big.NewInt(0).Sub(amount, messageFee.NativeFee)
	return finalAmount.Sub(finalAmount, txGasFees), nil
}

func (dapp *StargateBridge) Bridge(dstChain core.Chain, amount *big.Int, mode []byte, includeFees bool) bool {
	if amount == nil || includeFees {
		estimatedAmount, err := dapp.EstimateAmountBeforeFees(dstChain.LzEid, amount, mode)
		if err != nil {
			logoru.Error(fmt.Sprintf("Stargate Error: %s", err))
			return false
		}
		amount = estimatedAmount
	}

	logoru.Info(fmt.Sprintf("Stargate: Bridging %f ETH from %s to %s", core.WeiToEther(amount), dapp.client.Chain.Name, dstChain.Name))

	messageFee, err := dapp.getMessageFee(dstChain.LzEid, amount, mode)

	if err != nil {
		logoru.Error(fmt.Sprintf("Stargate Error: %s", err))
		return false
	}

	sendParam := &StargateSendParam{
		uint32(dstChain.LzEid),
		core.AddressToBytes32(dapp.client.Address),
		amount,
		core.ApplySlippage(amount, 0.005),
		[]byte{},
		[]byte{},
		mode,
	}

	txData, _ := dapp.abi.Pack("send", sendParam, messageFee, dapp.client.Address)
	txValue := big.NewInt(0).Add(amount, messageFee.NativeFee)

	_, err = dapp.client.SendTransaction(&core.TxParams{To: &dapp.contractAddress, Data: txData, Value: txValue})
	if err != nil {
		logoru.Error(fmt.Sprintf("Stargate Error: %v", err))
		return false
	}
	return true
}
