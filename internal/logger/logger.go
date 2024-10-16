package logger

import (
	"log"
	"net/http"
)

func printfMsg(level string, depth int, mes string, args ...interface{}) {
	log.Printf(level+" "+mes, args...)
}

func PrintfInfoMsg(mes string, args ...interface{}) {
	printfMsg("[INFO]", 0, mes, args...)
}

func PrintfDebugMsg(mes string, args ...interface{}) {
	printfMsg("[DEBUG]", 0, mes, args...)
}

func PrintfErrorMsg(mes string, args ...interface{}) {
	printfMsg("[ERROR]", 0, mes, args...)
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
