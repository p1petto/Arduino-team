package main

import (
	"arduinoteam/internal/server"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1090"
	}

	server := server.NewServer()
	server.Run(port)
}
