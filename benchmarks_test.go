package xor8_test

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	xor8 "github.com/getumen/xorfilter-leveldb-benchmark"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	k = 1024
	m = k * k
)

type keyGenerator interface {
	Key(i int) []byte
}

func newDB(b *testing.B, o *opt.Options) *leveldb.DB {

	now := time.Now()
	db, err := leveldb.OpenFile(fmt.Sprintf("/tmp/db_%s_%d", b.Name(), now.UnixNano()), o)

	if err != nil {
		b.Fatalf("fail to open leveldb: %+v", err)
	}

	return db
}

type randKeyGenerator struct {
	r *rand.Rand
}

func newRandKeyGenerator(r *rand.Rand) keyGenerator {
	return &randKeyGenerator{
		r: r,
	}
}

func (r *randKeyGenerator) Key(_ int) []byte {
	b := make([]byte, 8)
	// not thread safe
	u64 := r.r.Uint64()
	binary.LittleEndian.PutUint64(b, u64)
	return b
}

type seqKeyGenerator struct{}

func newSeqKeyGenerator() keyGenerator {
	return &seqKeyGenerator{}
}

func (s *seqKeyGenerator) Key(i int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func fullDB(db *leveldb.DB, g keyGenerator, num int) {
	for i := 0; i < num; i++ {
		db.Put(g.Key(i), []byte("ねむねむにゃんこパラダイス"), nil)
	}
}

func BenchmarkBloomFilterLevelDBGetRandom1M(b *testing.B) {

	// 8bits bloom filter offers a 0.3% false-positive probability
	filter := filter.NewBloomFilter(8)
	db := newDB(b, &opt.Options{
		Filter: filter,
	})

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	g := newRandKeyGenerator(r)

	fullDB(db, g, m)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = db.Get(g.Key(i), nil)
	}
}

func BenchmarkXorFilterLevelDBGetRandom1M(b *testing.B) {

	// 8bits bloom filter offers a 0.3% false-positive probability
	filter := xor8.NewXorFilter()
	db := newDB(b, &opt.Options{
		Filter: filter,
	})

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	g := newRandKeyGenerator(r)

	fullDB(db, g, m)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = db.Get(g.Key(i), nil)
	}
}

func BenchmarkBloomFilterLevelDBGetSequence1M(b *testing.B) {

	// 8bits bloom filter offers a 0.3% false-positive probability
	filter := filter.NewBloomFilter(8)
	db := newDB(b, &opt.Options{
		Filter: filter,
	})

	g := newSeqKeyGenerator()

	fullDB(db, g, m)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := db.Get(g.Key(i%m), nil)
		if err == leveldb.ErrNotFound {
			b.Fatalf("error: existing key not found")
		}
	}
}

func BenchmarkXorFilterLevelDBGetSequence1M(b *testing.B) {

	// 8bits bloom filter offers a 0.3% false-positive probability
	filter := xor8.NewXorFilter()
	db := newDB(b, &opt.Options{
		Filter: filter,
	})

	g := newSeqKeyGenerator()

	fullDB(db, g, m)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := db.Get(g.Key(i%m), nil)
		if err == leveldb.ErrNotFound {
			b.Fatalf("error: existing key not found")
		}
	}
}
