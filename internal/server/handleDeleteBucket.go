package server

import (
	"net/http"

	"triple-s/internal/utils"
	"triple-s/pkg/csvutil"
)

func (s *Server) HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
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

			// debug log
			s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")
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
	file.Close()

	// Searching for a bucket in metadata records
	bucketIndex, found := csvutil.FindInSlice(bucketName, records)
	if found {
		bucketDirName := s.config.data_directory + "/" + bucketName

		isDirEmpty, err := utils.IsDirEmpty(bucketDirName)
		if err != nil {
			s.logger.PrintfErrorMsg("error of IsDirEmpty(): " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Bucket Deleting Error", "Please check server logs and file permissions.", w, r)
			return
		}

		// deleting bucket directory
		if !isDirEmpty {
			s.WriteErrorResponse(http.StatusConflict, "Bucket is not empty", "The specified bucket is not empty.", w, r)

			s.logger.PrintfDebugMsg("(409 Conflict) Bucket with name '" + bucketName + "' is not empty")
			return
		} else {
			utils.RemoveDir(bucketDirName)
		}

		// deleting bucket info from records
		records = utils.RemoveValue(records, bucketIndex)

		// opening metadata to rewrite
		newMetadata, err := csvutil.OpenCSVForWrite(metadataPath)
		if err != nil {
			s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
			return
		}

		// rewriting
		newMetadata.RecordsToCSV(records)

		// status
		w.WriteHeader(http.StatusNoContent)

		// debug log
		s.logger.PrintfDebugMsg("(204 No Content) Bucket with the name '" + bucketName + "' is deleted")
	} else {
		// debug log
		s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")

		// info response
		s.WriteErrorResponse(http.StatusNotFound, "Bucket not found", "The requested bucket could not be found.", w, r)
	}
}
