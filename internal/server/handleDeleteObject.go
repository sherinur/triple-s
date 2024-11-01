package server

import "net/http"

func (s *Server) HandleDeleteObject(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if the bucket exists
	// TODO: Check if the object exists
	// TODO: Delete object file
	// TODO: Delete object metadata
	w.WriteHeader(http.StatusNoContent)
}
