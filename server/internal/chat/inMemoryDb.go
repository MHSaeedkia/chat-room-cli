package chat

import "sync"

type inMemoryDb struct {
	storage map[uuId]Client
	mtx     sync.RWMutex
}
