package main

import (
	"flag"
	"fmt"
	"os"

	"triple-s/internal/server"
)

var (
	configPath string
	port       string
	dir        string
)

func init() {
	flag.StringVar(&port, "port", "4400", "Port number")
	flag.StringVar(&dir, "dir", ".", "Path to the directory")
	flag.StringVar(&configPath, "cfg", "configs/server.yaml", "Path to config file")
}

func main() {
	flag.Parse()

	port = ":" + port

	cfg := server.NewConfig(configPath, port, dir)

	apiServer := server.New(cfg)
	err := apiServer.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
