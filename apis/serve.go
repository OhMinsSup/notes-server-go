package apis

import (
	"log"
	"net/http"
	"strings"

	"github.com/OhMinsSup/notes-server-go/core"
	"github.com/fatih/color"
)

type ServeOptions struct {
	ShowStartBanner bool
	HttpAddr        string
	HttpsAddr       string
	AllowedOrigins  []string // optional list of CORS origins (default to "*")
	BeforeServeFunc func(server *http.Server) error
}

// Serve starts a new app web server.
func Serve(app core.App, options *ServeOptions) error {
	if options == nil {
		options = &ServeOptions{}
	}

	if len(options.AllowedOrigins) == 0 {
		options.AllowedOrigins = []string{"*"}
	}

	// TODO: DB Migration

	router, err := InitApi(app)
	if err != nil {
		return err
	}

	if options.ShowStartBanner {
		schema := "http"
		if options.HttpsAddr != "" {
			schema = "https"
		}

		date := new(strings.Builder)
		log.New(date, "", log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s://%s", schema),
		)

		regular := color.New()
		regular.Printf(" ➜ REST API: %s\n", color.CyanString("%s://%s/api/", schema))
		regular.Printf(" ➜ Admin UI: %s\n", color.CyanString("%s://%s/_/", schema))
	}


	return router.Server.ListenAndServe()
}