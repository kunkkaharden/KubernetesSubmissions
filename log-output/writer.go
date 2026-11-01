package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

const writePeriod = 5 * time.Second

func runWriter(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("open log file %q: %v", filePath, err)
	}
	defer file.Close()

	randomValue := uuid.NewString()
	writeLine := func() {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		if _, err := fmt.Fprintf(file, "[%s] %s\n", timestamp, randomValue); err != nil {
			log.Fatalf("write log file %q: %v", filePath, err)
		}
	}

	writeLine()
	ticker := time.NewTicker(writePeriod)
	defer ticker.Stop()
	for range ticker.C {
		writeLine()
	}
}
