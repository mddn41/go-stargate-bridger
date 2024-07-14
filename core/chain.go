package core

import "math/big"

type Chain struct {
	Name       string
	ChainId    *big.Int
	coinSymbol string
	explorer   string
	eip1559    bool
	rpc        string
	LzEid      int
}

var MainnetChain = Chain{
	Name:       "ERC20",
	ChainId:    big.NewInt(1),
	coinSymbol: "ETH",
	explorer:   "https://etherscan.io/",
	eip1559:    true,
	rpc:        "https://rpc.ankr.com/eth",
}

var ArbitrumChain = Chain{
	Name:       "Arbitrum",
	ChainId:    big.NewInt(42161),
	coinSymbol: "ETH",
	explorer:   "https://arbiscan.io/",
	eip1559:    true,
	rpc:        "https://rpc.ankr.com/arbitrum",
	LzEid:      30110,
}

var OptimismChain = Chain{
	Name:       "Optimism",
	ChainId:    big.NewInt(10),
	coinSymbol: "ETH",
	explorer:   "https://optimistic.etherscan.io/",
	eip1559:    true,
	rpc:        "https://rpc.ankr.com/optimism",
	LzEid:      30111,
}

var BaseChain = Chain{
	Name:       "Base",
	ChainId:    big.NewInt(8453),
	coinSymbol: "ETH",
	explorer:   "https://basescan.org/",
	eip1559:    true,
	rpc:        "https://rpc.ankr.com/base",
	LzEid:      30184,
}

var LineaChain = Chain{
	Name:       "Linea",
	ChainId:    big.NewInt(59144),
	coinSymbol: "ETH",
	explorer:   "https://lineascan.build/",
	eip1559:    true,
	rpc:        "https://rpc.linea.build/",
	LzEid:      30183,
}

var ScrollChain = Chain{
	Name:       "Scroll",
	ChainId:    big.NewInt(534352),
	coinSymbol: "ETH",
	explorer:   "https://scrollscan.com/",
	eip1559:    true,
	rpc:        "https://rpc.ankr.com/scroll",
	LzEid:      30214,
}
