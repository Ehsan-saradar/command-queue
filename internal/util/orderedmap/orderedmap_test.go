package orderedmap

import (
	"reflect"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	om := NewOrderedMap()

	// Test Set and Get
	om.Set("b", 2)
	om.Set("a", 1)
	om.Set("c", 3)

	expectedKeys := []string{"b", "a", "c"}
	if keys := om.Keys(); !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Keys mismatch. Expected: %v, Got: %v", expectedKeys, keys)
	}

	value, ok := om.Get("a")
	if !ok || value != 1 {
		t.Errorf("Failed to get value for key 'a'. Expected: 1, Got: %v", value)
	}

	// Test DeleteItem
	om.DeleteItem("a")

	expectedKeysAfterDelete := []string{"b", "c"}
	if keys := om.Keys(); !reflect.DeepEqual(keys, expectedKeysAfterDelete) {
		t.Errorf("Keys mismatch after DeleteItem. Expected: %v, Got: %v", expectedKeysAfterDelete, keys)
	}

	// Test Get after DeleteItem
	value, ok = om.Get("a")
	if ok {
		t.Errorf("Key 'a' should not be present after DeleteItem, but got value: %v", value)
	}

	// Test Set to update existing value
	om.Set("b", 42)
	if value, ok := om.Get("b"); !ok || value != 42 {
		t.Errorf("Failed to update value for key 'b'. Expected: 42, Got: %v", value)
	}
}
