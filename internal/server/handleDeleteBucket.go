package server

import (
	"net/http"
)

func (s *Server) HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	// 	utils.CreateDir(s.config.data_directory)
	// 	bucketName := r.PathValue("BucketName")

	// 	records, err := utils.ParseCSV("./data/buckets.csv")
	// 	if err != nil {
	// 		s.logger.PrintfErrorMsg("error reading CSV: " + err.Error())
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// checking if bucket exists
	// 	isBucketExists := utils.FindItemByName(bucketName, records)
	// 	if !isBucketExists {
	// 		w.WriteHeader(http.StatusNotFound)

	// 		w.Header().Set("Content-Type", "application/xml")

	// 		errorResponse := types.NewErrorResponse("Bucket Not Found", "The specified bucket does not exist.")
	// 		output, err := xml.MarshalIndent(errorResponse, "", "  ")
	// 		if err != nil {
	// 			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}

	// 		w.Write(output)

	// 		s.logger.PrintfDebugMsg("(404 Not Found) Bucket with name '" + bucketName + "' does not exist")

	// 		return
	// 	}

	// 	// checking if bucket empty
	// 	isBucketEmpty, err := utils.IsDirEmpty("./data/" + bucketName)
	// 	if err != nil {
	// 		s.logger.PrintfErrorMsg("error of IsDirEmpty(): " + err.Error())
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if !isBucketEmpty {
	// 		w.WriteHeader(http.StatusConflict)

	// 		w.Header().Set("Content-Type", "application/xml")

	// 		errorResponse := types.NewErrorResponse("Bucket is not empty", "The specified bucket is not empty.")
	// 		output, err := xml.MarshalIndent(errorResponse, "", "  ")
	// 		if err != nil {
	// 			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}

	// 		w.Write(output)

	// 		s.logger.PrintfDebugMsg("(409 Conflict) Bucket with name '" + bucketName + "' is not empty")

	// 		return
	// 	}

	// 	for i, record := range records {
	// 		if record[0] == bucketName {
	// 			records = append(records[:i], records[i+1:]...)
	// 			break
	// 		}
	// 	}

	// 	utils.WriteCSVbyArr(records, false)
	// 	utils.RemoveDir("./data/" + bucketName)

	// s.logger.PrintfDebugMsg("(204 No Content) Bucket with the name '" + bucketName + "' is deleted")
	// w.WriteHeader(http.StatusNoContent)
}
