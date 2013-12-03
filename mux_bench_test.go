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
)

//
// Types / Methods needed by gocraft/web:
//
type BenchContext struct{}


//
// Benchmarks for gocraft/web:
//
func BenchmarkGocraftWebSimple(b *testing.B) {
	router := web.New(BenchContext{})
	router.Get("/action",func(rw web.ResponseWriter, r *web.Request) {
		fmt.Fprintf(rw, "hello")
	})
	
	rw, req := testRequest("GET", "/action")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

//
// Benchmarks for gorilla/mux:
//
func BenchmarkGorillaMuxSimple(b *testing.B) {
	router := mux.NewRouter()
	router.HandleFunc("/action", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "hello")
	}).Methods("GET")
	
	rw, req := testRequest("GET", "/action")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
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

//
// Benchmarks for rcrowley/go-tigertonic's tigertonic.TrieServeMux:
//
func BenchmarkTigerTonicTrieServeMux(b *testing.B) {
	mux := tigertonic.NewTrieServeMux()
	mux.HandleFunc("GET", "/action", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	w, r := testRequest("GET", "/action")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

//
// Helpers:
//
func testRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	return recorder, request
}
