package main

import "github.com/wen-bing/go-enterprise-web-toolkit/server"

func main() {
	app := server.New("")
	app.Start()
}
