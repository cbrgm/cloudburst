package main

import (
	"log"

	"github.com/zalando/skipper"
	"github.com/zalando/skipper/routing"
)

func main() {

	dataclient, err  := NewCloudburst()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(skipper.Run(skipper.Options{
		Address: ":9090",
		CustomDataClients: []routing.DataClient{dataclient},
	}))
}