package server

import "net/http"

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func HandleCreateBucket() {
	
}
