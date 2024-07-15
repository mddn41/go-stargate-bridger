package core

import "math/big"

// Common
var PrivateKeysFilePath = "data/private_keys.txt"
var DatabaseFilePath = "data/database.json"

// Dapps
var FullBridgeGasMultiplier *big.Float = big.NewFloat(1.5)
