// Package state is the core API for the blockchain and implements all the
// business rules and processing.
package state

import (
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/genesis"
	"github.com/ardanlabs/blockchain/foundation/blockchain/mempool"
	"github.com/ardanlabs/blockchain/foundation/blockchain/storage/disk"
)

// =============================================================================

// EventHandler defines a function that is called when events
// occur in the processing of persisting blocks.
type EventHandler func(v string, args ...any)

// =============================================================================

// Config represents the configuration required to start
// the blockchain node.
type Config struct {
	BeneficiaryID database.AccountID
	Host          string
	DBPath        string
	EvHandler     EventHandler
}

// State manages the blockchain database.
type State struct {
	mu sync.RWMutex

	beneficiaryID database.AccountID
	host          string
	dbPath        string
	evHandler     EventHandler

	genesis genesis.Genesis
	mempool *mempool.Mempool
	db      *database.Database
}

// New constructs a new blockchain for data management.
func New(cfg Config) (*State, error) {

	// Build a safe event handler function for use.
	ev := func(v string, args ...any) {
		if cfg.EvHandler != nil {
			cfg.EvHandler(v, args...)
		}
	}

	// Load the genesis file to get starting balances for
	// founders of the block chain.
	genesis, err := genesis.Load()
	if err != nil {
		return nil, err
	}

	// Access the storage for the blockchain.
	storage, err := disk.New(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	// Access the storage for the blockchain.
	db, err := database.New(genesis, storage, ev)
	if err != nil {
		return nil, err
	}

	// Construct a mempool.
	mempool, err := mempool.New()
	if err != nil {
		return nil, err
	}

	// Create the State to provide support for managing the blockchain.
	state := State{
		beneficiaryID: cfg.BeneficiaryID,
		host:          cfg.Host,
		dbPath:        cfg.DBPath,
		evHandler:     ev,

		genesis: genesis,
		mempool: mempool,
		db:      db,
	}

	return &state, nil

}
