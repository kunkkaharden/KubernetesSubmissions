package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomValue := uuid.NewString()

	for {
		fmt.Printf("[%s] %s\n", time.Now().UTC().Format(time.RFC3339), randomValue)
		time.Sleep(5 * time.Second)
	}
}
