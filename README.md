# xorfilter-leveldb-benchmark

LevelDB uses Bloom filters to reduce to load segment files.
Recently xor filter [Xor Filters: Faster and Smaller Than Bloom and Cuckoo Filters](https://arxiv.org/abs/1912.08258) was proposed whose performance is better than Bloom filters.

## benchmark results
### settings
- key num: 1M
- false positive rate: about 0.3%
- block cache off
- others: default values https://godoc.org/github.com/syndtr/goleveldb/leveldb/opt#pkg-variables

### results
TODO: benchmark test by filter bit size

```
go test -benchmem -cpuprofile -run='^$' github.com/getumen/xorfilter-leveldb-benchmark -bench .                                                                                                                                                                                                                                                           (git)-[master]
goos: linux
goarch: amd64
pkg: github.com/getumen/xorfilter-leveldb-benchmark
BenchmarkBloomFilterLevelDBGetRandom1M-8     	    6777	    156853 ns/op	  233231 B/op	      41 allocs/op
BenchmarkXorFilterLevelDBGetRandom1M-8       	    4052	    248835 ns/op	  373798 B/op	     398 allocs/op
BenchmarkBloomFilterLevelDBGetSequence1M-8   	    5586	    200323 ns/op	  318685 B/op	      41 allocs/op
BenchmarkXorFilterLevelDBGetSequence1M-8     	    4951	    268174 ns/op	  506220 B/op	     359 allocs/op
PASS
ok  	github.com/getumen/xorfilter-leveldb-benchmark	125.346s
```
