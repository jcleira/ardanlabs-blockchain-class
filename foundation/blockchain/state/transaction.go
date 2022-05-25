package state

import "github.com/ardanlabs/blockchain/foundation/blockchain/database"

// UpsertWalletTransaction accepts a transaction from a wallet for inclusion.
func (s *State) UpsertWalletTransaction(signedTx database.SignedTx) error {

	// CORE NOTE: Just check the signed transaction has a proper signature and
	// valid account for the recipient. It's up to the wallet to make sure the
	// account has a proper balance and nonce. Fees will be taken if this
	// transaction is mined into a block and those types of validation fail.

	if err := signedTx.Validate(); err != nil {
		return err
	}

	const oneUnitOfGas = 1
	tx := database.NewBlockTx(signedTx, s.genesis.GasPrice, oneUnitOfGas)
	if err := s.mempool.Upsert(tx); err != nil {
		return err
	}

	return nil
}
