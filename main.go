package main

import (
	"os"

  "https://github.com/maphy9/btc-utxo-indexer/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
