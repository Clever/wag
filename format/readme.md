# benchmarking json vs msgpack

```
BenchmarkFormatters/JSON--size_5-10                79628             14339 ns/op                 1.316 KB
BenchmarkFormatters/JSON--size_10-10               43141             27771 ns/op                 2.592 KB
BenchmarkFormatters/JSON--size_100-10               4393            262171 ns/op                25.53 KB
BenchmarkFormatters/JSON--size_1000-10               448           2659057 ns/op               254.9 KB
BenchmarkFormatters/JSON--size_10000-10               39          26617225 ns/op              2549 KB
BenchmarkFormatters/MsgPack--size_5-10            534350              2187 ns/op                 0.9258 KB
BenchmarkFormatters/MsgPack--size_10-10           292224              3911 ns/op                 1.819 KB
BenchmarkFormatters/MsgPack--size_100-10           33459             35535 ns/op                17.91 KB
BenchmarkFormatters/MsgPack--size_1000-10           3361            339654 ns/op               178.7 KB
BenchmarkFormatters/MsgPack--size_10000-10           372           3197495 ns/op              1787 KB
```

This is both serializing and deserializing the same data.