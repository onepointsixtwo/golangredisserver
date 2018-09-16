package expiry

import (
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"testing"
	"time"
)

// Tests

func TestExpiringKey(t *testing.T) {
	sut := createSut()

	sut.dataStore.SetString("mykey", "myvalue")
	_ = sut.ExpireKeyAfterSeconds("mykey", 1)

	time.Sleep(time.Duration(1500) * time.Millisecond)

	value, err := sut.dataStore.StringForKey("mykey")
	if err == nil {
		t.Errorf("Key should no longer exist in store but instead exists with value of %v", value)
	}
}

func TestFailsToSetExpiryForNonexistentKey(t *testing.T) {
	sut := createSut()
	err := sut.ExpireKeyAfterSeconds("mykey", 1)
	if err == nil {
		t.Error("Should have received error attempting to expire non-existent key")
	}
}

func TestCancelExpiringKey(t *testing.T) {
	sut := createSut()

	sut.dataStore.SetString("mykey", "myvalue")
	_ = sut.ExpireKeyAfterSeconds("mykey", 1)
	sut.CancelTimerForKeyIfExists("mykey")

	if sut.timersMap["mykey"] != nil {
		t.Error("Timer should have been removed when expiry cancelled")
	}

	time.Sleep(time.Duration(1500) * time.Microsecond)

	value, err := sut.dataStore.StringForKey("mykey")
	if err != nil || value != "myvalue" {
		t.Errorf("mykey should still exist but received error %v or incorrect value %v", err, value)
	}
}

func TestTTLForKey(t *testing.T) {
	sut := createSut()

	sut.dataStore.SetString("mykey", "myvalue")
	_ = sut.ExpireKeyAfterSeconds("mykey", 5)

	ttl, err := sut.RemainingExpiryTTLForKey("mykey")
	if err != nil || ttl != 4 {
		t.Errorf("Expected TTL for key to be 4 but was %v (err: %v)", ttl, err)
	}

	// cancel as cleanup
	sut.CancelTimerForKeyIfExists("mykey")
}

// Helpers

func createSut() *Handler {
	kvs := keyvaluestore.New()
	return New(kvs)
}
