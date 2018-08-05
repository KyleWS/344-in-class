package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/info344-a17/info344-in-class/zipsvr/handlers"

	"github.com/info344-a17/info344-in-class/zipsvr/models"
)

const zipsPath = "/zips/"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	w.Header().Add("Content-Type", "text/plain")

	fmt.Fprintf(w, "Hello %s!", name)
}

func memoryHandler(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	//get the value of the ADDR environment variable
	//and use that as the address this server will listen on
	addr := os.Getenv("ADDR")
	//if not set, default to ":443", which means listen for
	//all requests to all hosts on port 443
	if len(addr) == 0 {
		addr = ":443"
	}

	tlskey := os.Getenv("TLSKEY")
	tlscert := os.Getenv("TLSCERT")
	if len(tlskey) == 0 || len(tlscert) == 0 {
		log.Fatal("please set TLSKEY and TLSCERT")
	}

	//load the zips and report any errors
	zips, err := models.LoadZips("zips.csv")
	if err != nil {
		//if we get an error loading the zips, our
		//server can't function, so use log.Fatalf()
		//to log the error and exit the process
		log.Fatalf("error loading zips: %v", err)
	}

	//report how many zips we loaded
	log.Printf("loaded %d zips", len(zips))

	//build an index from city name (lower-cased) to
	//the zip codes for that city
	cityIndex := models.ZipIndex{}
	//the for...range loop iterates over a slice, assigning
	//the element index and element value to the variables on the left
	//since we don't need the index, we use _ to ignore it
	for _, z := range zips {
		//convert the city name to lower-case
		cityLower := strings.ToLower(z.City)
		//append the Zip struct to the ZipSlice for that city
		cityIndex[cityLower] = append(cityIndex[cityLower], z)
	}

	//fmt.Println("Hello World!")
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/memory", memoryHandler)

	//create a new handlers.CityHandler struct
	//since that is in a different package, use the
	//package name as a prefix, and import the package above
	cityHandler := &handlers.CityHandler{
		Index:      cityIndex,
		PathPrefix: zipsPath,
	}
	//add the handler to the mux using .Handle() instead
	//of .HandleFunc(). The former is used for structs that
	//implement the http.Handler interface, while the latter
	//is used for simple functions that conform to the
	//http.HandlerFunc type.
	//see https://drstearns.github.io/tutorials/goweb/#sechandlers
	mux.Handle(zipsPath, cityHandler)

	fmt.Printf("server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, mux))
}
