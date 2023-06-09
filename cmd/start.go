package main

import (
	"log"
	"net/http"

	"github.com/OhMinsSup/notes-server-go/apis"
	"github.com/OhMinsSup/notes-server-go/core"
)

func Start() {
	var allowedOrigins []string
	var httpAddr string
	var httpsAddr string
	var app core.App

	err := apis.Serve(app, &apis.ServeOptions{
		HttpAddr:        httpAddr,
		HttpsAddr:       httpsAddr,
		ShowStartBanner: true,
		AllowedOrigins:  allowedOrigins,
	})

	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}