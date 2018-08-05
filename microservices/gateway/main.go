package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
)

type User struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func GetCurrentUser(r *http.Request) *User {
	//does some magic with our sessions packages
	return &User{
		FirstName: "Test",
		LastName:  "User",
	}
}

func NewServiceProxy(addrs []string) *httputil.ReverseProxy {
	nextIndex := 0
	mx := sync.Mutex{}
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			//modify the request to indicate
			//remote host
			user := GetCurrentUser(r)
			userJSON, err := json.Marshal(user)
			if err != nil {
				log.Printf("error marshaling user: %v", err)
			}
			r.Header.Add("X-User", string(userJSON))

			mx.Lock()
			r.URL.Host = addrs[nextIndex%len(addrs)]
			nextIndex++
			mx.Unlock()
			r.URL.Scheme = "http"
		},
	}
}

//RootHandler handles requests for the root resource
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello from the gateway! Try requesting /v1/time")
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	//TODO: get network addresses for our
	//timesvc instances
	timesvcAddrs := os.Getenv("TIMESVC_ADDRS")
	splitTimeSvcAddrs := strings.Split(timesvcAddrs, ",")

	hellosvcAddrs := os.Getenv("HELLOSVC_ADDRS")
	splitHelloSvcAddrs := strings.Split(hellosvcAddrs, ",")

	nodeSvcAddrs := os.Getenv("NODESVC_ADDRS")
	splitNodeSvcAddrs := strings.Split(nodeSvcAddrs, ",")

	mux := http.NewServeMux()
	mux.HandleFunc("/", RootHandler)
	//TODO: add reverse proxy handler for `/v1/time`
	mux.Handle("/v1/time", NewServiceProxy(splitTimeSvcAddrs))
	mux.Handle("/v1/hello", NewServiceProxy(splitHelloSvcAddrs))
	mux.Handle("/v1/users/me/hello", NewServiceProxy(splitNodeSvcAddrs))
	mux.Handle("/v1/channels", NewServiceProxy(splitNodeSvcAddrs))

	log.Printf("server is listening at https://%s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, "tls/fullchain.pem", "tls/privkey.pem", mux))
}
