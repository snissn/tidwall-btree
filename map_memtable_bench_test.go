package btree

import (
	"encoding/binary"
	"testing"
)

type memtableBenchValue struct {
	value []byte
	flags byte
}

var memtableBenchSink int

func memtableBenchKeys(n int) []string {
	keys := make([]string, n)
	var buf [8]byte
	for i := range keys {
		binary.BigEndian.PutUint64(buf[:], uint64(i))
		keys[i] = string(buf[:])
	}
	return keys
}

func memtableBenchRandomOrder(n int) []int {
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	var x uint64 = 0x9e3779b97f4a7c15
	for i := n - 1; i > 0; i-- {
		x ^= x << 7
		x ^= x >> 9
		j := int(x % uint64(i+1))
		order[i], order[j] = order[j], order[i]
	}
	return order
}

func benchmarkMapMemtableSet(b *testing.B, order []int) {
	keys := memtableBenchKeys(len(order))
	value := memtableBenchValue{value: []byte("value")}
	opts := MapOptions{
		ReuseRightSplitCapacity:  true,
		ReuseSplitInsertCapacity: true,
		LeafItemArena:            true,
		NodeArena:                true,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMapWithOptions[string, memtableBenchValue](32, opts)
		for _, idx := range order {
			m.Set(keys[idx], value)
		}
		memtableBenchSink += m.Len()
	}
}

func BenchmarkMapMemtableSetRandom(b *testing.B) {
	const n = 65536
	benchmarkMapMemtableSet(b, memtableBenchRandomOrder(n))
}

func BenchmarkMapMemtableSetMixedAppendRandom(b *testing.B) {
	const n = 65536
	order := memtableBenchRandomOrder(n / 2)
	mixed := make([]int, 0, n)
	for i := 0; i < n/2; i++ {
		mixed = append(mixed, i+n/2)
		mixed = append(mixed, order[i])
	}
	benchmarkMapMemtableSet(b, mixed)
}

func BenchmarkMapMemtableLoadAppend(b *testing.B) {
	const n = 65536
	keys := memtableBenchKeys(n)
	value := memtableBenchValue{value: []byte("value")}
	opts := MapOptions{
		ReuseRightSplitCapacity:  true,
		ReuseSplitInsertCapacity: true,
		LeafItemArena:            true,
		NodeArena:                true,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMapWithOptions[string, memtableBenchValue](32, opts)
		for _, key := range keys {
			m.Load(key, value)
		}
		memtableBenchSink += m.Len()
	}
}
