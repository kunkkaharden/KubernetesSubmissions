package main

import (
	"cmp"
	"log"
	"os"
)

const defaultFile = "/data/log.txt"

func main() {
	filePath := cmp.Or(os.Getenv("LOG_FILE"), defaultFile)

	switch role := os.Getenv("ROLE"); role {
	case "writer":
		runWriter(filePath)
	case "server":
		runServer(filePath)
	default:
		log.Fatalf("ROLE must be either \"writer\" or \"server\", got %q", role)
	}
}
