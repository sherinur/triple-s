package server

import (
	"encoding/xml"
	"net/http"

	"triple-s/internal/types"

	"triple-s/internal/utils"
)

func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) HandleGetBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")

	records, err := utils.ParseCSV("./data/buckets.csv")
	if err != nil {
		s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result types.GetBucketResult
	var recordToFind []string
	isBucketFound := false

	for _, record := range records {
		if record[0] == bucketName {
			isBucketFound = true
			recordToFind = record
			break
		}
	}

	if !isBucketFound {
		w.WriteHeader(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/xml")

		errorResponse := types.NewErrorResponse("Bucket not found", "The requested bucket could not be found.")
		output, err := xml.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(output)
	} else {
		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "application/xml")

		bucket, err := utils.ConvertArrToBucket(recordToFind)
		if err != nil {
			s.logger.PrintfErrorMsg("error converting array to bucket (HandleGetBucket): " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result = types.GetBucketResult{
			Bucket: bucket,
		}

		output, err := xml.MarshalIndent(result, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

func (s *Server) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")

	// validation
	if !utils.IsValidBucketName(bucketName) {
		w.WriteHeader(http.StatusBadRequest)

		w.Header().Set("Content-Type", "application/xml")

		errorResponse := types.NewErrorResponse("Bucket name must be valid", "The provided bucket name is not valid.")
		output, err := xml.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(output)

		s.logger.PrintfDebugMsg("(400 Bad Request) Bucket with name '" + bucketName + "' is not valid")
		return
	}

	records, err := utils.ParseCSV("./data/buckets.csv")
	if err != nil {
		s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// uniqness check
	if utils.FindBucketByName(bucketName, records) {
		w.WriteHeader(http.StatusConflict)

		w.Header().Set("Content-Type", "application/xml")

		errorResponse := types.NewErrorResponse("Bucket name must be unique", "The provided bucket name is already in use.")
		output, err := xml.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(output)

		s.logger.PrintfDebugMsg("(409 Conflict) Bucket with name '" + bucketName + "' is not unique")
		return
	}

	err = utils.CreateBucket(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.PrintfErrorMsg(err.Error())
		return
	}

	s.logger.PrintfDebugMsg("Creation of bucket with the name '" + bucketName + "'")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	records, err := utils.ParseCSV("./data/buckets.csv")
	if err != nil {
		s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buckets []types.Bucket

	if len(records) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, record := range records {
		bucket, err := utils.ConvertArrToBucket(record)
		if err != nil {
			s.logger.PrintfErrorMsg("error converting array to bucket (HandleListBuckets): " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		buckets = append(buckets, bucket)
	}

	result := types.ListAllBucketsResult{
		Buckets: struct {
			Bucket []types.Bucket `xml:"Bucket"`
		}{
			Bucket: buckets,
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	output, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func (s *Server) HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	// TODO: The server checks if {BucketName} exists and is empty
	w.WriteHeader(http.StatusNoContent)
	// TODO: The bucket is removed from the ./data/buckets.csv
	w.Write([]byte("XML: BUCKET DELETED: " + bucketName))
}
