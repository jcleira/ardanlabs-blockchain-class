package mempool

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
)

// Mempool represents a cache of transactions organized by account:nonce.
type Mempool struct {
	mu   sync.RWMutex
	pool map[string]database.BlockTx
}

// New constructs a new mempool using the default sort strategy.
func New() (*Mempool, error) {
	return &Mempool{
		pool: make(map[string]database.BlockTx),
	}, nil
}

// Count returns the current number of transaction in the pool.
func (mp *Mempool) Count() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return len(mp.pool)
}

// Upsert adds or replaces a transaction from the mempool.
func (mp *Mempool) Upsert(tx database.BlockTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// CORE NOTE: Different blockchains have different algorithms to limit the
	// size of the mempool. Some limit based on the amount of memory being
	// consumed and some may limit based on the number of transaction. If a limit
	// is met, then either the transaction that has the least return on investment
	// or the oldest will be dropped from the pool to make room for new the transaction.

	// For now, the Ardan blockchain in not imposing any limits.
	key, err := mapKey(tx)
	if err != nil {
		return err
	}

	// Ethereum requires a 10% bump in the tip to replace an existing
	// transaction in the mempool and so do we. We want to limit users
	// from this sort of behavior.
	if etx, exists := mp.pool[key]; exists {
		if tx.Tip < uint64(math.Round(float64(etx.Tip)*1.10)) {
			return errors.New("replacing a transaction requires a 10% bump in the tip")
		}
	}

	mp.pool[key] = tx

	return nil
}

// Delete removed a transaction from the mempool.
func (mp *Mempool) Delete(tx database.BlockTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	key, err := mapKey(tx)
	if err != nil {
		return err
	}

	delete(mp.pool, key)

	return nil
}

// PickBest uses the configured sort strategy to return a set of transactions.
// If 0 is passed, all transactions in the mempool will be returned.
func (mp *Mempool) PickBest() []database.BlockTx {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	var txs []database.BlockTx

	for _, tx := range mp.pool {
		txs = append(txs, tx)
	}

	return txs
}

// =============================================================================

// mapKey is used to generate the map key.
func mapKey(tx database.BlockTx) (string, error) {
	account, err := tx.FromAccount()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", account, tx.Nonce), nil
}
