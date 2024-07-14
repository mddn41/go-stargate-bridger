package core

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gleich/logoru"
)

var eip1559rewardPercentiles [5]float64 = [5]float64{10, 30, 50, 70, 90}

type EIP1559Params struct {
	maxPriorityFeePerGas big.Int
	maxFeePerGas         big.Int
}

type TxParams struct {
	To    *common.Address
	Data  []byte
	Value *big.Int
}

type EvmClient struct {
	Provider   *ethclient.Client
	privateKey *ecdsa.PrivateKey
	Address    common.Address
	Chain      Chain
}

func NewClient(privateKey string, chain Chain) (*EvmClient, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	provider, err := ethclient.Dial(chain.rpc)

	if err != nil {
		return nil, err
	}

	pk, err := crypto.HexToECDSA(privateKey)

	if err != nil {
		return nil, err
	}
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		return nil, err
	}
	return &EvmClient{privateKey: pk, Provider: provider, Address: crypto.PubkeyToAddress(*publicKeyECDSA), Chain: chain}, nil
}

func (c *EvmClient) EstimateTransaction(params *TxParams) (uint64, error) {
	gas, err := c.Provider.EstimateGas(context.Background(), ethereum.CallMsg{
		To:    params.To,
		Data:  params.Data,
		Value: params.Value,
	})

	if err != nil {
		return 0, err
	}
	return gas, nil
}

func (c *EvmClient) BuildTransaction(params *TxParams) (*types.Transaction, error) {
	nonce, err := c.Provider.PendingNonceAt(context.Background(), c.Address)
	if err != nil {
		return nil, err
	}

	gasLimit, err := c.EstimateTransaction(params)

	if err != nil {
		return nil, err
	}

	if c.Chain.eip1559 {
		eip1559Params := getEIP1559Params(c.Provider)

		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   c.Chain.ChainId,
			Nonce:     nonce,
			Gas:       gasLimit,
			GasTipCap: &eip1559Params.maxPriorityFeePerGas,
			GasFeeCap: &eip1559Params.maxFeePerGas,
			To:        params.To,
			Value:     params.Value,
			Data:      params.Data,
		}), nil
	}

	gasPrice, err := c.Provider.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		To:       params.To,
		Value:    params.Value,
		Data:     params.Data,
	}), nil
}

func (c *EvmClient) SendTransaction(params *TxParams) (bool, error) {
	tx, err := c.BuildTransaction(params)
	if err != nil {
		return false, fmt.Errorf("failed to build transaction: %s", err)
	}

	var signer types.Signer

	switch tx.Type() {
	case types.LegacyTxType:
		signer = types.NewEIP155Signer(c.Chain.ChainId)
	case types.DynamicFeeTxType:
		signer = types.NewLondonSigner(c.Chain.ChainId)
	}

	signedTx, err := types.SignTx(tx, signer, c.privateKey)
	if err != nil {
		return false, fmt.Errorf("failed to sign transaction: %s", err)
	}

	err = c.Provider.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return false, fmt.Errorf("failed to send transaction: %v", err)
	}

	return c.VerifyTransaction(signedTx.Hash()), nil
}

func (c *EvmClient) VerifyTransaction(txHash common.Hash) bool {
	receipt, err := waitForTxReceipt(c.Provider, txHash, 600, 0.1)

	if err != nil {
		logoru.Error(fmt.Sprintf("Failed to get transaction status: %s", err))
		return false
	}

	if receipt == nil {
		return false
	}

	if receipt.Status == 1 {
		logoru.Success(fmt.Sprintf("Transaction was successfull: %stx/%s", c.Chain.explorer, txHash.Hex()))
		return true
	}
	logoru.Error(fmt.Sprintf("Transaction failed: %s", txHash.Hex()))
	return false
}

func (c *EvmClient) GetNativeBalance() (*big.Int, error) {
	bal, err := c.Provider.BalanceAt(context.Background(), c.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get native balance: %s", err)
	}
	return bal, nil
}

func getEIP1559Params(client *ethclient.Client) *EIP1559Params {
	feeHistory, _ := client.FeeHistory(context.Background(), 1, nil, eip1559rewardPercentiles[:])

	var flatArray []*big.Int
	for _, row := range feeHistory.Reward {
		flatArray = append(flatArray, row...)
	}

	var rewardSum *big.Int = big.NewInt(0)
	for _, item := range flatArray {
		rewardSum.Add(item, rewardSum)
	}

	avgReward := rewardSum.Div(rewardSum, big.NewInt(int64(len(flatArray))))

	//https://github.com/ethereum/go-ethereum/issues/26542
	lastBlockHeader, _ := client.HeaderByNumber(context.Background(), nil)
	nextBaseFee := lastBlockHeader.BaseFee

	return &EIP1559Params{
		maxPriorityFeePerGas: *avgReward,
		maxFeePerGas:         *big.NewInt(0).Add(nextBaseFee, avgReward),
	}

}

func waitForTxReceipt(client *ethclient.Client, txHash common.Hash, timeout int, pollLatency float64) (*types.Receipt, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	for {
		tx, err := client.TransactionReceipt(ctxTimeout, txHash)

		if errors.Is(err, context.DeadlineExceeded) {
			logoru.Error(fmt.Sprintf("Transaction %s was not found in the blockchain after %d seconds", txHash.Hex(), timeout))
		}

		if err != nil && err != ethereum.NotFound {
			return nil, err
		}
		if tx != nil {
			return tx, nil
		}
		time.Sleep(time.Duration(pollLatency) * time.Second)
	}
}
