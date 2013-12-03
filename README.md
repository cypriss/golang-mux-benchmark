golang-mux-benchmark
====================

A collection of benchmarks for popular Go web frameworks.

# Frameworks
*  [gocraft/web](https://github.com/gocraft/web)
*  [gorilla/mux](https://github.com/gorilla/mux)
*  [Martini](https://github.com/codegangsta/martini)

# Output
On my Macbook Air (1.8ghz):

```
BenchmarkGocraftWebSimple	 1000000	      2517 ns/op
BenchmarkGorillaMuxSimple	  500000	      3585 ns/op
BenchmarkCodegangstaMartiniSimple	  200000	      8026 ns/op
```

On @rcrowley's 4GB VirtualBox VM on a 2GHz Core i7 MacBook Air:

```
BenchmarkGocraftWebSimple        1000000              2457 ns/op
BenchmarkGorillaMuxSimple         500000              6270 ns/op
BenchmarkCodegangstaMartiniSimple         200000              8565 ns/op
BenchmarkTigerTonicTrieServeMux  5000000               350 ns/op
```
