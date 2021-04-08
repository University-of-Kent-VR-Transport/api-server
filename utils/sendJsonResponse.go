package utils

import (
	"os"
	"fmt"
	"encoding/json"
	"compress/gzip"
	"net/http"
)

// SendJSONResponse sends back a JSON response and compresses using gzip
// where possible
func SendJSONResponse(
	w http.ResponseWriter, successStatus int, compress bool, body interface{},
) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !compress {
		w.WriteHeader(successStatus)

		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Fprint(os.Stderr, "Error encoding (json-plain) body", err)
		}

		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(successStatus)

	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(body); err != nil {
		fmt.Fprint(os.Stderr, "Error encoding (json-gzip) body", err)
	}
}