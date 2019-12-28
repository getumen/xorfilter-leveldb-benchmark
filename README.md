# xorfilter-leveldb-benchmark

LevelDB uses Bloom filters to reduce to load segment files.
Recently xor filter [Xor Filters: Faster and Smaller Than Bloom and Cuckoo Filters](https://arxiv.org/abs/1912.08258) was proposed whose performance is better than Bloom filters.

## benchmark results
### settings
- key: 1M
- false positive rate: about 0.3%

### results
TODO: 