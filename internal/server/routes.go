package server

import (
	"net/http"

	"triple-s/internal/utils"
)

func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")

	if !utils.IsValidBucketName(bucketName) {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.PrintfInfoMsg("Bucket with name '" + bucketName + "' is not valid")
		return
	}

	records, err := utils.ParseCSV("buckets.csv")
	if err != nil {
		return
	}

	if utils.FindBucketByName(bucketName, records) {
		w.WriteHeader(http.StatusConflict)
		s.logger.PrintfInfoMsg("Bucket with name '" + bucketName + "' is not unique")
		return
	}

	// TODO: Create dir using config data_directory
	// TODO: Create ./data/buckets.csv metadata file
	err = utils.CreateBucketMeta(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.PrintfErrorMsg(err.Error())
	}

	utils.CreateDir(bucketName)

	s.logger.PrintfInfoMsg("Bucket with name '" + bucketName + "' is created")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	// TODO: Read ./data/buckets.csv
	// TODO: Format and return a XML list of buckets
	// TODO: 200 OK
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	// TODO: The server checks if {BucketName} exists and is empty
	w.WriteHeader(http.StatusNoContent)
	// TODO: The bucket is removed from the ./data/buckets.csv
	w.Write([]byte("XML: BUCKET DELETED: " + bucketName))
}
