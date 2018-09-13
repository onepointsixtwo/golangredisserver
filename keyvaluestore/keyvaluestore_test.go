package keyvaluestore

import (
	"testing"
)

func TestKeyValueStoreSetKey(t *testing.T) {
	store := New()

	store.SetString("key", "value")

	if store.stringStore["key"] != "value" {
		t.Errorf("Value in story for key 'key' should be 'value' but is %v", store.stringStore["key"])
	}
}

func TestKeyValueStoreReadKey(t *testing.T) {
	store := New()

	store.SetString("key1", "value1")

	value, err := store.StringForKey("key1")

	if err != nil || value != "value1" {
		t.Errorf("Value should be 'value1' but was %v", value)
	}
}
