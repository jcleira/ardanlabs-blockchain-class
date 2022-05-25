package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/merkle"
	"github.com/ardanlabs/blockchain/foundation/blockchain/signature"
	"github.com/ardanlabs/blockchain/foundation/blockchain/storage/disk"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// https://goethereumbook.org/signature-verify/

func main() {
	err := readBlock()

	if err != nil {
		log.Fatalln(err)
	}
}

func readBlock() error {
	d, err := disk.New("zblock/miner1")
	if err != nil {
		return err
	}

	blockData, err := d.GetBlock(1)
	if err != nil {
		return err
	}

	fmt.Println(blockData)

	block, err := database.ToBlock(blockData)
	if err != nil {
		return err
	}

	if blockData.Header.TransRoot != block.MerkleTree.RootHex() {
		return errors.New("merkle tree wrong")
	}

	fmt.Println("merkle tree matches")

	return nil
}

func writeBlock() error {
	txs := []database.Tx{
		{
			ChainID: 1,
			Nonce:   1,
			ToID:    "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
			Value:   100,
			Tip:     50,
		},
		{
			ChainID: 1,
			Nonce:   2,
			ToID:    "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
			Value:   100,
			Tip:     50,
		},
	}

	blockTxs := make([]database.BlockTx, len(txs))
	for i, tx := range txs {
		blockTx, err := signToBlockTx(tx, 15)
		if err != nil {
			return err
		}

		blockTxs[i] = blockTx
	}

	tree, err := merkle.NewTree(blockTxs)
	if err != nil {
		return err
	}

	beneficiaryID, err := database.ToAccountID("0xF01813E4B85e178A83e29B8E7bF26BD830a25f32")
	if err != nil {
		return err
	}

	block := database.Block{
		Header: database.BlockHeader{
			Number:        1,
			PrevBlockHash: signature.ZeroHash,
			TimeStamp:     uint64(time.Now().UTC().Unix()),
			BeneficiaryID: beneficiaryID,
			Difficulty:    6,
			MiningReward:  700,
			StateRoot:     "not defined",
			TransRoot:     tree.RootHex(), //
			Nonce:         0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	bd := database.NewBlockData(block)

	d, err := disk.New("zblock/miner1")
	if err != nil {
		return err
	}

	if err := d.Write(bd); err != nil {
		return err
	}

	return nil
}

func signToBlockTx(tx database.Tx, gas uint64) (database.BlockTx, error) {
	pk, err := crypto.HexToECDSA("fae85851bdf5c9f49923722ce38f3c1defcfd3619ef5453230a58ad805499959")
	if err != nil {
		return database.BlockTx{}, err
	}

	signedTx, err := tx.Sign(pk)
	if err != nil {
		return database.BlockTx{}, err
	}

	return database.NewBlockTx(signedTx, gas, 1), nil
}

func sign() error {

	// Load the private key from disk.
	path := fmt.Sprintf("%s%s.ecdsa", "zblock/accounts/", "kennedy")
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	// Display the address of the account.
	address := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	fmt.Println("Address:", address)

	// =========================================================================
	// Stamp and sign the data.

	v := struct {
		Name string
		Age  int
	}{
		Name: "Bill",
		Age:  10,
	}

	data, err := stamp(v)
	if err != nil {
		return fmt.Errorf("stamp: %w", err)
	}

	// Sign the hash with the private key to produce a signature.
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	fmt.Printf("SIG: %s\n", hexutil.Encode(sig))

	// =========================================================================
	// Pretent we are the node receiving the data and the signature. We need
	// the address of the person who signed this.

	// Capture the public key associated with this data and signature.
	publicKey, err := crypto.SigToPub(data, sig)
	if err != nil {
		return fmt.Errorf("sigToPub: %w", err)
	}

	// Extract the account address from the public key.
	fmt.Println("Address:", crypto.PubkeyToAddress(*publicKey).String())

	// =========================================================================
	// If we want to validate the data and signature, we need access to the
	// original public key. We can't use the public key returned by SigToPub
	// since that function is calculating the public key based on the data and
	// signature.

	// Extract the bytes for the original public key.
	publicKeyOrg := privateKey.Public()
	publicKeyECDSA, ok := publicKeyOrg.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// Check the public key validates the data and signature.
	rs := sig[:crypto.RecoveryIDOffset]
	if !crypto.VerifySignature(publicKeyBytes, data, rs) {
		return errors.New("invalid signature")
	}

	fmt.Println("Signature Validated")

	return nil
}

func stamp(value any) ([]byte, error) {

	// Marshal the data.
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// Hash the transaction data into a 32 byte array. This will provide
	// a data length consistency with all transactions.
	txHash := crypto.Keccak256Hash(data)

	// Convert the stamp into a slice of bytes. This stamp is
	// used so signatures we produce when signing transactions
	// are always unique to the Ardan blockchain.
	stamp := []byte("\x19Ardan Signed Message:\n32")

	// Hash the stamp and txHash together in a final 32 byte array
	// that represents the transaction data.
	tran := crypto.Keccak256Hash(stamp, txHash.Bytes())

	return tran.Bytes(), nil
}
