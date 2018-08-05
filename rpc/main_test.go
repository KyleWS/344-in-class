package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/rpc"
	"testing"
)

//run these benchmarks using:
// go test -bench=. -benchmem

const urlToSummarize = "http://example.com"

func BenchmarkRPC(b *testing.B) {
	client, err := rpc.Dial("tcp", rpcAddr)
	if err != nil {
		b.Fatalf("error dialing RPC server: %v", err)
	}
	defer client.Close()
	psum := &PageSummary{}
	for i := 0; i < b.N; i++ {
		if err := client.Call("SummaryService.GetPageSummary", "http://ogp.me", psum); err != nil {
			b.Fatalf("error calling RPC: %v", err)
		}
		if psum.URL != "http://ogp.me" {
			b.Fatalf("incorrect data returned from RPC: expected http://ogp.me but got %s", psum.URL)
		}
	}
}

func BenchmarkHTTP(b *testing.B) {
	summaryURL := fmt.Sprintf("http://%s?url=http://ogp.me", httpAddr)
	psum := &PageSummary{}
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(summaryURL)
		if err != nil {
			b.Fatalf("error getting page summary: %v", err)
		}
		if err := json.NewDecoder(resp.Body).Decode(psum); err != nil {
			b.Fatalf("error decoding JSON response: %v", err)
		}
		if psum.URL != "http://ogp.me" {
			b.Fatalf("incorrect URL returned: expected http://ogp.me but got %s", psum.URL)
		}
		resp.Body.Close()
	}
}
