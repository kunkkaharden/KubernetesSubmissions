package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// fetchPingPongCount asks Ping-pong how many pings it has served.
func fetchPingPongCount(ctx context.Context, client *http.Client, url string) (int, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status %s", response.Status)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, 64))
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		return 0, fmt.Errorf("invalid count %q: %w", body, err)
	}
	return count, nil
}
