package logger

import (
	"log"
	"net/http"
)

// TODO: Save logs to the ./logs/triple-s.log path
type iLogger interface {
	PrintfInfoMsg(string, ...interface{})
	PrintfDebugMsg(string, ...interface{})
	PrintfErrorMsg(string, ...interface{})
	LogRequestMiddleware(http.Handler) http.Handler
}

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func printfMsg(level string, mes string, args ...interface{}) {
	log.Printf(level+" "+mes, args...)
}

func (l *Logger) PrintfInfoMsg(mes string, args ...interface{}) {
	printfMsg("[INFO]", mes, args...)
}

func (l *Logger) PrintfDebugMsg(mes string, args ...interface{}) {
	printfMsg("[DEBUG]", mes, args...)
}

func (l *Logger) PrintfErrorMsg(mes string, args ...interface{}) {
	printfMsg("[ERROR]", mes, args...)
}

func (l *Logger) LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
