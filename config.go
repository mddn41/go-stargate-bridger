package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type GeneralSettings struct {
	WalletDelayRange          [2]int
	AfterFailDelayRange       [2]int
	SrcChain                  string
	DstChain                  string
	StargateBridgeMode        string
	BalancePercentageToBridge [2]int
	IncludeFees               bool
	UseFullBridge             bool
}

type Config struct {
	GeneralSettings GeneralSettings
	RpcEndpoints    map[string]string
}

var UserConfig GeneralSettings

func init() {
	fmt.Println("initting")
	fmt.Printf("\n%v", UserConfig)
	_, err := toml.DecodeFile("config.toml", &UserConfig)

	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%v", UserConfig)
}
