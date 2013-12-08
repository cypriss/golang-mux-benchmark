golang-mux-benchmark
====================

A collection of benchmarks for popular Go web frameworks.

# Frameworks

*  [gocraft/web](https://github.com/gocraft/web)
*  [gorilla/mux](https://github.com/gorilla/mux)
*  [Martini](https://github.com/codegangsta/martini)
*  [Tiger Tonic](https://github.com/rcrowley/go-tigertonic)
*  [Traffic](https://github.com/pilu/traffic)

# Benchmarks

*  **Simple** - A single route, GET /action. Renders 'hello'.
*  **RouteN** - Where N is the total number of routes that roughly approximate a REST API. The routes are as follows:
   *  Three namespaces: /admin, /api, /site
   *  Within each namespace, N/15 resources. Each resource is a random string. Each resource specifies 5 routes:
      *  GET /resources
      *  GET /resources/:id
      *  POST /resources
      *  PUT /resources/:id
      *  DELETE /resources/:id
*  **Middleware** - Run 6 middleware functions before invoking a hello handler.
*  **Composite**
   *  6 Middleware functions
   *  150 Routes
   *  The first middleware function sets a value that the handler must read and render.

# Output

`go test -bench=. 2>/dev/null` in @cypriss's 1.8 GHz i7 Macbook Air:

```
BenchmarkGocraftWebSimple	 1000000	      2612 ns/op
BenchmarkGocraftWebRoute15	  500000	      3991 ns/op
BenchmarkGocraftWebRoute75	  500000	      4048 ns/op
BenchmarkGocraftWebRoute150	  500000	      4042 ns/op
BenchmarkGocraftWebRoute300	  500000	      4090 ns/op
BenchmarkGocraftWebRoute3000	  500000	      4597 ns/op
BenchmarkGocraftWebMiddleware	  200000	      9525 ns/op
BenchmarkGocraftWebComposite	  200000	     10165 ns/op

BenchmarkGorillaMuxSimple	  500000	      3571 ns/op
BenchmarkGorillaMuxRoute15	  100000	     18029 ns/op
BenchmarkGorillaMuxRoute75	  100000	     28711 ns/op
BenchmarkGorillaMuxRoute150	   50000	     45426 ns/op
BenchmarkGorillaMuxRoute300	   50000	     70562 ns/op
BenchmarkGorillaMuxRoute3000	    5000	    654030 ns/op

BenchmarkCodegangstaMartiniSimple	  200000	      8290 ns/op
BenchmarkCodegangstaMartiniRoute15	  100000	     15445 ns/op
BenchmarkCodegangstaMartiniRoute75	  100000	     19526 ns/op
BenchmarkCodegangstaMartiniRoute150	  100000	     23655 ns/op
BenchmarkCodegangstaMartiniRoute300	   50000	     33418 ns/op
BenchmarkCodegangstaMartiniRoute3000	   10000	    206212 ns/op
BenchmarkCodegangstaMartiniMiddleware	  100000	     18662 ns/op
BenchmarkCodegangstaMartiniComposite	   50000	     37604 ns/op

BenchmarkRcrowleyTigerTonicSimple	 5000000	       378 ns/op
BenchmarkRcrowleyTigerTonicRoute15	  200000	     26696 ns/op
BenchmarkRcrowleyTigerTonicRoute75	  500000	     14189 ns/op
BenchmarkRcrowleyTigerTonicRoute150	  500000	      9485 ns/op
BenchmarkRcrowleyTigerTonicRoute300	  500000	      7124 ns/op
BenchmarkRcrowleyTigerTonicRoute3000	  500000	      4464 ns/op

BenchmarkPiluTrafficSimple	  500000	      3624 ns/op
BenchmarkPiluTrafficRoute15	  200000	     10896 ns/op
BenchmarkPiluTrafficRoute75	  100000	     20687 ns/op
BenchmarkPiluTrafficRoute150	   50000	     33320 ns/op
BenchmarkPiluTrafficRoute300	   50000	     57293 ns/op
BenchmarkPiluTrafficRoute3000	   10000	    670427 ns/op
BenchmarkPiluTrafficMiddleware	  500000	      4102 ns/op
BenchmarkPiluTrafficComposite	   50000	     34593 ns/op
```
