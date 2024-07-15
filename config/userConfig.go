package config

import (
	"github.com/BurntSushi/toml"
)

type generalSettings struct {
	WalletDelayRange          [2]int
	AfterFailDelayRange       [2]int
	SrcChain                  string
	DstChain                  string
	StargateBridgeMode        string
	BalancePercentageToBridge [2]int
	IncludeFees               bool
	UseFullBridge             bool
}

type config struct {
	GeneralSettings generalSettings
	RpcEndpoints    map[string]string
}

var UserConfig generalSettings

func init() {
	config := &config{}
	_, err := toml.DecodeFile("config.toml", &config)

	if err != nil {
		panic(err)
	}
	UserConfig = config.GeneralSettings
}
