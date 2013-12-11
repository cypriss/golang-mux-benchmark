package mux_bench_test

import (
	"crypto/sha1"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/gocraft/web"
	"github.com/gorilla/mux"
	"github.com/pilu/traffic"
	"github.com/rcrowley/go-tigertonic"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//
// Types used by any/all frameworks:
//
type RouterBuilder func(namespaces []string, resources []string) http.Handler

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

func BenchmarkGocraftWeb_Simple(b *testing.B) {
	router := web.New(BenchContext{})
	router.Get("/action", gocraftWebHandler)

	rw, req := testRequest("GET", "/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

func BenchmarkGocraftWeb_Route15(b *testing.B) {
	benchmarkRoutesN(b, 1, gocraftWebRouterFor)
}

func BenchmarkGocraftWeb_Route75(b *testing.B) {
	benchmarkRoutesN(b, 5, gocraftWebRouterFor)
}

func BenchmarkGocraftWeb_Route150(b *testing.B) {
	benchmarkRoutesN(b, 10, gocraftWebRouterFor)
}

func BenchmarkGocraftWeb_Route300(b *testing.B) {
	benchmarkRoutesN(b, 20, gocraftWebRouterFor)
}

func BenchmarkGocraftWeb_Route3000(b *testing.B) {
	benchmarkRoutesN(b, 200, gocraftWebRouterFor)
}

func BenchmarkGocraftWeb_Middleware(b *testing.B) {
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

func BenchmarkGocraftWeb_Composite(b *testing.B) {
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
func gorillaMuxRouterFor(namespaces []string, resources []string) http.Handler {
	router := mux.NewRouter()
	for _, ns := range namespaces {
		subrouter := router.PathPrefix("/" + ns).Subrouter()
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

func BenchmarkGorillaMux_Simple(b *testing.B) {
	router := mux.NewRouter()
	router.HandleFunc("/action", helloHandler).Methods("GET")

	rw, req := testRequest("GET", "/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

func BenchmarkGorillaMux_Route15(b *testing.B) {
	benchmarkRoutesN(b, 1, gorillaMuxRouterFor)
}

func BenchmarkGorillaMux_Route75(b *testing.B) {
	benchmarkRoutesN(b, 5, gorillaMuxRouterFor)
}

func BenchmarkGorillaMux_Route150(b *testing.B) {
	benchmarkRoutesN(b, 10, gorillaMuxRouterFor)
}

func BenchmarkGorillaMux_Route300(b *testing.B) {
	benchmarkRoutesN(b, 20, gorillaMuxRouterFor)
}

func BenchmarkGorillaMux_Route3000(b *testing.B) {
	benchmarkRoutesN(b, 200, gorillaMuxRouterFor)
}

//
// Benchmarks for codegangsta/martini:
//
type martiniContext struct {
	MyField string
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

func BenchmarkCodegangstaMartini_Simple(b *testing.B) {
	r := martini.NewRouter()
	m := martini.New()
	m.Action(r.Handle)

	r.Get("/action", helloHandler)

	rw, req := testRequest("GET", "/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(rw, req)
	}
}

func BenchmarkCodegangstaMartini_Route15(b *testing.B) {
	benchmarkRoutesN(b, 1, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartini_Route75(b *testing.B) {
	benchmarkRoutesN(b, 5, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartini_Route150(b *testing.B) {
	benchmarkRoutesN(b, 10, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartini_Route300(b *testing.B) {
	benchmarkRoutesN(b, 20, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartini_Route3000(b *testing.B) {
	benchmarkRoutesN(b, 200, codegangstaMartiniRouterFor)
}

func BenchmarkCodegangstaMartini_Middleware(b *testing.B) {
	martiniMiddleware := func(rw http.ResponseWriter, r *http.Request, c martini.Context) {
		c.Next()
	}

	r := martini.NewRouter()
	m := martini.New()
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Action(r.Handle)

	r.Get("/action", helloHandler)

	rw, req := testRequest("GET", "/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(rw, req)
		if rw.Code != 200 {
			panic("no good")
		}
	}
}

func BenchmarkCodegangstaMartini_Composite(b *testing.B) {
	namespaces, resources, requests := resourceSetup(10)

	martiniMiddleware := func(rw http.ResponseWriter, r *http.Request, c martini.Context) {
		c.Next()
	}

	handler := func(rw http.ResponseWriter, r *http.Request, c *martiniContext) {
		fmt.Fprintf(rw, c.MyField)
	}

	r := martini.NewRouter()
	m := martini.New()
	m.Use(func(rw http.ResponseWriter, r *http.Request, c martini.Context) {
		c.Map(&martiniContext{MyField: r.URL.Path})
		c.Next()
	})
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Use(martiniMiddleware)
	m.Action(r.Handle)

	for _, ns := range namespaces {
		for _, res := range resources {
			r.Get("/"+ns+"/"+res, handler)
			r.Post("/"+ns+"/"+res, handler)
			r.Get("/"+ns+"/"+res+"/:id", handler)
			r.Put("/"+ns+"/"+res+"/:id", handler)
			r.Delete("/"+ns+"/"+res+"/:id", handler)
		}
	}
	benchmarkRoutes(b, m, requests)
}

//
// Benchmarks for rcrowley/go-tigertonic's tigertonic.TrieServeMux:
//
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

func BenchmarkRcrowleyTigerTonic_Simple(b *testing.B) {
	mux := tigertonic.NewTrieServeMux()
	mux.HandleFunc("GET", "/action", helloHandler)
	rw, r := testRequest("GET", "/action")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(rw, r)
	}
}

func BenchmarkRcrowleyTigerTonic_Route15(b *testing.B) {
	benchmarkRoutesN(b, 1, tigertonicRouterFor)
}

func BenchmarkRcrowleyTigerTonic_Route75(b *testing.B) {
	benchmarkRoutesN(b, 5, tigertonicRouterFor)
}

func BenchmarkRcrowleyTigerTonic_Route150(b *testing.B) {
	benchmarkRoutesN(b, 10, tigertonicRouterFor)
}

func BenchmarkRcrowleyTigerTonic_Route300(b *testing.B) {
	benchmarkRoutesN(b, 20, tigertonicRouterFor)
}

func BenchmarkRcrowleyTigerTonic_Route3000(b *testing.B) {
	benchmarkRoutesN(b, 200, tigertonicRouterFor)
}

//
// Benchmarks for pilu/traffic:
//
func piluTrafficHandler(rw traffic.ResponseWriter, r *traffic.Request) {
	fmt.Fprintf(rw, "hello")
}

func piluTrafficCompositeHandler(rw traffic.ResponseWriter, r *traffic.Request) {
	fieldVal := rw.GetVar("field").(string)
	fmt.Fprintf(rw, fieldVal)
}

type trafficMiddleware struct{}
type trafficCompositeMiddleware struct{}

func (middleware *trafficMiddleware) ServeHTTP(w traffic.ResponseWriter, r *traffic.Request, next traffic.NextMiddlewareFunc) {
	if nextMiddleware := next(); nextMiddleware != nil {
		nextMiddleware.ServeHTTP(w, r, next)
	}
}

func (middleware *trafficCompositeMiddleware) ServeHTTP(w traffic.ResponseWriter, r *traffic.Request, next traffic.NextMiddlewareFunc) {
	if nextMiddleware := next(); nextMiddleware != nil {
		w.SetVar("field", r.URL.Path)
		nextMiddleware.ServeHTTP(w, r, next)
	}
}

func piluTrafficRouterFor(namespaces []string, resources []string) http.Handler {
	traffic.SetVar("env", "production")
	router := traffic.New()
	for _, ns := range namespaces {
		for _, res := range resources {
			router.Get("/"+ns+"/"+res, piluTrafficHandler)
			router.Post("/"+ns+"/"+res, piluTrafficHandler)
			router.Get("/"+ns+"/"+res+"/:id", piluTrafficHandler)
			router.Put("/"+ns+"/"+res+"/:id", piluTrafficHandler)
			router.Delete("/"+ns+"/"+res+"/:id", piluTrafficHandler)
		}
	}
	return router
}

func BenchmarkPiluTraffic_Simple(b *testing.B) {
	traffic.SetVar("env", "production")
	router := traffic.New()
	router.Get("/action", piluTrafficHandler)
	rw, r := testRequest("GET", "/action")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, r)
	}
}

func BenchmarkPiluTraffic_Route15(b *testing.B) {
	benchmarkRoutesN(b, 1, piluTrafficRouterFor)
}

func BenchmarkPiluTraffic_Route75(b *testing.B) {
	benchmarkRoutesN(b, 5, piluTrafficRouterFor)
}

func BenchmarkPiluTraffic_Route150(b *testing.B) {
	benchmarkRoutesN(b, 10, piluTrafficRouterFor)
}

func BenchmarkPiluTraffic_Route300(b *testing.B) {
	benchmarkRoutesN(b, 20, piluTrafficRouterFor)
}

func BenchmarkPiluTraffic_Route3000(b *testing.B) {
	benchmarkRoutesN(b, 200, piluTrafficRouterFor)
}

func BenchmarkPiluTraffic_Middleware(b *testing.B) {
	traffic.SetVar("env", "production")
	router := traffic.New()
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Get("/action", piluTrafficHandler)

	rw, req := testRequest("GET", "/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
		if rw.Code != 200 {
			panic("no good")
		}
	}
}

func BenchmarkPiluTraffic_Composite(b *testing.B) {
	namespaces, resources, requests := resourceSetup(10)

	traffic.SetVar("env", "production")
	router := traffic.New()
	router.Use(&trafficCompositeMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})
	router.Use(&trafficMiddleware{})

	for _, ns := range namespaces {
		for _, res := range resources {
			router.Get("/"+ns+"/"+res, piluTrafficCompositeHandler)
			router.Post("/"+ns+"/"+res, piluTrafficCompositeHandler)
			router.Get("/"+ns+"/"+res+"/:id", piluTrafficCompositeHandler)
			router.Put("/"+ns+"/"+res+"/:id", piluTrafficCompositeHandler)
			router.Delete("/"+ns+"/"+res+"/:id", piluTrafficCompositeHandler)
		}
	}
	benchmarkRoutes(b, router, requests)
}

//
// Helpers:
//
func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "hello")
}

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

		if recorder.Code != 200 {
			panic("wat")
		}

		reqId += 1
	}
}
