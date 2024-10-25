package server

import (
	"net/http"

	"triple-s/internal/utils"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("XML: BUCKET CREATED: " + bucketName))

	utils.CreateDir("data")
	s.logger.PrintfInfoMsg("Bucket with name '" + bucketName + "' is created")
}

func HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketNameString := r.URL.Query().Get("BucketName")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("XML: BUCKET DELETED: " + bucketNameString))
}
