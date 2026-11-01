package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// newTodoProxy forwards /todos requests to the todo backend. The browser talks
// to this app only, so no Ingress or CORS configuration is needed for the API.
func newTodoProxy(target string) (http.Handler, error) {
	backendURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	if backendURL.Scheme == "" || backendURL.Host == "" {
		return nil, errors.New("TODO_BACKEND_URL must be an absolute URL")
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy todo request: %v", err)
		http.Error(w, "Todo backend is not available", http.StatusBadGateway)
	}
	return proxy, nil
}
