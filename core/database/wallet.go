package database

import (
	"github.com/mddn41/go-stargate-bridger/core"
	"github.com/mddn41/go-stargate-bridger/core/chains"
)

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
	SrcChain   string `json:"srcChain"`
	DstChain   string `json:"dstChain"`
	BridgeSent bool   `json:"bridgeSent"`
}

func (w *Wallet) ToClient(chain *chains.Chain) (*core.EvmClient, error) {
	return core.NewClient(w.PrivateKey, chain)
}
