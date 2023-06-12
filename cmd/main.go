package main

import (
	"log"
	"net/http"

	"github.com/OhMinsSup/notes-server-go/apis"
	"github.com/OhMinsSup/notes-server-go/tools/config"
	"github.com/OhMinsSup/notes-server-go/tools/serve"
)

func main() {
	config, err := config.ReadConfigFile("./config.json")
	log.Println("config: ", config)
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
		return
	}
	options := &apis.BaseAppConfig{
		Config: config,
		ServerOptions: &serve.ServeOptions{
			ShowStartBanner: true,
			HttpAddr:        config.ServerRoot,
			HttpsAddr:       "",
			AllowedOrigins:  []string{"*"},
			BeforeServeFunc: func(server *http.Server) error {
				log.Println("BeforeServeFunc")
				return nil
			},
		},
		DataMaxOpenConns: 120,
		DataMaxIdleConns: 10,
	}
	bootstrap := apis.NewBaseApp(options)

	err = bootstrap.Bootstrap()
	if err != nil {
		log.Fatal("Unable to bootstrap the application: ", err)
		return
	}

	log.Panicln("Application terminated")
}
