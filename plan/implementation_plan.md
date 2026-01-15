# Implementation Plan: LevelDB from Scratch

We will build LevelDB step-by-step, starting with the core in-memory component: the **Skip List**. This will serve as our **Memtable**, where all writes go first to stay sorted and fast.

## Proposed Changes

### [Phase 1] Skip List Implementation
The Skip List is a probabilistic data structure that allows $O(\log n)$ search and insertion. It's easier to implement than balanced trees and works great for LevelDB.

#### [NEW] [skiplist.go](file:///home/duongvct/Documents/workspace/goland/go-leveldb-from-scratch/internal/memtable/skiplist.go)
- Implement `SkipList` struct.
- Implement `Node` struct with forward pointers.
- Functions: `Insert`, `Get`, `Delete`, `NewIterator`.

#### [NEW] [skiplist_test.go](file:///home/duongvct/Documents/workspace/goland/go-leveldb-from-scratch/internal/memtable/skiplist_test.go)
- Comprehensive tests for concurrent/sequential access.

## Verification Plan

### Automated Tests
- Run `go test ./internal/memtable/...` to verify Skip List correctness.
- Property-based testing for sorted order maintenance.

### Manual Verification
- I will provide code snippets for you to run and see the Skip List in action with different data loads.
