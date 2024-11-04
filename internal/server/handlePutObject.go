package server

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"triple-s/internal/objects"
	"triple-s/internal/types"
	"triple-s/internal/utils"
	"triple-s/pkg/csvutil"
)

func (s *Server) HandlePutObject(w http.ResponseWriter, r *http.Request) {
	// Values from endpoint
	bucketName := r.PathValue("BucketName")
	objectKey := r.PathValue("ObjectKey")

	if !objects.ValidateObjectKey(objectKey) {
		s.WriteErrorResponse(http.StatusBadRequest, "Invalid Object Key", "Object key must be between 1 and 1024 characters and cannot contain invalid characters.", w, r)
		s.logger.PrintfDebugMsg("(400 Bad Request) Invalid object key")
		return
	}

	// Response headers
	contentLength := strconv.FormatInt(r.ContentLength, 10)
	contentType := r.Header.Get("Content-Type")

	// content type check
	if contentType == "" {
		s.WriteErrorResponse(http.StatusUnsupportedMediaType, "Unsupported Media Type", "The request is missing the Content-Type header.", w, r)
		s.logger.PrintfDebugMsg("(415 Unsupported Media Type) The request is missing the Content-Type header")
		return
	}

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

	// Reading content
	content, err := io.ReadAll(r.Body)
	if err != nil {
		s.WriteErrorResponse(http.StatusBadRequest, "Content Reading Error", "Please check server logs and file permissions.", w, r)
		return
	}
	defer r.Body.Close()

	// searching for a object in objects metadata records to overwrite
	objectRecordIndex, objectFound := csvutil.FindInSlice(objectKey, objectRecords)
	if objectFound {
		// updating existing record in CSV
		object, err := objects.ConvertArrToObject(objectRecords[objectRecordIndex])
		if err != nil {
			s.logger.PrintfErrorMsg("error creating object: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "Object Creating Error", "Please check server logs and file permissions.", w, r)
			return
		}
		object.LastModifiedTime = time.Now()
		object.Size = contentLength
		object.ContentType = contentType

		objectRecords[objectRecordIndex] = objects.ConvertObjectToArr(&object)

		// open objects metadata to overwrite
		updatedObjectsFile, err := csvutil.OpenCSVForWrite(objectsMetadataPath)
		if err != nil {
			s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
			return
		}

		updatedObjectsFile.RecordsToCSV(objectRecords)
		updatedObjectsFile.Close()
	} else {
		// appending new record to CSV
		object := types.NewObject(objectKey, contentLength, contentType)
		newObjectRecord := objects.ConvertObjectToArr(object)

		// open objects metadata to append
		updatedObjectsFile, err := csvutil.OpenCSVForAppend(objectsMetadataPath)
		if err != nil {
			s.logger.PrintfErrorMsg("error opening CSV file: " + err.Error())
			s.WriteErrorResponse(http.StatusInternalServerError, "File Opening Error", "Please check server logs and file permissions.", w, r)
			return
		}

		updatedObjectsFile.AppendToCSV(newObjectRecord)
		updatedObjectsFile.Close()
	}

	// Writing object content
	contentPath := bucketDirPath + "/" + objectKey
	contentFile, err := os.Create(contentPath)
	if err != nil {
		s.logger.PrintfErrorMsg("error creating file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Content Creation Error", "Please check server logs and file permissions.", w, r)
		return
	}
	defer contentFile.Close()

	_, err = contentFile.Write(content)
	if err != nil {
		s.logger.PrintfErrorMsg("error writing file: " + err.Error())
		s.WriteErrorResponse(http.StatusInternalServerError, "Content Creation Error", "Please check server logs and file permissions.", w, r)
		return
	}

	s.WriteInfoResponse(http.StatusOK, "The object has been successfully stored.", w, r)
	s.logger.PrintfDebugMsg("Putting object with the key '" + objectKey + "'")
}
