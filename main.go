package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/korylprince/hasura-ad-webhook/httpapi"
)

func main() {
	config := new(httpapi.Config)
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}

	s := httpapi.NewServer(config)

	log.Println("Listening on:", config.ListenAddr)

	log.Println(http.ListenAndServe(config.ListenAddr, s.Router()))
}
