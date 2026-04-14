package main

import (
	"fmt"
	"log"

	"github.com/tyler-smith/go-bip39"
)

func main() {
	entropy, err := bip39.NewEntropy(256) // 24 words
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(mnemonic)
}
