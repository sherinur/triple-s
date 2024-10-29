package server

import (
	"encoding/xml"
	"net/http"

	"triple-s/internal/types"
	"triple-s/internal/utils"
)

func (s *Server) HandleGetBucket(w http.ResponseWriter, r *http.Request) {
	utils.CreateDir(s.config.data_directory)

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

func (s *Server) HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
	utils.CreateDir(s.config.data_directory)
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

	// ? CONTINUE HERE
	// ! Не читает кастомный путь директории конфига
	// metadata records parsing
	records, err := utils.ParseCSV(s.config.data_directory + "/buckets.csv")
	if err != nil {
		s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// uniqness check
	if utils.FindItemByName(bucketName, records) {
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

	w.Header().Set("Content-Type", "application/xml")

	infoResponse := types.NewInfoResponse("Bucket has been successfully created")
	output, err := xml.MarshalIndent(infoResponse, "", "  ")
	if err != nil {
		s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

func (s *Server) HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	utils.CreateDir(s.config.data_directory)
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
	utils.CreateDir(s.config.data_directory)
	bucketName := r.PathValue("BucketName")

	records, err := utils.ParseCSV("./data/buckets.csv")
	if err != nil {
		s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// checking if bucket exists
	isBucketExists := utils.FindItemByName(bucketName, records)
	if !isBucketExists {
		w.WriteHeader(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/xml")

		errorResponse := types.NewErrorResponse("Bucket Not Found", "The specified bucket does not exist.")
		output, err := xml.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(output)

		s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")

		return
	}

	// checking if bucket empty
	isBucketEmpty, err := utils.IsDirEmpty("./data/" + bucketName)
	if err != nil {
		s.logger.PrintfErrorMsg("error of IsDirEmpty(): " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isBucketEmpty {
		w.WriteHeader(http.StatusConflict)

		w.Header().Set("Content-Type", "application/xml")

		errorResponse := types.NewErrorResponse("Bucket is not empty", "The specified bucket is not empty.")
		output, err := xml.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(output)

		s.logger.PrintfDebugMsg("(409 Conflict) Bucket with name '" + bucketName + "' is not empty")

		return
	}

	for i, record := range records {
		if record[0] == bucketName {
			records = append(records[:i], records[i+1:]...)
			break
		}
	}

	utils.WriteCSVbyArr(records, false)
	utils.RemoveDir("./data/" + bucketName)

	s.logger.PrintfDebugMsg("(204 No Content) Bucket with the name '" + bucketName + "' is deleted")
	w.WriteHeader(http.StatusNoContent)
}
