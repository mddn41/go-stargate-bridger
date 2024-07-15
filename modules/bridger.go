package modules

import (
	"fmt"
	"math/big"
	"math/rand/v2"

	"github.com/gleich/logoru"
	"github.com/mddn41/go-stargate-bridger/config"
	"github.com/mddn41/go-stargate-bridger/core"
	"github.com/mddn41/go-stargate-bridger/core/chains"
	"github.com/mddn41/go-stargate-bridger/core/database"
	"github.com/mddn41/go-stargate-bridger/dapps"
)

func BridgeBatch() {
	db, err := database.LoadDatabase()

	if err != nil {
		logoru.Error(err)
	}

	var bridgeMode []byte

	if config.UserConfig.StargateBridgeMode == "BUS" {
		bridgeMode = dapps.StargateBusBridgeMode
	} else {
		bridgeMode = dapps.StargateTaxiBridgeMode
	}

	for wallet := db.GetRandomItem(); wallet != nil; wallet = db.GetRandomItem() {
		client, err := wallet.ToClient(chains.ChainByName(wallet.SrcChain))

		if err != nil {
			logoru.Error(fmt.Sprintf("Error while initializing client: %s", err))
			core.Sleep(config.UserConfig.AfterFailDelayRange)
			continue
		}

		logoru.Info(fmt.Sprintf("Wallet: %s", wallet.Address))
		logoru.Debug(fmt.Sprintf("Transactions left: %d", db.WalletsLeft()))

		var amountToBridge *big.Int
		var includeFees bool

		if !config.UserConfig.UseFullBridge {
			includeFees = config.UserConfig.IncludeFees

			userBalance, err := client.GetNativeBalance()

			if err != nil {
				logoru.Error(err)
			}

			percentToBridge := rand.IntN(config.UserConfig.BalancePercentageToBridge[1]-config.UserConfig.BalancePercentageToBridge[0]) + config.UserConfig.BalancePercentageToBridge[0]
			println(percentToBridge)

			amountToBridge, _ = new(big.Float).Mul(
				new(big.Float).SetInt(userBalance),
				big.NewFloat(float64(percentToBridge)/float64(100)),
			).Int(nil)
			fmt.Println(amountToBridge)
		}

		stargate := dapps.Stargate(client)

		var sleepRange [2]int

		if stargate.Bridge(*chains.ChainByName(config.UserConfig.DstChain), amountToBridge, bridgeMode, includeFees) {
			wallet.BridgeSent = true
			db.Save()
			sleepRange = config.UserConfig.WalletDelayRange
		} else {
			sleepRange = config.UserConfig.AfterFailDelayRange
		}

		core.Sleep(sleepRange)
	}

	logoru.Success("No more active wallets left")
}
