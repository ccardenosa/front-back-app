package main

import (
	"github.com/ccardenosa/front-back-app/backend"
	"github.com/ccardenosa/front-back-app/database"
	"github.com/ccardenosa/front-back-app/frontend"
)

func main() {

	dbConfig := database.Config{
		ListenUri: "0.0.0.0:28082",
	}
	go database.StartDatabase(dbConfig)

	beConfig := backend.Config{
		ListenUri:        "0.0.0.0:28081",
		DababaseEndpoint: "127.0.0.1:28082",
	}
	go backend.StartBackend(beConfig)

	ftConfig := frontend.Config{
		ListenUri:       "0.0.0.0:28080",
		BackendEndpoint: "127.0.0.1:28081",
	}
	frontend.StartFrontend(ftConfig)
}
