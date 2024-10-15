package main

import (
	"flag"
	"log"
	"net/http"

	"triple-s/pkg/logger"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SHERIII SHERIII GOOO"))
}

func main() {
	port := flag.String("port", "4400", "Port number")
	dir := flag.String("dir", ".", "Path to the directory")

	flag.Parse()

	logger.PrintfInfoMsg("Starting server on port :" + *port)
	logger.PrintfInfoMsg("Path to the directory set: " + *dir)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	err := http.ListenAndServe(":"+*port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
