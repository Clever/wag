# benchmarking json vs msgpack

```
BenchmarkFormatters/JSON--size_5-10                80086             14627 ns/op                 1.316 KB
BenchmarkFormatters/JSON--size_10-10               42763             27946 ns/op                 2.592 KB
BenchmarkFormatters/JSON--size_100-10               4359            264715 ns/op                25.53 KB
BenchmarkFormatters/JSON--size_1000-10               448           2673102 ns/op               254.9 KB
BenchmarkFormatters/JSON--size_10000-10               39          27467369 ns/op              2549 KB
BenchmarkFormatters/MsgPack--size_5-10            207471              5649 ns/op                 0.9258 KB
BenchmarkFormatters/MsgPack--size_10-10           113991             10400 ns/op                 1.819 KB
BenchmarkFormatters/MsgPack--size_100-10           12346             97254 ns/op                17.91 KB
BenchmarkFormatters/MsgPack--size_1000-10           1246            952653 ns/op               178.7 KB
BenchmarkFormatters/MsgPack--size_10000-10           126           9482415 ns/op              1787 KB
```

This is both serializing and deserializing the same data.