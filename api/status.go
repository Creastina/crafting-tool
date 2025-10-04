package api

import "net/http"

func getHealth(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
