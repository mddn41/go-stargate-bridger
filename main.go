package main

import (
	"fmt"

	"github.com/gleich/logoru"
	"github.com/mddn41/go-stargate-bridger/core"
	"github.com/mddn41/go-stargate-bridger/core/database"
	"github.com/mddn41/go-stargate-bridger/modules"
)

//"fmt"

func main() {
	core.PrintGreeting()

	var moduleNum int
	_, err := fmt.Scan(&moduleNum)

	if err != nil {
		logoru.Error("Invalid module number")
		return
	}

	switch moduleNum {
	case 1:
		db := database.NewDatabase()

		if db == nil {
			return
		}

		db.Save()
		logoru.Success("Database has been successfully created")
	case 2:
		modules.BridgeBatch()
	default:
		core.Sleep([2]int{10, 20})
		logoru.Error("Invalid module number")
	}
}
