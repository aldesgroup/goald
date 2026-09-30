// ------------------------------------------------------------------------------------------------
// Here we implement the generic endpoint service the API's documentation
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func serveDocForAPI(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Read the local YAML file
	html, err := os.ReadFile("data/api-doc.html")
	if err != nil {
		http.Error(w, "Failed to read API spec: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(html))
}
