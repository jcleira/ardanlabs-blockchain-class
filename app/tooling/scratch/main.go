package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := sign(); err != nil {
		log.Fatalln(err)
	}
}

func sign() error {
	path := fmt.Sprintf("%s%s.ecdsa", "zblock/accounts/", "kennedy")
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println(address)

	data, err := json.Marshal(
		struct{ Name string }{Name: "Bill"},
	)

	dataBytes := crypto.Keccak256(data)
	sig, err := crypto.Sign(dataBytes, privateKey)
	if err != nil {
		return err
	}
	fmt.Println(sig)

	sigPublicKey, err := crypto.Ecrecover(dataBytes, sig)
	if err != nil {
		return fmt.Errorf("ecrecover, %w", err)
	}
	x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	sigPublicKeyVal := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	recoveredAddress := crypto.PubkeyToAddress(sigPublicKeyVal).String()
	fmt.Println(recoveredAddress)

	return nil
}
