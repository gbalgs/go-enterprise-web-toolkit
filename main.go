package main

import (
	"flag"
	"fmt"
	"github.com/wen-bing/go-enterprise-web-toolkit/server"
	"os"
)

func main() {
	env := flag.String("e", "development", "")
	configDir := flag.String("c", ".", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode} -c {config file dir}")
		os.Exit(1)
	}
	flag.Parse()
	app := server.New(*env, *configDir)
	app.Start()
}
