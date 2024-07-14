package main

import (
	"fmt"

	"github.com/mddn41/go-stargate-bridger/core"
	"github.com/mddn41/go-stargate-bridger/dapps"
)

func main() {
	pk := ""
	c, _ := core.NewClient(pk, core.LineaChain)

	s := dapps.Stargate(c)
	fmt.Println(s.Bridge(core.ScrollChain, nil, dapps.StargateTaxiBridgeMode, false))
}
