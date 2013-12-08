package mux_bench_test

import (
	"github.com/gocraft/web"
	"github.com/gorilla/mux"
	"github.com/codegangsta/martini"
	"github.com/rcrowley/go-tigertonic"
	"testing"
	"fmt"
	"net/http"
	"net/http/httptest"
	"crypto/sha1"
	"io"
)

type RouterBuilder func(namespaces []string, resources []string) http.Handler

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "hello")
}

//
// Benchmarks for gocraft/web:
//
type BenchContext struct {
	MyField string
}
type BenchContextB struct {
	*BenchContext
}
type BenchContextC struct {
	*BenchContextB
}
func (c *BenchContext) Action(w web.ResponseWriter, r *web.Request) {
	fmt.Fprintf(w, "hello")
}

func (c *BenchContextB) Action(w web.ResponseWriter, r *web.Request) {
	fmt.Fprintf(w, c.MyField)
}

func gocraftWebHandler(rw web.ResponseWriter, r *web.Request) {
	fmt.Fprintf(rw, "hello")
}

func BenchmarkGocraftWebSimple(b *testing.B) {
	router := web.New(BenchContext{})
	router.Get("/action", gocraftWebHandler)
	
	rw, req := testRequest("GET", "/action")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

func BenchmarkGocraftWebRoute15(b *testing.B) {
	benchmarkRoutesN(b, 1, gocraftWebRouterFor)
}

func BenchmarkGocraftWebRoute75(b *testing.B) {
	benchmarkRoutesN(b, 5, gocraftWebRouterFor)
}

func BenchmarkGocraftWebRoute150(b *testing.B) {
	benchmarkRoutesN(b, 10, gocraftWebRouterFor)
}

func BenchmarkGocraftWebRoute300(b *testing.B) {
	benchmarkRoutesN(b, 20, gocraftWebRouterFor)
}

func BenchmarkGocraftWebRoute3000(b *testing.B) {
	benchmarkRoutesN(b, 200, gocraftWebRouterFor)
}

func BenchmarkGocraftWebMiddleware(b *testing.B) {
	nextMw := func(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		next(rw, r)
	}

	router := web.New(BenchContext{})
	router.Middleware(nextMw)
	router.Middleware(nextMw)
	routerB := router.Subrouter(BenchContextB{}, "/b")
	routerB.Middleware(nextMw)
	routerB.Middleware(nextMw)
	routerC := routerB.Subrouter(BenchContextC{}, "/c")
	routerC.Middleware(nextMw)
	routerC.Middleware(nextMw)
	routerC.Get("/action", gocraftWebHandler)

	rw, req := testRequest("GET", "/b/c/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
		// if rw.Code != 200 { panic("no good") }
	}
}

// Composite scenario:
// - 6 middlewares
// - top middleware allocates a context and sets a field on it
// - handler reads that value and renders it
// - 150 routes (10 resources on 3 namespaces)
func BenchmarkGocraftWebComposite(b *testing.B) {
	namespaces, resources, requests := resourceSetup(10)
	
	nextMw := func(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		next(rw, r)
	}
	
	router := web.New(BenchContext{})
	router.Middleware(func(c *BenchContext, rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		c.MyField = r.URL.Path
		next(rw, r)
	})
	router.Middleware(nextMw)
	router.Middleware(nextMw)
	
	for _, ns := range namespaces {
		subrouter := router.Subrouter(BenchContextB{}, "/"+ns)
		subrouter.Middleware(nextMw)
		subrouter.Middleware(nextMw)
		subrouter.Middleware(nextMw)
		for _, res := range resources {
			subrouter.Get("/"+res, (*BenchContextB).Action)
			subrouter.Post("/"+res, (*BenchContextB).Action)
			subrouter.Get("/"+res+"/:id", (*BenchContextB).Action)
			subrouter.Put("/"+res+"/:id", (*BenchContextB).Action)
			subrouter.Delete("/"+res+"/:id", (*BenchContextB).Action)
		}
	}
	benchmarkRoutes(b, router, requests)
}

//
// Benchmarks for gorilla/mux:
//
func BenchmarkGorillaMuxSimple(b *testing.B) {
	router := mux.NewRouter()
	router.HandleFunc("/action", helloHandler).Methods("GET")
	
	rw, req := testRequest("GET", "/action")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

func BenchmarkGorillaMuxRoute15(b *testing.B) {
	benchmarkRoutesN(b, 1, gorillaMuxRouterFor)
}

func BenchmarkGorillaMuxRoute75(b *testing.B) {
	benchmarkRoutesN(b, 5, gorillaMuxRouterFor)
}

func BenchmarkGorillaMuxRoute150(b *testing.B) {
	benchmarkRoutesN(b, 10, gorillaMuxRouterFor)
}

func BenchmarkGorillaMuxRoute300(b *testing.B) {
	benchmarkRoutesN(b, 20, gorillaMuxRouterFor)
}

func BenchmarkGorillaMuxRoute3000(b *testing.B) {
	benchmarkRoutesN(b, 200, gorillaMuxRouterFor)
}



//
// Benchmarks for codegangsta/martini:
//
func BenchmarkCodegangstaMartiniSimple(b *testing.B) {
	r := martini.NewRouter()
	m := martini.New()
	m.Action(r.Handle)
	
	r.Get("/action", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "hello")
	})
	
	rw, req := testRequest("GET", "/action")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(rw, req)
	}
}

func BenchmarkCodegangstaMartiniRoute15(b *testing.B) {
	benchmarkRoutesN(b, 1, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartiniRoute75(b *testing.B) {
	benchmarkRoutesN(b, 5, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartiniRoute150(b *testing.B) {
	benchmarkRoutesN(b, 10, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartiniRoute300(b *testing.B) {
	benchmarkRoutesN(b, 20, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartiniRoute3000(b *testing.B) {
	benchmarkRoutesN(b, 200, codegangstaMartiniRouterFor)
}

//
// Benchmarks for rcrowley/go-tigertonic's tigertonic.TrieServeMux:
//
func BenchmarkTigerTonicTrieServeMux(b *testing.B) {
	mux := tigertonic.NewTrieServeMux()
	mux.HandleFunc("GET", "/action", helloHandler)
	rw, r := testRequest("GET", "/action")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(rw, r)
	}
}

func BenchmarkTigerTonicRoute15(b *testing.B) {
	benchmarkRoutesN(b, 1, tigertonicRouterFor)
}

func BenchmarkTigerTonicRoute75(b *testing.B) {
	benchmarkRoutesN(b, 5, tigertonicRouterFor)
}

func BenchmarkTigerTonicRoute150(b *testing.B) {
	benchmarkRoutesN(b, 10, tigertonicRouterFor)
}

func BenchmarkTigerTonicRoute300(b *testing.B) {
	benchmarkRoutesN(b, 20, tigertonicRouterFor)
}

func BenchmarkTigerTonicRoute3000(b *testing.B) {
	benchmarkRoutesN(b, 200, tigertonicRouterFor)
}

//
// Helpers:
//
func testRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	return recorder, request
}

func benchmarkRoutesN(b *testing.B, N int, builder RouterBuilder) {
	namespaces, resources, requests := resourceSetup(N)
	router := builder(namespaces, resources)
	benchmarkRoutes(b, router, requests)
}

// Returns a routeset with N *resources per namespace*. so N=1 gives about 15 routes
func resourceSetup(N int) (namespaces []string, resources []string, requests []*http.Request) {
	namespaces = []string{"admin", "api", "site"}
	resources = []string{}
	
	for i := 0; i < N; i += 1 {
		sha1 := sha1.New()
		io.WriteString(sha1, fmt.Sprintf("%d", i))
		strResource := fmt.Sprintf("%x", sha1.Sum(nil))
		resources = append(resources, strResource)
	}
	
	for _, ns := range namespaces {
		for _, res := range resources {
			req, _ := http.NewRequest("GET", "/"+ns+"/"+res, nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("POST", "/"+ns+"/"+res, nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("GET", "/"+ns+"/"+res+"/3937", nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("PUT", "/"+ns+"/"+res+"/3937", nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("DELETE", "/"+ns+"/"+res+"/3937", nil)
			requests = append(requests, req)
		}
	}
	
	return
}

func benchmarkRoutes(b *testing.B, handler http.Handler, requests []*http.Request) {
	recorder := httptest.NewRecorder()
	reqId := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if reqId >= len(requests) {
			reqId = 0
		}
		req := requests[reqId]
		handler.ServeHTTP(recorder, req)

		// if recorder.Code != 200 {
		// 	panic("wat")
		// }

		reqId += 1
	}
}

func gocraftWebRouterFor(namespaces []string, resources []string) http.Handler {
	router := web.New(BenchContext{})
	for _, ns := range namespaces {
		subrouter := router.Subrouter(BenchContext{}, "/"+ns)
		for _, res := range resources {
			subrouter.Get("/"+res, (*BenchContext).Action)
			subrouter.Post("/"+res, (*BenchContext).Action)
			subrouter.Get("/"+res+"/:id", (*BenchContext).Action)
			subrouter.Put("/"+res+"/:id", (*BenchContext).Action)
			subrouter.Delete("/"+res+"/:id", (*BenchContext).Action)
		}
	}
	return router
}

func gorillaMuxRouterFor(namespaces []string, resources []string) http.Handler {
	router := mux.NewRouter()
	for _, ns := range namespaces {
		subrouter := router.PathPrefix("/"+ns).Subrouter()
		for _, res := range resources {
			subrouter.HandleFunc("/"+res, helloHandler).Methods("GET")
			subrouter.HandleFunc("/"+res, helloHandler).Methods("POST")
			subrouter.HandleFunc("/"+res+"/:id", helloHandler).Methods("GET")
			subrouter.HandleFunc("/"+res+"/:id", helloHandler).Methods("PUT")
			subrouter.HandleFunc("/"+res+"/:id", helloHandler).Methods("DELETE")
		}
	}
	return router
}

func codegangstaMartiniRouterFor(namespaces []string, resources []string) http.Handler {
	router := martini.NewRouter()
	martini := martini.New()
	martini.Action(router.Handle)
	for _, ns := range namespaces {
		for _, res := range resources {
			router.Get("/"+ns+"/"+res, helloHandler)
			router.Post("/"+ns+"/"+res, helloHandler)
			router.Get("/"+ns+"/"+res+"/:id", helloHandler)
			router.Put("/"+ns+"/"+res+"/:id", helloHandler)
			router.Delete("/"+ns+"/"+res+"/:id", helloHandler)
		}
	}
	return martini
}

func tigertonicRouterFor(namespaces []string, resources []string) http.Handler {
	mux := tigertonic.NewTrieServeMux()
	for _, ns := range namespaces {
		for _, res := range resources {
			mux.HandleFunc("GET", "/"+ns+"/"+res, helloHandler)
			mux.HandleFunc("POST", "/"+ns+"/"+res, helloHandler)
			mux.HandleFunc("GET", "/"+ns+"/"+res+"/{id}", helloHandler)
			mux.HandleFunc("POST", "/"+ns+"/"+res+"/{id}", helloHandler)
			mux.HandleFunc("DELETE", "/"+ns+"/"+res+"/{id}", helloHandler)
		}
	}
	return mux
}
