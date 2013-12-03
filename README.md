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
