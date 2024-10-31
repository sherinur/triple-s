package server

import (
	"net/http"

	"triple-s/internal/buckets"
	"triple-s/internal/types"
	"triple-s/internal/utils"
	"triple-s/pkg/csvutil"
)

func (s *Server) HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	// checking the directory and creating if not exists
	if exists, err := utils.FileExists(s.config.data_directory); err != nil {
		s.logger.PrintfErrorMsg("error checking directory: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Directory Initialization Error", "Please check server logs and directory permissions.", w, r)
		return
	} else if !exists {
		if err := utils.CreateDir(s.config.data_directory); err != nil {
			s.logger.PrintfErrorMsg("error creating directory: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Directory Initialization Error", "Please check server logs and directory permissions.", w, r)
			return
		}
	}

	// checking the bucket metadata and creating if not exists
	metadataPath := s.config.data_directory + "/buckets.csv"

	if exists, err := utils.FileExists(metadataPath); err != nil {
		s.logger.PrintfErrorMsg("error checking file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Error", "Please check server logs and file permissions.", w, r)
		return
	} else if !exists {
		if err := utils.CreateFile(metadataPath); err != nil {
			s.logger.PrintfErrorMsg("error creating file: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Creation Error", "Please check server logs and file permissions.", w, r)
			return
		}
	}

	// Opening CSV metadata file
	file, err := csvutil.OpenCSVForRead(metadataPath)
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
		return
	}

	// Parsing CSV
	records, err := file.ReadAllRecords()
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Reading Error", "Please check server logs and file permissions.", w, r)
		return
	}
	file.Close()

	// Check for empty metadata
	if len(records) == 0 {
		s.WriteInfoResponse(http.StatusOK, "No buckets available", w, r)
		return
	}

	// Convert records arr to array of Buckets
	var bucketsSlice []types.Bucket
	for _, record := range records {
		bucket, err := buckets.ConvertArrToBucket(record)
		if err != nil {
			s.logger.PrintfErrorMsg("error converting array to bucket (HandleListBuckets): " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Checking Error", "Please check server logs and file permissions.", w, r)
			return
		}

		bucketsSlice = append(bucketsSlice, bucket)
	}

	listAllBucketsResult := types.ListAllBucketsResult{
		Buckets: struct {
			Bucket []types.Bucket `xml:"Bucket"`
		}{
			Bucket: bucketsSlice,
		},
	}

	// Write response - > StatusOK and XML of result to the body
	s.WriteXMLResponse(http.StatusOK, listAllBucketsResult, w, r)
}
