package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println(address)

	data, err := json.Marshal(
		struct{ Name string }{Name: "Jose"},
	)

	sig, err := crypto.Sign(crypto.Keccak256(data), privateKey)
	if err != nil {
		return err
	}
	fmt.Println(sig)

	return nil
}
