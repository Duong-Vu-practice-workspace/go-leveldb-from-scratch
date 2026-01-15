package memtable

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSkipList_Basic(t *testing.T) {
	sl := NewSkipList()

	// Insert some keys
	keys := []string{"key1", "key3", "key2", "key5", "key4"}
	for _, k := range keys {
		sl.Put([]byte(k), []byte("val-"+k))
	}

	// Verify Get
	for _, k := range keys {
		val, ok := sl.Get([]byte(k))
		if !ok {
			t.Errorf("expected to find key %s", k)
		}
		if string(val) != "val-"+k {
			t.Errorf("expected val-%s, got %s", k, string(val))
		}
	}

	// Verify non-existent key
	if _, ok := sl.Get([]byte("non-existent")); ok {
		t.Error("did not expect to find non-existent key")
	}
}

func TestSkipList_Iterator(t *testing.T) {
	sl := NewSkipList()
	sl.Put([]byte("b"), []byte("2"))
	sl.Put([]byte("a"), []byte("1"))
	sl.Put([]byte("d"), []byte("4"))
	sl.Put([]byte("c"), []byte("3"))

	it := sl.NewIterator()
	
	// Test SeekToFirst
	it.SeekToFirst()
	expected := []struct{ k, v string }{
		{"a", "1"}, {"b", "2"}, {"c", "3"}, {"d", "4"},
	}

	for i, exp := range expected {
		if !it.Valid() {
			t.Fatalf("iterator should be valid at index %d", i)
		}
		if string(it.Key()) != exp.k {
			t.Errorf("expected key %s, got %s", exp.k, string(it.Key()))
		}
		if string(it.Value()) != exp.v {
			t.Errorf("expected value %s, got %s", exp.v, string(it.Value()))
		}
		it.Next()
	}
	if it.Valid() {
		t.Error("iterator should be invalid after end")
	}

	// Test Seek
	it.Seek([]byte("b"))
	if !it.Valid() || string(it.Key()) != "b" {
		t.Errorf("expected to seek to 'b', got %s", string(it.Key()))
	}

	it.Seek([]byte("cc")) // should go to 'd'
	if !it.Valid() || string(it.Key()) != "d" {
		t.Errorf("expected to seek to 'd' when seeking 'cc', got %s", string(it.Key()))
	}

	it.Seek([]byte("z")) // should be invalid
	if it.Valid() {
		t.Error("expected invalid iterator when seeking past end")
	}
}

func TestSkipList_Update(t *testing.T) {
	sl := NewSkipList()
	sl.Put([]byte("key"), []byte("val1"))
	sl.Put([]byte("key"), []byte("val2"))

	val, ok := sl.Get([]byte("key"))
	if !ok || string(val) != "val2" {
		t.Errorf("expected val2, got %s", string(val))
	}
}

func TestSkipList_Large(t *testing.T) {
	sl := NewSkipList()
	n := 1000
	for i := 0; i < n; i++ {
		k := fmt.Sprintf("%04d", i)
		sl.Put([]byte(k), []byte(k))
	}

	for i := 0; i < n; i++ {
		k := fmt.Sprintf("%04d", i)
		val, ok := sl.Get([]byte(k))
		if !ok || !bytes.Equal(val, []byte(k)) {
			t.Errorf("failed to get key %s", k)
		}
	}
}
