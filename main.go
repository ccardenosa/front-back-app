package main

import (
	"github.com/ccardenosa/front-back-app/backend"
	"github.com/ccardenosa/front-back-app/database"
	"github.com/ccardenosa/front-back-app/frontend"
)

var dbConfig = database.Config{
	ListenUri: "0.0.0.0:28082",
}

var beConfig = backend.Config{
	ListenUri:        "0.0.0.0:28081",
	DatabaseEndpoint: "127.0.0.1:28082",
}

var ftConfig = frontend.Config{
	ListenUri:       "0.0.0.0:28080",
	BackendEndpoint: "127.0.0.1:28081",
}

func main() {

	go database.StartDatabase(dbConfig)
	go backend.StartBackend(beConfig)
	frontend.StartFrontend(ftConfig)
}
