package main

import (
	"flag"
	"retroHub/data/json"
	"retroHub/server"
)

func main() {
	contentPath := flag.String("content", "content.json", "Set the content configuration path, this file is encoded in JSON")
	serverPort := flag.Int("port", 8080, "Set the HTTP Server port")

	if contentPath == nil {
		contentPath = new(string)
		*contentPath = "content.json"
	}
	if serverPort == nil {
		serverPort = new(int)
		*serverPort = 8080
	}

	provider, err := json.New(*contentPath)
	if err != nil {
		panic(err)
	}
	err = server.Serve(provider, uint(*serverPort))
	if err != nil {
		panic(err)
	}
}
