package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
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

	v := struct {
		Name string
	}{
		Name: "Bill",
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Hash the transaction data into a 32 byte array. This will provide
	// a data length consistency with all transactions.
	txHash := crypto.Keccak256(data)

	// Sign the hash with the private key to produce a signature.
	sig, err := crypto.Sign(txHash, privateKey)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	fmt.Println("SIG:", sig)

	// =========================================================================

	v2 := struct {
		Name string
	}{
		Name: "Bill",
	}

	data2, err := json.Marshal(v2)
	if err != nil {
		return err
	}

	// Hash the transaction data into a 32 byte array. This will provide
	// a data length consistency with all transactions.
	txHash2 := crypto.Keccak256(data2)

	sigPublicKey, err := crypto.Ecrecover(txHash2, sig)
	if err != nil {
		return err
	}

	rs := sig[:crypto.RecoveryIDOffset]
	if !crypto.VerifySignature(sig, txHash2, rs) {
		return errors.New("invalid signature")
	}

	// Capture the public key associated with this signature.
	x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	publicKey := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	// Extract the account address from the public key.
	address = crypto.PubkeyToAddress(publicKey).String()
	fmt.Println(address)

	return nil
}
