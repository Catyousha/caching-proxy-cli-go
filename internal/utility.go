package internal

import (
	"fmt"
	"net/http"
)

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal server error: %v\n", err)
}
