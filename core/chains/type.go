package chains

import "math/big"

type Chain struct {
	Name       string
	ChainId    *big.Int
	CoinSymbol string
	Explorer   string
	ERIP1559   bool
	RPC        string
	LzEid      int
}

func ChainByName(name string) *Chain {
	switch name {
	case ArbitrumChain.Name:
		return ArbitrumChain
	case OptimismChain.Name:
		return OptimismChain
	case BaseChain.Name:
		return BaseChain
	case LineaChain.Name:
		return LineaChain
	case ScrollChain.Name:
		return ScrollChain
	default:
		return nil
	}
}
