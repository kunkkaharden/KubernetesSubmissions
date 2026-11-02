package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type imageCache struct {
	mu    sync.Mutex
	path  string
	ttl   time.Duration
	fetch func(context.Context) ([]byte, error)
	now   func() time.Time
}

func newImageCache(dir string, ttl time.Duration, fetch func(context.Context) ([]byte, error)) *imageCache {
	return &imageCache{
		path:  filepath.Join(dir, "image"),
		ttl:   ttl,
		fetch: fetch,
		now:   time.Now,
	}
}

// get returns the saved image. It downloads and saves a new one only when there
// is no image yet or the saved one is older than the configured TTL. If the
// download fails, the old image is shown.
func (c *imageCache) get(ctx context.Context) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	saved, _ := os.ReadFile(c.path)
	info, err := os.Stat(c.path)
	if err == nil && len(saved) > 0 && c.now().Sub(info.ModTime()) < c.ttl {
		return saved, nil
	}

	image, err := c.fetch(ctx)
	if err != nil {
		if len(saved) > 0 {
			return saved, nil
		}
		return nil, err
	}
	if err := c.save(image); err != nil {
		log.Printf("save image: %v", err)
	}
	return image, nil
}

func (c *imageCache) save(image []byte) error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(c.path, image, 0o644); err != nil {
		return err
	}
	return os.Chtimes(c.path, c.now(), c.now())
}

func fetchImage(ctx context.Context, imageURL string, timeout time.Duration, maxImageSize int64) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New("image provider returned " + response.Status)
	}
	if contentType := response.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "image/") {
		return nil, errors.New("image provider returned a non-image response")
	}

	image, err := io.ReadAll(io.LimitReader(response.Body, maxImageSize+1))
	if err != nil {
		return nil, err
	}
	if len(image) == 0 || int64(len(image)) > maxImageSize {
		return nil, errors.New("image provider returned an empty or oversized image")
	}
	return image, nil
}
