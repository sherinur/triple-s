package server

import "net/http"

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
}

func HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("XML: LIST OF BUCKETS"))
}

func HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.)
}