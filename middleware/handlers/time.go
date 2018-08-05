package handlers

import (
	"fmt"
	"net/http"
	"time"
)

//TimeHandler handles requests for the /time resource
func TimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "the current time is %s", time.Now().Format(time.Kitchen))
}
