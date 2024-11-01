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

func CustomUsage() {
	fmt.Println(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>] [-cfg <S>] 
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory
- --cfg S    Path to the config file`)
}

func init() {
	flag.StringVar(&port, "port", "4400", "Port number")
	flag.StringVar(&dir, "dir", "./data", "Path to the directory")
	flag.StringVar(&configPath, "cfg", "configs/server.yaml", "Path to the config file")

	flag.Usage = CustomUsage
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
