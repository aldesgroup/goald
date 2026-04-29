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
	yamlData, err := os.ReadFile("data/api-doc.yaml")
	if err != nil {
		http.Error(w, "Failed to read API spec: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build the HTML page embedding the YAML inline and using Scalar via CDN
	html := fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script id="api-reference" type="application/yaml">
%s
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`, string(yamlData))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}
