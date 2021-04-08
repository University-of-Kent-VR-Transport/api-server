package utils

import (
	"fmt"
	"net/http"
	"strings"
)

// OptionsResponse creates and sends back to the client a valid OPTION response
func OptionsResponse(
	w http.ResponseWriter,
	acceptedMethods []string,
	contentType string,
) {
	w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
	w.Header().Set("Accept", contentType)
	w.Header().Set("Accept-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, http.StatusText(http.StatusOK))
}