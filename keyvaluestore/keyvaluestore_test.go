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

func TestDeleteValueForKey(t *testing.T) {
	store := New()

	deleted := store.DeleteString("key")
	if deleted {
		t.Error("Should not have been able to delete key 'key' - it shouldn't exist!\n")
	}

	store.SetString("key2", "value")
	deleted = store.DeleteString("key2")
	if !deleted {
		t.Error("Store should have deleted key2 successfully")
	}
}
