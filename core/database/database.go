package database

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/gleich/logoru"
	"github.com/mddn41/go-stargate-bridger/config"
	"github.com/mddn41/go-stargate-bridger/core"
)

type JsonDatabase struct {
	Data []*Wallet
}

func NewDatabase() *JsonDatabase {
	if config.UserConfig.SrcChain == "" || config.UserConfig.DstChain == "" {
		logoru.Error("Database: SrcChain and DstChain must be specified")
		return nil
	}

	if config.UserConfig.SrcChain == config.UserConfig.DstChain {
		logoru.Error("Database: SrcChain and DstChain must have different values")
		return nil
	}

	privateKeys, err := core.ReadLinesFromTxt(core.PrivateKeysFilePath)

	if err != nil {
		logoru.Error(fmt.Sprintf("Database: failed to read private keys file: %s", err))
		return nil
	}

	var data []*Wallet

	for _, privateKey := range privateKeys {
		client, err := core.NewClient(privateKey, nil)

		if err != nil {
			logoru.Error(fmt.Sprintf("Database: failed to initialize client: %s", err))
			return nil
		}

		data = append(data, &Wallet{
			PrivateKey: privateKey,
			Address:    client.Address.Hex(),
			SrcChain:   config.UserConfig.SrcChain,
			DstChain:   config.UserConfig.DstChain,
		})
	}
	return &JsonDatabase{Data: data}
}

func LoadDatabase() (*JsonDatabase, error) {
	jsonData, err := os.ReadFile(core.DatabaseFilePath)
	if err != nil {
		return nil, fmt.Errorf("database: failed to read database file: %s", err)
	}

	var dbData []*Wallet
	err = json.Unmarshal(jsonData, &dbData)

	if err != nil {
		return nil, fmt.Errorf("database: failed to deserialize data: %s", err)
	}

	return &JsonDatabase{Data: dbData}, nil
}

func (db *JsonDatabase) GetRandomItem() *Wallet {
	var activeWallets = db.getActiveWallets()

	walletsLeft := len(activeWallets)

	if walletsLeft == 0 {
		return nil
	}

	randomIndex := rand.Intn(walletsLeft)
	return activeWallets[randomIndex]
}

func (db *JsonDatabase) Save() error {
	dataJson, err := json.MarshalIndent(db.Data, "", "\t")
	if err != nil {
		return fmt.Errorf("database: failed to serialize data: %s", err)
	}

	err = os.WriteFile(core.DatabaseFilePath, dataJson, 0644)
	if err != nil {
		return fmt.Errorf("database: failed to write data into file: %s", err)
	}
	return nil
}

func (db *JsonDatabase) WalletsLeft() int {
	return len(db.getActiveWallets())
}

func (db *JsonDatabase) getActiveWallets() []*Wallet {
	var activeWallets []*Wallet
	for _, w := range db.Data {
		if !w.BridgeSent {
			activeWallets = append(activeWallets, w)
		}
	}
	return activeWallets
}
