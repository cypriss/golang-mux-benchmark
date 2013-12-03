golang-mux-benchmark
====================

A collection of benchmarks for popular Go web frameworks.

# Frameworks
*  [gocraft/web](https://github.com/gocraft/web)
*  [gorilla/mux](https://github.com/gorilla/mux)
*  [Martini](https://github.com/codegangsta/martini)
*  [TigerTonic](https://github.com/rcrowley/go-tigertonic)

# Output
```go test -bench=. -v``` on my Macbook Air (1.8ghz):

```
BenchmarkGocraftWebSimple	 1000000	      2522 ns/op
BenchmarkGocraftWebRoute15	  500000	      3813 ns/op
BenchmarkGocraftWebRoute75	  500000	      3887 ns/op
BenchmarkGocraftWebRoute150	  500000	      3942 ns/op
BenchmarkGocraftWebRoute300	  500000	      3903 ns/op
BenchmarkGocraftWebRoute3000	  500000	      4429 ns/op

BenchmarkGorillaMuxSimple	  500000	      3564 ns/op
BenchmarkGorillaMuxRoute15	  100000	     17762 ns/op
BenchmarkGorillaMuxRoute75	  100000	     31095 ns/op
BenchmarkGorillaMuxRoute150	   50000	     41708 ns/op
BenchmarkGorillaMuxRoute300	   50000	     74333 ns/op
BenchmarkGorillaMuxRoute3000	    5000	    651023 ns/op

BenchmarkCodegangstaMartiniSimple	  200000	      8665 ns/op
BenchmarkCodegangstaMartiniRoute15	  100000	     16192 ns/op
BenchmarkCodegangstaMartiniRoute75	  100000	     20267 ns/op
BenchmarkCodegangstaMartiniRoute150	  100000	     25183 ns/op
BenchmarkCodegangstaMartiniRoute300	   50000	     34053 ns/op
BenchmarkCodegangstaMartiniRoute3000	   10000	    208872 ns/op

BenchmarkTigerTonicTrieServeMux	  5000000	       422 ns/op
```
