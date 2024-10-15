package main

import (
	"flag"
	"log"
	"net/http"
	"triple-s/pkg/logger"
)

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr) // Логируем метод, путь и IP-адрес
		next.ServeHTTP(w, r)                                                      // Вызываем следующий обработчик
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SHERIII SHERIII GOOO"))
}

func gay(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("YOU ARE GAY"))
}

func main() {
	var err error

	port := flag.String("port", "4400", "Port number")
	dir := flag.String("dir", ".", "Path to the directory")

	flag.Parse()

	logger.PrintfInfoMsg("Starting server on port :" + *port)
	logger.PrintfInfoMsg("Path to the directory set: " + *dir)

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/gay", gay)

	loggedMux := logRequestMiddleware(mux)

	err = http.ListenAndServe(":"+*port, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}
