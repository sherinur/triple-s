package server

import (
	"net/http"
)

// ? CONTINUE HERE
// ! Не пишется objects.csv
func (s *Server) HandlePutObject(w http.ResponseWriter, r *http.Request) {
	// bucketName := r.PathValue("BucketName")
	// objectKey := r.PathValue("ObjectKey")

	// // // TODO: Check if the bucket exists
	// // // TODO: Check if the request has content-type (error in other case), and get the content-type
	// // // TODO: Validate the object key (check is it valid) for uniqueness
	// // TODO: Update objects.csv of the bucket
	// // TODO: Create object file

	// // parsing buckets.csv
	// bucketMetaRecords, err := utils.ParseCSV("./data/buckets.csv")
	// if err != nil {
	// 	s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// // checking if bucket exists
	// isBucketExists := utils.FindItemByName(bucketName, bucketMetaRecords)
	// if !isBucketExists {
	// 	w.WriteHeader(http.StatusNotFound)

	// 	w.Header().Set("Content-Type", "application/xml")

	// 	errorResponse := types.NewErrorResponse("Bucket Not Found", "The specified bucket does not exist.")
	// 	output, err := xml.MarshalIndent(errorResponse, "", "  ")
	// 	if err != nil {
	// 		s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Write(output)

	// 	s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")

	// 	return
	// }

	// if isBucketExists {
	// 	s.logger.PrintfDebugMsg(bucketName + " exists.")
	// }

	// // parsing objects.csv
	// objectMetaRecords, err := utils.ParseCSV("./data/" + bucketName + "/objects.csv")
	// if err != nil {
	// 	s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// // checking the object for uniqueness
	// isObjectNotUnique := utils.FindItemByName(objectKey, objectMetaRecords)
	// if isObjectNotUnique {
	// 	w.WriteHeader(http.StatusConflict)

	// 	w.Header().Set("Content-Type", "application/xml")

	// 	errorResponse := types.NewErrorResponse("Object key must be unique", "The provided object key is already in use.")
	// 	output, err := xml.MarshalIndent(errorResponse, "", "  ")
	// 	if err != nil {
	// 		s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Write(output)

	// 	s.logger.PrintfDebugMsg("(409 Conflict) Object with key '" + objectKey + "' is not unique")
	// 	return
	// }

	// contentType := r.Header.Get("Content-Type")
	// if contentType == "" {
	// 	w.WriteHeader(http.StatusUnsupportedMediaType)

	// 	errorResponse := types.NewErrorResponse("Unsupported Media Type", "The request is missing the Content-Type header.")
	// 	output, err := xml.MarshalIndent(errorResponse, "", "  ")
	// 	if err != nil {
	// 		s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	s.logger.PrintfDebugMsg("(415 Unsupported Media Type) The request is missing the Content-Type header")

	// 	w.Write(output)
	// 	return
	// }

	// sizeStr := strconv.FormatInt(r.ContentLength, 10)

	// utils.CreateObject(bucketName, objectKey, sizeStr, contentType)

	// w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleGetObject(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if the bucket exists
	// TODO: Check if the object exists
	// TODO: Return the binary content of file
	// TODO: Set Content-Type header
}

func (s *Server) HandleDeleteObject(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if the bucket exists
	// TODO: Check if the object exists
	// TODO: Delete object file
	// TODO: Delete object metadata
	w.WriteHeader(http.StatusNoContent)
}
