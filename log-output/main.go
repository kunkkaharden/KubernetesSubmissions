package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const defaultPort = "8080"

func currentStatus(randomValue string) string {
	return fmt.Sprintf("[%s] %s", time.Now().UTC().Format(time.RFC3339), randomValue)
}

func statusHandler(randomValue string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, currentStatus(randomValue))
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid PORT %q: must be a number between 1 and 65535", port)
	}

	randomValue := uuid.NewString()
	http.Handle("/", statusHandler(randomValue))

	go func() {
		for {
			fmt.Println(currentStatus(randomValue))
			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Printf("Server started in port %d\n", portNumber)
	if err := http.ListenAndServe(":"+strconv.Itoa(portNumber), nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
