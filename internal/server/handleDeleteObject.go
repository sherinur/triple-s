package server

import (
	"net/http"
	"triple-s/internal/utils"
	"triple-s/pkg/csvutil"
)

func (s *Server) HandleDeleteObject(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if the bucket exists
	// Values from endpoint
	bucketName := r.PathValue("BucketName")
	objectKey := r.PathValue("ObjectKey")

	// checking the storage directory and creating if not exists
	if storageExists, err := utils.FileExists(s.config.data_directory); err != nil {
		s.logger.PrintfErrorMsg("error checking directory: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Directory Initialization Error", "Please check server logs and directory permissions.", w, r)
		return
	} else if !storageExists {
		if err := utils.CreateDir(s.config.data_directory); err != nil {
			s.logger.PrintfErrorMsg("error creating directory: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Directory Initialization Error", "Please check server logs and directory permissions.", w, r)
			return
		}
	}

	// checking the bucket metadata and creating if not exists
	bucketsMetadataPath := s.config.data_directory + "/buckets.csv"

	if exists, err := utils.FileExists(bucketsMetadataPath); err != nil {
		s.logger.PrintfErrorMsg("error checking file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Error", "Please check server logs and file permissions.", w, r)
	} else if !exists {
		if err := utils.CreateFile(bucketsMetadataPath); err != nil {
			s.logger.PrintfErrorMsg("error creating file: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Creation Error", "Please check server logs and file permissions.", w, r)
			return
		}

		s.WriteErrorResponse(http.StatusNotFound, "Bucket not found", "The requested bucket could not be found.", w, r)

		// debug log
		s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")

		return
	}

	// Opening buckets CSV metadata file
	file, err := csvutil.OpenCSVForRead(bucketsMetadataPath)
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
		return
	}

	// Parsing buckets CSV
	bucketRecords, err := file.ReadAllRecords()
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Reading Error", "Please check server logs and file permissions.", w, r)
		return
	}

	file.Close()

	// Searching for a bucket in metadata records
	_, bucketFound := csvutil.FindInSlice(bucketName, bucketRecords)
	if !bucketFound {
		// debug log
		s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")

		s.WriteErrorResponse(http.StatusNotFound, "Bucket not found", "The requested bucket could not be found.", w, r)
		return
	}

	// checking the bucket directory and creating if not exists
	bucketDirPath := s.config.data_directory + "/" + bucketName

	if bucketDirExists, err := utils.FileExists(bucketDirPath); err != nil {
		s.logger.PrintfErrorMsg("error checking directory: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Directory Initialization Error", "Please check server logs and directory permissions.", w, r)
		return
	} else if !bucketDirExists {
		if err := utils.CreateDir(bucketDirPath); err != nil {
			s.logger.PrintfErrorMsg("error creating directory: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Directory Initialization Error", "Please check server logs and directory permissions.", w, r)
			return
		}
	}

	// checking the objects metadata and creating if not exists
	objectsMetadataPath := bucketDirPath + "/objects.csv"

	if exists, err := utils.FileExists(objectsMetadataPath); err != nil {
		s.logger.PrintfErrorMsg("error checking file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Error", "Please check server logs and file permissions.", w, r)
		return
	} else if !exists {
		if err := utils.CreateFile(objectsMetadataPath); err != nil {
			s.logger.PrintfErrorMsg("error creating file: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Metadata Creation Error", "Please check server logs and file permissions.", w, r)
			return
		}
	}

	// Opening CSV metadata file
	objectsFile, err := csvutil.OpenCSVForRead(objectsMetadataPath)
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
		return
	}

	// Parsing objects CSV
	objectRecords, err := objectsFile.ReadAllRecords()
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Reading Error", "Please check server logs and file permissions.", w, r)
		return
	}

	objectsFile.Close()

	// searching for an object in objects metadata records to overwrite
	objectRecordIndex, objectFound := csvutil.FindInSlice(objectKey, objectRecords)
	if !objectFound {
		// debug log
		s.logger.PrintfDebugMsg("(404 Not Found) Object with key '" + bucketName + "' does not exist")

		// not found error
		s.WriteErrorResponse(http.StatusNotFound, "Object not found", "The requested object could not be found.", w, r)
		return
	}

	// deleting from records slice
	objectRecords = utils.RemoveValue(objectRecords, objectRecordIndex)

	// opening metadata to rewrite the objects csv
	newObjectMetadata, err := csvutil.OpenCSVForWrite(objectsMetadataPath)
	if err != nil {
		s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
		return
	}

	// rewriting
	newObjectMetadata.RecordsToCSV(objectRecords)

	// debug log
	s.logger.PrintfDebugMsg("(204 No Content) Bucket with the name '" + bucketName + "' is deleted")

	// status
	w.WriteHeader(http.StatusNoContent)
}
