# Walkthrough: Phase 1 - Skip List (Memtable)

I have completed the first major phase of our LevelDB implementation: the **Skip List**. This serves as the core of our **Memtable**, providing sorted in-memory storage for fast writes.

## Changes Made

### Skip List Implementation
Implemented a probabilistic Skip List that supports:
- **Keys and Values**: Both stored as `[]byte` to handle arbitrary data.
- **Ordered Insertion**: Keys are kept in sorted order.
- **Fast Lookups**: $O(\log n)$ average time complexity for `Get` and `Put`.
- **Iterators**: Support for sequential scans, `SeekToFirst`, and `Seek` to a specific key.

Files created:
- [skiplist.go](file:///home/duongvct/Documents/workspace/goland/go-leveldb-from-scratch/internal/memtable/skiplist.go)
- [skiplist_test.go](file:///home/duongvct/Documents/workspace/goland/go-leveldb-from-scratch/internal/memtable/skiplist_test.go)

## Verification Results

### Automated Tests
I ran a suite of basic and edge-case tests to ensure the Skip List behaves as expected.

```bash
go test -v ./internal/memtable/...
```

**Results:**
```text
=== RUN   TestSkipList_Basic
--- PASS: TestSkipList_Basic (0.00s)
=== RUN   TestSkipList_Iterator
--- PASS: TestSkipList_Iterator (0.00s)
=== RUN   TestSkipList_Update
--- PASS: TestSkipList_Update (0.00s)
=== RUN   TestSkipList_Large
--- PASS: TestSkipList_Large (0.00s)
PASS
ok      github.com/duongvct/go-leveldb-from-scratch/internal/memtable        0.005s
```

## Next Steps
Now that we have a sorted in-memory container, the next step is **Phase 2: Write-Ahead Log (WAL)**. This ensures that even if the process crashes, the data in the Memtable isn't lost.
