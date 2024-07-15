# Stargate Bridger

Software for one-way bridges via Stargate using Bus/Taxi strategies.

## Configuring
1. Write down private keys into `data/private_keys.txt` file, one private key per line
2. Set up config in the `config.toml` file

## Usage
In the project's folder root execute the following command via preffered CL: `go run main.go`.

Database must be created before running main module. Database is used for storing wallets state, which allows user to stop soft execution almost any time and re-launch without repeating bridges on wallets, that have already performed bridge. User also can manually check or change (only when soft is not running) current database state by editing `data/database.json` file.

### ⚠️ Disclaimer
> This is a rewritten verison of the [Stargate Bridger](https://github.com/sybil-v-zakone/stargate-v2-bridger) originally written by me a few months ago using Python. This project is considered as first experince with Go and I am not responsible for any potential money loss. Use at your own risk and only if you know what you do. Any bugs, code improvements purposals and etc. are welcomed in the [Issues](https://github.com/mddn41/go-stargate-bridger/issues) section.
