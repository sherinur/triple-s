package server

import (
	"net/http"

	"triple-s/internal/logger"
)

type Server struct {
	config *Config
	logger *logger.Logger
}

// New server
func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logger.New(),
	}
}

// TODO: Продолжить по видео REST API на Golang

// Start the server
func (s *Server) Start() error {
	s.logger.PrintfInfoMsg("Starting server on port " + s.config.port)
	s.logger.PrintfInfoMsg("Path to the directory set: " + s.config.data_directory)
	s.logger.PrintfInfoMsg("Path to the config set: " + s.config.cfg_file)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", HandleHealth)
	mux.HandleFunc("GET /", HandleListBuckets)
	mux.HandleFunc("PUT /{BucketName}", HandleCreateBucket)
	mux.HandleFunc("DELETE /{BucketName}", HandleDeleteBucket)

	loggedMux := s.logger.LogRequestMiddleware(mux)

	err := http.ListenAndServe(s.config.port, loggedMux)
	if err != nil {
		return err
	}

	return nil
}
