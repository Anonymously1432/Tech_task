package helper

import "sync"

type WalletLocker struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewWalletLocker() *WalletLocker {
	return &WalletLocker{
		locks: make(map[string]*sync.Mutex),
	}
}

// Получить мьютекс для конкретного кошелька
func (wl *WalletLocker) Lock(walletID string) {
	wl.mu.Lock()
	m, exists := wl.locks[walletID]
	if !exists {
		m = &sync.Mutex{}
		wl.locks[walletID] = m
	}
	wl.mu.Unlock()

	m.Lock() // блокируем конкретный кошелёк
}

func (wl *WalletLocker) Unlock(walletID string) {
	wl.mu.Lock()
	m := wl.locks[walletID]
	wl.mu.Unlock()

	m.Unlock()
}
