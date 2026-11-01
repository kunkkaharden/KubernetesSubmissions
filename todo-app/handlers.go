package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

type pageData struct {
	Title string
}

//go:embed templates/index.html
var templateFiles embed.FS

var pageTemplate = template.Must(template.ParseFS(templateFiles, "templates/index.html"))

func newHandler(cache *imageCache) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := pageTemplate.Execute(w, pageData{Title: "Todo App"}); err != nil {
			log.Printf("render page: %v", err)
		}
	})
	mux.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		image, err := cache.get(r.Context())
		if err != nil {
			log.Printf("get image: %v", err)
			http.Error(w, "Image is not available", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(image))
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(image); err != nil {
			log.Printf("write image response: %v", err)
		}
	})
	return mux
}
