package server

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"triple-s/internal/logger"
	"triple-s/internal/types"
)

type Server struct {
	config *Config
	logger *logger.Logger
	mux    *http.ServeMux
}

// New server
func New(config *Config) *Server {
	s := &Server{
		config: config,
		logger: logger.New(true),
		mux:    http.NewServeMux(),
	}

	s.registerRoutes()
	return s
}

// TODO: Продолжить по видео REST API на Golang

// Start the server
func (s *Server) Start() error {
	s.logger.PrintfInfoMsg("Starting server on port " + s.config.port)
	s.logger.PrintfInfoMsg("Path to the directory set: " + s.config.data_directory)
	s.logger.PrintfInfoMsg("Path to the config set: " + s.config.cfg_file)

	mux := s.RequestMiddleware(s.mux)

	return http.ListenAndServe(s.config.port, mux)
}

func (s *Server) registerRoutes() {
	// basic routes
	s.mux.HandleFunc("GET /health", s.HandleHealth)

	// bucket routes
	s.mux.HandleFunc("GET /{BucketName}", s.HandleGetBucket)
	s.mux.HandleFunc("GET /", s.HandleListBuckets)
	s.mux.HandleFunc("PUT /{BucketName}", s.HandleCreateBucket)
	s.mux.HandleFunc("DELETE /{BucketName}", s.HandleDeleteBucket)

	// object routes
	s.mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", s.HandlePutObject)
	s.mux.HandleFunc("GET /{BucketName}/{ObjectKey}", s.HandleGetObject)
	s.mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", s.HandleDeleteObject)
}

func (s *Server) RequestMiddleware(next http.Handler) http.Handler {
	allowedMethods := map[string]bool{
		http.MethodGet:    true,
		http.MethodPost:   true,
		http.MethodPut:    true,
		http.MethodDelete: true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.PrintfInfoMsg(fmt.Sprintf("Request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr))

		if !allowedMethods[r.Method] {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusMethodNotAllowed)

			errorResponse := types.NewErrorResponse("Method Not Allowed", "The HTTP method is not allowed for this endpoint.")
			output, err := xml.MarshalIndent(errorResponse, "", "  ")
			if err != nil {
				s.logger.PrintfErrorMsg("Error encoding XML: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(output)
			return
		}

		next.ServeHTTP(w, r)
	})
}
