package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

//HelloHandler handles requests for the /hello resource
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(name) == 0 {
		http.Error(w, "please provide a 'name' parameter", http.StatusBadRequest)
		return
	}
	name = strings.Title(name)
	fmt.Fprintf(w, "Hello, %s!", name)
}
