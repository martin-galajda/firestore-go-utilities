package httputils

import (
	"net/http"
)

// NewHTTPClient returns new instance of http.Client
func NewHTTPClient() *http.Client {
	return makeHTTPClient()
}

func makeHTTPClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}
