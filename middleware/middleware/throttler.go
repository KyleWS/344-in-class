package middleware

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

//Throttler is a middleware that throttles requests
type Throttler struct {
	handler     http.Handler
	store       *redis.Client
	maxRequests int64
	duration    time.Duration
}

//NewThrottler constructs a Throttler
func NewThrottler(handler http.Handler, store *redis.Client, maxRequests int64, duration time.Duration) *Throttler {
	return &Throttler{handler, store, maxRequests, duration}
}

//ServeHTTP implements the http.Handler interface for Throttler
func (t *Throttler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO: if the client has made too many requests
	//respond with an http.StatusTooManyRequests error;
	//else, call the wrapped handler to process the request.
	//Use the redis database to track the number of requests
	//made within the duration. Hint: use the Incr() command:
	//https://godoc.org/github.com/go-redis/redis#Client.Incr

	//Use r.RemoteAddr to identify the client, though beware
	//that this will return the IP:port of the machine that
	//contacted our server, which could be a proxy sitting
	//in-between our server and the client. For a more complete
	//but much more complicated solution, see
	//https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html

}
