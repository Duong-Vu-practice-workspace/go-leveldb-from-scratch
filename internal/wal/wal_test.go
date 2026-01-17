package wal

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestWAL_WriteRead(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "wal_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := NewWriter(tmpFile)

	records := []Record{
		{Type: TypePut, Key: []byte("key1"), Value: []byte("value1")},
		{Type: TypeDelete, Key: []byte("key2"), Value: nil},
		{Type: TypePut, Key: []byte("key3"), Value: []byte("value3")},
	}

	for _, r := range records {
		if err := writer.Append(r); err != nil {
			t.Fatalf("failed to append record: %v", err)
		}
	}

	// Seek back to start for reading
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek: %v", err)
	}

	reader := NewReader(tmpFile)
	for i, expected := range records {
		got, err := reader.Next()
		if err != nil {
			t.Fatalf("failed to read record %d: %v", i, err)
		}

		if got.Type != expected.Type {
			t.Errorf("record %d: expected type %d, got %d", i, expected.Type, got.Type)
		}
		if !bytes.Equal(got.Key, expected.Key) {
			t.Errorf("record %d: expected key %s, got %s", i, string(expected.Key), string(got.Key))
		}
		if !bytes.Equal(got.Value, expected.Value) {
			t.Errorf("record %d: expected value %s, got %s", i, string(expected.Value), string(got.Value))
		}
	}

	// Check EOF
	_, err = reader.Next()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}
