package server

import (
	"encoding/xml"
	"net/http"

	"triple-s/internal/types"
)

func (s *Server) WriteErrorResponse(httpStatus int, errorMsg string, message string, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(httpStatus)

	if message != "" && errorMsg != "" {
		w.Header().Set("Content-Type", "application/xml")

		errorResponse := types.NewErrorResponse(errorMsg, message)
		output, err := xml.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		w.Write(output)
	}

	return nil
}

func (s *Server) WriteInfoResponse(httpStatus int, message string, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(httpStatus)

	if message != "" {
		w.Header().Set("Content-Type", "application/xml")

		infoResponse := types.NewInfoResponse(message)
		output, err := xml.MarshalIndent(infoResponse, "", "  ")
		if err != nil {
			s.logger.PrintfErrorMsg("error encoding XML: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		w.Write(output)
	}

	return nil
}
