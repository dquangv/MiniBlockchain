package blockchain

import "sync"

// PendingTxs stores transactions that haven't been included in a block yet.
// These will be processed by the leader node during the next block creation cycle.
var (
	PendingTxs   []*Transaction // Slice holding pending transactions
	pendingMutex sync.Mutex     // Mutex to ensure thread-safe access to PendingTxs
)

// AddPendingTx safely adds a transaction to the pending pool.
func AddPendingTx(tx *Transaction) {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()
	PendingTxs = append(PendingTxs, tx)
}

// GetAndClearPendingTxs retrieves all pending transactions and clears the pool.
// This is called by the leader when it's ready to create a new block.
// The returned list is used to construct the block's transaction set.
func GetAndClearPendingTxs() []*Transaction {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()

	txs := PendingTxs
	PendingTxs = nil
	return txs
}
