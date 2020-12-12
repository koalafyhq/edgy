package main

import (
	"log"

	"github.com/koalafy/edgy/internal/helpers"
	"github.com/koalafy/edgy/internal/server"
)

func init() {
	if helpers.GetIPFSGateway() == "" {
		log.Fatal("IPFS_GATEWAY need to be initialized")
	}
}

func main() {
	app := server.New()

	if err := server.Run(app); err != nil {
		log.Println(err)
	}
}
