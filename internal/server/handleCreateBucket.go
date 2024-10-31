package server

import (
	"net/http"

	"triple-s/internal/buckets"
	"triple-s/internal/types"
	"triple-s/internal/utils"
	"triple-s/pkg/csvutil"
)

func (s *Server) HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")

	// bucket name validation
	if !buckets.ValidateBucketName(bucketName) {
		s.logger.PrintfDebugMsg("(400 Bad Request) Bucket with name '" + bucketName + "' is not valid")
		s.WriteErrorResponse(http.StatusBadRequest, "Bucket name must be valid", "The provided bucket name is not valid.", w, r)
		return
	}

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

	// Searching for a bucket in metadata records
	// (checking if the bucket is unique)
	_, found := csvutil.FindInSlice(bucketName, records)
	if found {
		s.logger.PrintfDebugMsg("(409 Conflict) Bucket with name '" + bucketName + "' is not unique")
		s.WriteErrorResponse(http.StatusConflict, "Bucket name must be unique", "The provided bucket name is already in use.", w, r)
		return
	}

	// Appending new bucket to csv metadata
	s.logger.PrintfDebugMsg("Creation of bucket with the name '" + bucketName + "'")

	newBucket := types.NewBucket(bucketName)

	metadata, err := csvutil.OpenCSVForAppend(metadataPath)
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Reading Error", "Please check server logs and file permissions.", w, r)
	}
	defer file.Close()

	newRecord := buckets.ConvertBucketToArr(newBucket)
	metadata.AppendToCSV(newRecord)
	if err != nil {
		s.logger.PrintfErrorMsg("error writing CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Writing Error", "Please check server logs and file permissions.", w, r)
	}
	metadata.Close()

	// Creating bucket directory
	err = utils.CreateDir(s.config.data_directory + "/" + bucketName)
	if err != nil {
		s.logger.PrintfErrorMsg("error creating bucket dir: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Dir Creating Error", "Please check server logs and directory permissions.", w, r)
	}

	// Writing info response
	s.WriteInfoResponse(http.StatusOK, "Bucket has been successfully created", w, r)

	// err = utils.CreateBucket(bucketName)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	s.logger.PrintfErrorMsg(err.Error())
	// 	return
	// }

	// s.logger.PrintfDebugMsg("Creation of bucket with the name '" + bucketName + "'")

	// w.WriteHeader(http.StatusOK)

	// w.Header().Set("Content-Type", "application/xml")

	// infoResponse := types.NewInfoResponse("Bucket has been successfully created")
	// output, err := xml.MarshalIndent(infoResponse, "", "  ")
	// if err != nil {
	// 	s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(output)
}
