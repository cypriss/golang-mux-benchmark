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
BenchmarkGocraftWebSimple	 1000000	      2526 ns/op
BenchmarkGocraftWebRoute15	  500000	      3954 ns/op
BenchmarkGocraftWebRoute75	  500000	      4286 ns/op
BenchmarkGocraftWebRoute150	  500000	      4101 ns/op
BenchmarkGocraftWebRoute300	  500000	      4006 ns/op
BenchmarkGocraftWebRoute3000	  500000	      4482 ns/op
BenchmarkGocraftWebMiddleware	  200000	      9539 ns/op
BenchmarkGocraftWebComposite	  200000	     10310 ns/op

BenchmarkGorillaMuxSimple	  500000	      3567 ns/op
BenchmarkGorillaMuxRoute15	  100000	     18429 ns/op
BenchmarkGorillaMuxRoute75	  100000	     28917 ns/op
BenchmarkGorillaMuxRoute150	   50000	     42757 ns/op
BenchmarkGorillaMuxRoute300	   50000	     70676 ns/op
BenchmarkGorillaMuxRoute3000	    5000	    665517 ns/op

BenchmarkCodegangstaMartiniSimple	  200000	      8007 ns/op
BenchmarkCodegangstaMartiniRoute15	  100000	     15344 ns/op
BenchmarkCodegangstaMartiniRoute75	  100000	     18813 ns/op
BenchmarkCodegangstaMartiniRoute150	  100000	     23315 ns/op
BenchmarkCodegangstaMartiniRoute300	   50000	     32734 ns/op
BenchmarkCodegangstaMartiniRoute3000	   10000	    271873 ns/op

BenchmarkRcrowleyTigerTonicSimple	 5000000	       497 ns/op
BenchmarkRcrowleyTigerTonicRoute15	  200000	     35149 ns/op
BenchmarkRcrowleyTigerTonicRoute75	  500000	     18845 ns/op
BenchmarkRcrowleyTigerTonicRoute150	  500000	     11731 ns/op
BenchmarkRcrowleyTigerTonicRoute300	  500000	      8649 ns/op
BenchmarkRcrowleyTigerTonicRoute3000	  500000	      5312 ns/op

BenchmarkPiluTrafficSimple	  500000	      4158 ns/op
BenchmarkPiluTrafficRoute15	  200000	     13150 ns/op
BenchmarkPiluTrafficRoute75	  100000	     24709 ns/op
BenchmarkPiluTrafficRoute150	   50000	     37977 ns/op
BenchmarkPiluTrafficRoute300	   50000	     68799 ns/op
BenchmarkPiluTrafficRoute3000	   10000	    756036 ns/op
```
