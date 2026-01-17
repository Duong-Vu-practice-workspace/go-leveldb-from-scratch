package wal

import (
	"encoding/binary"
	"io"
	"os"
)

const (
	TypePut    byte = 0
	TypeDelete byte = 1
)

// Record represents a single entry in the WAL.
type Record struct {
	Type  byte
	Key   []byte
	Value []byte
}

// Writer handles appending records to the WAL file.
type Writer struct {
	file *os.File
}

// NewWriter creates a new WAL writer.
func NewWriter(file *os.File) *Writer {
	return &Writer{file: file}
}

// Append writes a record to the log.
func (w *Writer) Append(r Record) error {
	// Calculate total length of the payload
	// Type (1) + KeyLen (4) + len(Key) + ValueLen (4) + len(Value)
	payloadSize := 1 + 4 + len(r.Key) + 4 + len(r.Value)
	
	buf := make([]byte, 4+payloadSize)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(payloadSize))
	buf[4] = r.Type
	
	offset := 5
	binary.LittleEndian.PutUint32(buf[offset:offset+4], uint32(len(r.Key)))
	offset += 4
	copy(buf[offset:], r.Key)
	offset += len(r.Key)
	
	binary.LittleEndian.PutUint32(buf[offset:offset+4], uint32(len(r.Value)))
	offset += 4
	copy(buf[offset:], r.Value)

	_, err := w.file.Write(buf)
	if err != nil {
		return err
	}
	
	// Sync to ensure it's written to disk
	return w.file.Sync()
}

// Reader handles reading records from the WAL file.
type Reader struct {
	reader io.Reader
}

// NewReader creates a new WAL reader.
func NewReader(reader io.Reader) *Reader {
	return &Reader{reader: reader}
}

// Next reads the next record from the log.
func (r *Reader) Next() (Record, error) {
	var record Record
	
	// Read payload size
	var payloadSize uint32
	err := binary.Read(r.reader, binary.LittleEndian, &payloadSize)
	if err != nil {
		return record, err
	}
	
	// Read payload
	payload := make([]byte, payloadSize)
	_, err = io.ReadFull(r.reader, payload)
	if err != nil {
		return record, err
	}
	
	record.Type = payload[0]
	
	offset := 1
	keyLen := binary.LittleEndian.Uint32(payload[offset : offset+4])
	offset += 4
	record.Key = payload[offset : offset+int(keyLen)]
	offset += int(keyLen)
	
	valueLen := binary.LittleEndian.Uint32(payload[offset : offset+4])
	offset += 4
	record.Value = payload[offset : offset+int(valueLen)]
	
	return record, nil
}
