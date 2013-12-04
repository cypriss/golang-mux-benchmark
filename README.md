golang-mux-benchmark
====================

A collection of benchmarks for popular Go web frameworks.

# Frameworks

*  [gocraft/web](https://github.com/gocraft/web)
*  [gorilla/mux](https://github.com/gorilla/mux)
*  [Martini](https://github.com/codegangsta/martini)
*  [Tiger Tonic](https://github.com/rcrowley/go-tigertonic)

# Output

`go test -bench=. 2>/dev/null` in @rcrowley's 4GB VirtualBox VM on a 2 GHz Core i7 Macbook Air:

```
BenchmarkGocraftWebSimple        1000000              2488 ns/op
BenchmarkGocraftWebRoute15        500000              4449 ns/op
BenchmarkGocraftWebRoute75        500000              4533 ns/op
BenchmarkGocraftWebRoute150       500000              4765 ns/op
BenchmarkGocraftWebRoute300       500000              4764 ns/op
BenchmarkGocraftWebRoute3000      500000              5009 ns/op

BenchmarkGorillaMuxSimple         500000              4028 ns/op
BenchmarkGorillaMuxRoute15        100000             17143 ns/op
BenchmarkGorillaMuxRoute75        100000             26580 ns/op
BenchmarkGorillaMuxRoute150        50000             39561 ns/op
BenchmarkGorillaMuxRoute300        50000             64948 ns/op
BenchmarkGorillaMuxRoute3000        5000            885580 ns/op

BenchmarkCodegangstaMartiniSimple         200000              9059 ns/op
BenchmarkCodegangstaMartiniRoute15        100000             17695 ns/op
BenchmarkCodegangstaMartiniRoute75        100000             18826 ns/op
BenchmarkCodegangstaMartiniRoute150       100000             22027 ns/op
BenchmarkCodegangstaMartiniRoute300       100000             30392 ns/op
BenchmarkCodegangstaMartiniRoute3000       10000            213453 ns/op

BenchmarkTigerTonicTrieServeMux  5000000               355 ns/op
BenchmarkTigerTonicRoute15        200000             34310 ns/op
BenchmarkTigerTonicRoute75        500000             18440 ns/op
BenchmarkTigerTonicRoute150       500000             12542 ns/op
BenchmarkTigerTonicRoute300       500000              9716 ns/op
BenchmarkTigerTonicRoute3000      500000              5151 ns/op
```
