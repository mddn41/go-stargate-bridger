package chains

import (
	"math/big"

	"github.com/mddn41/go-stargate-bridger/config"
)

var MainnetChain = &Chain{
	Name:       "ERC20",
	ChainId:    big.NewInt(1),
	CoinSymbol: "ETH",
	Explorer:   "https://etherscan.io/",
	ERIP1559:   true,
	RPC:        config.RpcEndpoints["ERC20"],
}

var ArbitrumChain = &Chain{
	Name:       "Arbitrum",
	ChainId:    big.NewInt(42161),
	CoinSymbol: "ETH",
	Explorer:   "https://arbiscan.io/",
	ERIP1559:   true,
	RPC:        config.RpcEndpoints["Arbitrum"],
	LzEid:      30110,
}

var OptimismChain = &Chain{
	Name:       "Optimism",
	ChainId:    big.NewInt(10),
	CoinSymbol: "ETH",
	Explorer:   "https://optimistic.etherscan.io/",
	ERIP1559:   true,
	RPC:        config.RpcEndpoints["Optimism"],
	LzEid:      30111,
}

var BaseChain = &Chain{
	Name:       "Base",
	ChainId:    big.NewInt(8453),
	CoinSymbol: "ETH",
	Explorer:   "https://basescan.org/",
	ERIP1559:   true,
	RPC:        config.RpcEndpoints["Base"],
	LzEid:      30184,
}

var LineaChain = &Chain{
	Name:       "Linea",
	ChainId:    big.NewInt(59144),
	CoinSymbol: "ETH",
	Explorer:   "https://lineascan.build/",
	ERIP1559:   true,
	RPC:        config.RpcEndpoints["Linea"],
	LzEid:      30183,
}

var ScrollChain = &Chain{
	Name:       "Scroll",
	ChainId:    big.NewInt(534352),
	CoinSymbol: "ETH",
	Explorer:   "https://scrollscan.com/",
	ERIP1559:   true,
	RPC:        config.RpcEndpoints["Scroll"],
	LzEid:      30214,
}
