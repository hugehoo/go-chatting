package main

import (
	"chat_backend/config"
	"chat_backend/network"
	"chat_backend/repository"
	"chat_backend/service"
	"flag"
	"log"
)

var pathFlag = flag.String("config", "./config.toml", "config set up")
var port = flag.String("port", ":1010", "port set up")

func init() {
	log.Print("init first")
}

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		server := network.NewServer(service.NewService(rep), *port)
		server.StartServer()
	}

}
