package keyvaluestore

import (
	"fmt"
	"sync"
)

type Store interface {
	StringForKey(key string) (string, error)
	SetString(key, value string) error
}

type KeyValueStore struct {
	stringStore     map[string]string
	stringStoreLock *sync.Mutex
}

func New() *KeyValueStore {
	return &KeyValueStore{stringStore: make(map[string]string), stringStoreLock: &sync.Mutex{}}
}

func (store *KeyValueStore) StringForKey(key string) (string, error) {
	store.stringStoreLock.Lock()
	defer store.stringStoreLock.Unlock()

	value, exists := store.stringStore[key]

	if exists {
		return value, nil
	} else {
		return "", fmt.Errorf("No value set for key %v\n", key)
	}
}

func (store *KeyValueStore) SetString(key, value string) error {
	store.stringStoreLock.Lock()
	defer store.stringStoreLock.Unlock()

	store.stringStore[key] = value
	return nil
}
