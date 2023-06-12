package serve

import "net/http"

type ServeOptions struct {
	ShowStartBanner bool
	HttpAddr        string
	HttpsAddr       string
	AllowedOrigins  []string // optional list of CORS origins (default to "*")
	BeforeServeFunc func(server *http.Server) error
}
