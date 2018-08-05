package handlers

import "github.com/info344-a17/info344-in-class/zipsvr/models"
import "net/http"
import "strings"
import "encoding/json"

//CityHandler is an HTTP handler that returns zip codes
//for a given city name. It requires a ZipIndex from
//city name (lower-case) to a ZipSlice, as well as
//a PathPrefix, which is the resource path prefix (e.g., "/zips/"")
type CityHandler struct {
	PathPrefix string
	Index      models.ZipIndex
}

//ServeHTTP handles HTTP requests for the CityHandler. This is a method
//of the CityHandler struct defined above. Methods in Go use a receiver
//parameter defined on the left, which will be an instance of the struct.
//The receiver parameter is exactly like the `this` pointer in Java, just
//more explicitly defined. For more details on receiver parameters, see
//https://drstearns.github.io/tutorials/golang/#secreceivers
func (ch *CityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// URL: /zips/city-name
	//slice off just the city name from the resource path
	//since we know the length of the PathPrefix, we can slice off
	//just the characters that follow that PathPrefix
	cityName := r.URL.Path[len(ch.PathPrefix):]

	//convert the city name to lower case since the ZipIndex map
	//keys are all lower-cased as well
	cityName = strings.ToLower(cityName)

	//if the city name is zero-length respond with an error
	if len(cityName) == 0 {
		//the http.Error() method writes an error message to the response
		//and sets the HTTP status code to the value of the third parameter
		http.Error(w, "please provide a city name", http.StatusBadRequest)
		//since http.Error() writes a response, we should return to
		//stop processing this request.
		return
	}

	//add the header `Content-Type: application/json`
	w.Header().Add(headerContentType, contentTypeJSON)
	//add the CORS header `Access-Control-Allow-Origin: *`
	//see https://drstearns.github.io/tutorials/cors/
	w.Header().Add(headerAccessControlAllowOrigin, "*")

	//get the ZipSlice for the requested city name
	zips := ch.Index[cityName]
	//write that slice to the response, encoded as JSON
	json.NewEncoder(w).Encode(zips)
}
