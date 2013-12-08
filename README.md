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
BenchmarkGocraftWebSimple	 1000000	      2886 ns/op
BenchmarkGocraftWebRoute15	  500000	      4230 ns/op
BenchmarkGocraftWebRoute75	  500000	      4204 ns/op
BenchmarkGocraftWebRoute150	  500000	      4236 ns/op
BenchmarkGocraftWebRoute300	  500000	      4317 ns/op
BenchmarkGocraftWebRoute3000	  500000	      5133 ns/op
BenchmarkGocraftWebMiddleware	  200000	     10436 ns/op
BenchmarkGocraftWebComposite	  200000	     11022 ns/op

BenchmarkGorillaMuxSimple	  500000	      3884 ns/op
BenchmarkGorillaMuxRoute15	  100000	     20357 ns/op
BenchmarkGorillaMuxRoute75	   50000	     31369 ns/op
BenchmarkGorillaMuxRoute150	   50000	     46707 ns/op
BenchmarkGorillaMuxRoute300	   20000	     77188 ns/op
BenchmarkGorillaMuxRoute3000	    5000	    729515 ns/op

BenchmarkCodegangstaMartiniSimple	  200000	      8923 ns/op
BenchmarkCodegangstaMartiniRoute15	  100000	     17264 ns/op
BenchmarkCodegangstaMartiniRoute75	  100000	     21219 ns/op
BenchmarkCodegangstaMartiniRoute150	  100000	     26335 ns/op
BenchmarkCodegangstaMartiniRoute300	   50000	     36935 ns/op
BenchmarkCodegangstaMartiniRoute3000	   10000	    228868 ns/op

BenchmarkRcrowleyTigerTonicSimple	 5000000	       414 ns/op
BenchmarkRcrowleyTigerTonicRoute15	  200000	     28473 ns/op
BenchmarkRcrowleyTigerTonicRoute75	  500000	     15128 ns/op
BenchmarkRcrowleyTigerTonicRoute150	  500000	     10088 ns/op
BenchmarkRcrowleyTigerTonicRoute300	  500000	      7445 ns/op
BenchmarkRcrowleyTigerTonicRoute3000	  500000	      4900 ns/op

BenchmarkPiluTrafficSimple	  500000	      3921 ns/op
BenchmarkPiluTrafficRoute15	  200000	     11760 ns/op
BenchmarkPiluTrafficRoute75	  100000	     21150 ns/op
BenchmarkPiluTrafficRoute150	   50000	     35636 ns/op
BenchmarkPiluTrafficRoute300	   50000	     60219 ns/op
BenchmarkPiluTrafficRoute3000	   10000	    671760 ns/op
BenchmarkPiluTrafficMiddleware	  500000	      3950 ns/op
BenchmarkPiluTrafficComposite	   50000	     33970 ns/op
```
