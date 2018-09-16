package expiry

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/ttltimer"
	"sync"
)

type Handler struct {
	dataStore     keyvaluestore.Store
	timersMap     map[string]*ttltimer.TTLTimer
	timersMapLock *sync.Mutex
}

func New(dataStore keyvaluestore.Store) *Handler {
	return &Handler{dataStore: dataStore, timersMap: make(map[string]*ttltimer.TTLTimer), timersMapLock: &sync.Mutex{}}
}

func (handler *Handler) ExpireKeyAfterSeconds(key string, afterSeconds int) error {
	_, err := handler.dataStore.StringForKey(key)
	if err != nil {
		return fmt.Errorf("Cannot set expiry for nonexistent key %v", key)
	}

	// Cancel existing timer
	handler.CancelTimerForKeyIfExists(key)

	// Start new timer
	timer := ttltimer.New(afterSeconds)
	handler.storeTimerForKey(timer, key)
	go handler.runTimer(timer, key)
	return nil
}

func (handler *Handler) CancelTimerForKeyIfExists(key string) {
	handler.timersMapLock.Lock()
	defer handler.timersMapLock.Unlock()

	timer, exists := handler.timersMap[key]
	if exists {
		timer.Stop()
		delete(handler.timersMap, key)
	}
}

func (handler *Handler) RemainingExpiryTTLForKey(key string) (int, error) {
	handler.timersMapLock.Lock()
	defer handler.timersMapLock.Unlock()

	timer, exists := handler.timersMap[key]
	if exists {
		return timer.RemainingTTL(), nil
	}

	return 0, fmt.Errorf("No timer exists for key %v", key)
}

func (handler *Handler) runTimer(timer *ttltimer.TTLTimer, key string) {
	<-timer.GetTimerChannel()

	handler.removeTimerForKey(timer, key)

	fmt.Printf("Deleting expiring key: %v\n", key)
	handler.dataStore.DeleteString(key)
}

func (handler *Handler) storeTimerForKey(timer *ttltimer.TTLTimer, key string) {
	handler.timersMapLock.Lock()
	defer handler.timersMapLock.Unlock()

	handler.timersMap[key] = timer
}

func (handler *Handler) removeTimerForKey(timer *ttltimer.TTLTimer, key string) {
	handler.timersMapLock.Lock()
	defer handler.timersMapLock.Unlock()

	_, exists := handler.timersMap[key]
	if exists {
		delete(handler.timersMap, key)
	}
}
