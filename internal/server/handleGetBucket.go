package server

import (
	"encoding/xml"
	"net/http"

	"triple-s/internal/buckets"
	"triple-s/internal/types"
	"triple-s/internal/utils"
	"triple-s/pkg/csvutil"
)

func (s *Server) HandleGetBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")

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

	// Searching for a bucket in metadata records
	bucketIndex, found := csvutil.FindInSlice(bucketName, records)
	if found {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/xml")

		bucketValues := records[bucketIndex]
		bucket, err := buckets.ConvertArrToBucket(bucketValues)
		if err != nil {
			s.logger.PrintfErrorMsg("error converting arr to bucket: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Response formatting error", "Please check server logs and file permissions.", w, r)
			return
		}

		getBucketResult := types.GetBucketResult{
			Bucket: bucket,
		}

		// formatting to xml
		output, err := xml.MarshalIndent(getBucketResult, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Response formatting error", "Please check server logs and file permissions.", w, r)
			return
		}

		w.Write(output)
	} else {
		s.WriteErrorResponse(http.StatusNotFound, "Bucket not found", "The requested bucket could not be found.", w, r)
	}
}
