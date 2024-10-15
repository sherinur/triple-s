package logger

import (
	"log"
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
