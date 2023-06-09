package main

import (
	"log"

	"github.com/OhMinsSup/notes-server-go/tools/config"
)

func main() {
	config, err := config.ReadConfigFile("./config.json")
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
		return
	}

	log.Fatal(config)
}
