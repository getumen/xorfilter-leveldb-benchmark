package xor8

import (
	"encoding/binary"

	"github.com/FastFilter/xorfilter"
	"github.com/syndtr/goleveldb/leveldb/filter"
)

type xorFilter struct {
}

type xorFilterGenerator struct {
	keys []uint64
}

func (f *xorFilter) Name() string {
	return "leveldb.Xor8"
}

func (f *xorFilter) NewGenerator() filter.FilterGenerator {
	return &xorFilterGenerator{
		keys: make([]uint64, 0),
	}
}

func (f *xorFilter) Contains(filter, key []byte) bool {
	fingerprints := make([]uint8, len(filter)-12)
	for i := 0; i < len(filter)-12; i++ {
		fingerprints[i] = uint8(filter[i+12])
	}
	xorfilter := &xorfilter.Xor8{
		Seed:         binary.LittleEndian.Uint64(filter[0:8]),
		BlockLength:  binary.LittleEndian.Uint32(filter[8:12]),
		Fingerprints: fingerprints,
	}
	return xorfilter.Contains(binary.LittleEndian.Uint64(key))
}

func (g *xorFilterGenerator) Add(key []byte) {
	b := make([]byte, (8-len(key)%8)%8)
	key = append(key, b...)
	arr := make([]uint64, len(key)/8)
	for i := 0; i < len(arr); i++ {
		arr[i] = binary.LittleEndian.Uint64(key[8*i : 8*(i+1)])
	}
	g.keys = append(g.keys, arr...)
}

func (g *xorFilterGenerator) Generate(buf filter.Buffer) {
	xorfilter := xorfilter.Populate(g.keys)
	// reset keys
	g.keys = make([]uint64, 0)

	filterBlob := make([]byte, 12)

	binary.LittleEndian.PutUint64(filterBlob[0:8], xorfilter.Seed)
	binary.LittleEndian.PutUint32(filterBlob[8:12], xorfilter.BlockLength)
	for i := 0; i < len(xorfilter.Fingerprints); i++ {
		filterBlob = append(filterBlob, byte(xorfilter.Fingerprints[i]))
	}

	buf.Write(filterBlob)
}

// NewXorFilter generates xor filter
func NewXorFilter() filter.Filter {
	return &xorFilter{}
}
