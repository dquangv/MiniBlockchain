package blockchain

import "sync"

var (
	PendingTxs   []*Transaction
	pendingMutex sync.Mutex
)

func AddPendingTx(tx *Transaction) {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()
	PendingTxs = append(PendingTxs, tx)
}

func GetAndClearPendingTxs() []*Transaction {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()

	txs := PendingTxs
	PendingTxs = nil // hoặc []*Transaction{}
	return txs
}
